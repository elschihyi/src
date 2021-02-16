package lrsms_util

import (
  "container/list"
  "go-coap"
  "log"
  "net"
  "lrsms_coap"
  "math/rand"
  "encoding/json"
  "time"
)

type Device struct{
  Apps map[string]map[string]Resource  //[appID]map[resourceID]Resource
  AppServerPort string
  LRSMSServerPort string
  Channel chan interface{}
}

const (
  Dev string = "Dev" //modify connected device and resource
  Ref string = "Ref" //modify resource info
	Res string = "Res" //modify resource
)

const (
  localhost string = "localhost"
)

//******************************************************************************
//Public Functions
//******************************************************************************
func StartDevice(appServerPort string, lrsmsServerPort string)*Device{
  var nDevice Device
  nDevice.Apps = make(map[string]map[string]Resource)
  nDevice.AppServerPort = appServerPort
  nDevice.LRSMSServerPort = lrsmsServerPort
  //nDevice.Connected = false
  go startLRSMS(lrsmsServerPort) //start lrsms
  go startAppServer(appServerPort, &nDevice) //start StartAppServer
  //log.Printf("Device %v started", localhost+nDevice.AppServerPort)
  return &nDevice
}

func (device Device) Connect(ConnectedDevices *list.List){
  log.Printf("Device %v connected", localhost+device.AppServerPort)
  //sync with all connected device
  for e := ConnectedDevices.Front(); e != nil; e = e.Next() {
    otherDevice := e.Value.(*Device)
     //no need to sync with self
	   if otherDevice.LRSMSServerPort == device.LRSMSServerPort {
       break
     }

     //make device sync with the otherDevice
     //1. tell device a new Otherdevice is connected
     otherDeviceURL := localhost+otherDevice.LRSMSServerPort
     mapOtherDeviceURL := map[string]string{"LRSMSServerAddress": otherDeviceURL}
     otherDevicejsonByte, _ := json.Marshal(mapOtherDeviceURL)
     sendCoAP(device.LRSMSServerPort, Dev, coap.POST, otherDevicejsonByte)

     //2. tell otherdevice a new device is connected
     deviceURL := localhost+device.LRSMSServerPort
     mapDeviceURL := map[string]string{"LRSMSServerAddress": deviceURL}
     devicejsonByte, _ := json.Marshal(mapDeviceURL)
     sendCoAP(otherDevice.LRSMSServerPort, Dev, coap.POST, devicejsonByte)
	}
  //device.Connected = true
}


func (device Device) Disconnect(ConnectedDevices *list.List){
  //log.Printf("Device %v Disconnected", localhost+device.AppServerPort)
  for e := ConnectedDevices.Front(); e != nil; e = e.Next() {
    otherDevice := e.Value.(*Device)
    //otherDeviceURL := localhost+otherDevice.LRSMSServerPort
    deviceURL := localhost+device.LRSMSServerPort
    mapDeviceURL := map[string]string{"DeviceURL": deviceURL}
    devicejsonByte, _ := json.Marshal(mapDeviceURL)
    sendCoAP(otherDevice.LRSMSServerPort, Dev, coap.DELETE, devicejsonByte)
  }
}

func (device Device)AddApp(appID string){
  device.Apps[appID] = make(map[string]Resource)
}

func (device Device)CreateResource(appID string, resource *Resource){
  device.Apps[appID][resource.URI] = *resource
  sendCoAP(device.LRSMSServerPort, Ref, coap.POST,
    device.makeResourceJson(appID,*resource))
  //log.Printf("Resource %v manual created in %v", resource.URI, appID)
}

func (device Device)AlertResource(appID string, resourceID string){
  resource := device.Apps[appID][resourceID]
  resource.Alert()
}

func (device Device)UpdateResource(appID string, resourceID string){
  //log.Printf("Resource %v in %v update", resourceID, appID)
  resource := device.Apps[appID][resourceID]
  resource.Update()
  sendCoAP(device.LRSMSServerPort, Ref, coap.PUT, device.makeResourceJson(appID ,resource))
  if (device.Channel != nil){
    device.Channel <- time.Now()
  }
}

func (device Device)UpdateCacheResource(appID string, resourceID string,
  createTime string, content string){
  resource := device.Apps[appID][resourceID]
  resource.Content = []byte(content)
  resource.CreateTime = createTime
  log.Printf("Resource %v update in %v", resource.URI, appID)
  sendCoAP(device.LRSMSServerPort, Ref, coap.PUT, device.makeResourceJson(appID ,resource))
  //log.Printf("Resource %v manual update in %v", resource.URI, appID)
}

func (device Device)DeleteResource(appID string, resourceID string){
  payload := device.makeResourceJson(appID ,device.Apps[appID][resourceID])
  delete(device.Apps[appID], resourceID)
  log.Printf("Resource %v deleted in %v", resourceID, appID)
  sendCoAP(device.LRSMSServerPort, Ref, coap.DELETE,payload)
}

//******************************************************************************
//Private Functions
//******************************************************************************
func startLRSMS(lrsmsServerPort string){
  lrsms_coap.CoAPServerStart(lrsmsServerPort)
}

func startAppServer(appServerPort string, myDevice *Device){
  log.Fatal(coap.ListenAndServe("udp", appServerPort, //":5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
      //log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
      //log.Printf(port+" Got message path=%q: from %v", m.Path(), a)
      //log.Printf("payload is: %v", string(m.Payload))

			var payload map[string]interface{}
			if err := json.Unmarshal(m.Payload, &payload); err != nil {
				panic(err)
			}

      if payload["Action"].(string) == "Alert" && m.Code == coap.PUT {
        appID := m.Path()[0]
        resourceID := localhost+appServerPort+"/"+m.Path()[0]+"/"+m.Path()[1]
        go myDevice.AlertResource(appID, resourceID)
      }

      if payload["Action"].(string) == "Update" && m.Code == coap.PUT {
        appID := m.Path()[0]
        resourceID := localhost+appServerPort+"/"+m.Path()[0]+"/"+m.Path()[1]
        go myDevice.UpdateResource(appID, resourceID)
      }

      if payload["Action"].(string) == "UpdateCache" && m.Code == coap.PUT {
        appID := m.Path()[0]
        resourceID := localhost+appServerPort+"/"+m.Path()[0]+"/"+m.Path()[1]
        go myDevice.UpdateCacheResource(appID, resourceID,
          payload["CreateTime"].(string), payload["Content"].(string))
      }

      res := &coap.Message{
        Type:      coap.Acknowledgement,
        Code:      coap.Content,
        MessageID: m.MessageID,
        //Token:     m.Token,
        Payload:   []byte(""),
      }
      return res
		})))
}
//******************************************************************************
//Private utility Functions
//******************************************************************************
type resourceJson struct {
    URI        string
    Depended   []string
    CreateTime string
}

func (device Device) makeResourceJson(appID string, resource Resource)[]byte{
  //log.Printf("appID %v reosuceID %v", appID, resource.URI)
  dependedArray := make([]string, resource.Depended.Len())
  i := 0
  for e := resource.Depended.Front(); e != nil; e = e.Next() {
    dependedArray[i]=e.Value.(string)
    i++
  }
  theResource := &resourceJson{
    URI :       resource.URI,
    Depended:   dependedArray,
    CreateTime: resource.CreateTime}
  resourceJsonByte, _ := json.Marshal(theResource)
  return resourceJsonByte
}

func sendCoAP(port string, path string, code coap.COAPCode, payload []byte){
  req := coap.Message{
		Type:      coap.Confirmable,
		Code:      code,
		MessageID: uint16(rand.Intn(10000)),
		Payload:   payload,
	}
  req.SetPathString(path)
	c, err := coap.Dial("udp", localhost+port)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	if rv != nil {
		//log.Printf("Response payload: %s", rv.Payload)
	}
}
