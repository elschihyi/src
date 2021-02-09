package lrsms_util
import (
  //"container/list"
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
}

const (
  //modify connected device and resource
  Dev string = "/Dev"

  //modify resource info
  Ref string ="/Ref"
)

//******************************************************************************
//Public Functions
//******************************************************************************
func StartDevice(appServerPort string, lrsmsServerPort string)*Device{
  var nDevice Device
  nDevice.Apps = make(map[string]map[string]Resource)
  nDevice.AppServerPort = appServerPort
  nDevice.LRSMSServerPort = lrsmsServerPort
  go startLRSMS(lrsmsServerPort) //start lrsms
  go startAppServer(appServerPort) //start StartAppServer
  return &nDevice
}

func (device Device)AddApp(appID string, resources map[string]Resource){
  device.Apps[appID]=resources
  //register resources
  for _, v := range resources {
    sendCoAP(device.LRSMSServerPort, Ref, coap.POST, makeResourceJson(v))
  }
}

func (device Device)UpdateResource(appID string, resourceID string,
   newContent []byte){
  resource := device.Apps[appID][resourceID]
  resource.Content = newContent
  sendCoAP(device.LRSMSServerPort, Ref, coap.PUT,[]byte(""))
}

func (device Device)DeleteResource(appID string, resourceID string){
  delete(device.Apps[appID],resourceID)
  sendCoAP(device.LRSMSServerPort, Ref, coap.DELETE,[]byte(""))
}

func (device Device)CreateResource(appID string, resource *Resource){
  device.Apps[appID][resource.URI] = *resource
  sendCoAP(device.LRSMSServerPort, Ref, coap.POST, makeResourceJson(*resource))
}

//******************************************************************************
//Private Functions
//******************************************************************************
func startLRSMS(lrsmsServerPort string){
  lrsms_coap.CoAPServerStart(lrsmsServerPort)
}

func startAppServer(appServerPort string){
  log.Fatal(coap.ListenAndServe("udp", appServerPort, //":5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			return nil
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

func makeResourceJson(resource Resource)[]byte{
  dependedArray := make([]string,resource.Depended.Len())
  i := 0
  for e := resource.Depended.Front(); e != nil; e = e.Next() {
    dependedArray[i]=e.Value.(string)
    i++
  }
  theResource := &resourceJson{
    URI :       resource.URI,
    Depended:   dependedArray,
    CreateTime: resource.CreateTime.Format(time.RFC3339)}
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
	c, err := coap.Dial("udp", "localhost"+port)
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
