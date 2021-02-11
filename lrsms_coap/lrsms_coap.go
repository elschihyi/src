package lrsms_coap

import (
	"log"
	"net"
  "go-coap"
  "lrsms"
  "container/list"
  "encoding/json"
  //"time"
	"regexp"
	"math/rand"
	//"github.com/dustin/go-coap"
)

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
func CoAPServerStart(port string){
  mylrsms := lrsms.NewLRSMS(GetDeviceRefs, SentRefs, GetResourceOtherDevice,
		SentResource, UpdateCache)

  log.Fatal(coap.ListenAndServe("udp", port,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			//log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
      //log.Printf(port+" Got message path=%q: from %v", m.Path(), a)
      //log.Printf("payload is: %v", string(m.Payload))

			//get payload
			var payload map[string]interface{}
			if err := json.Unmarshal(m.Payload, &payload); err != nil {
				panic(err)
			}

			//create Ref
      if m.Path()[0] == Ref && m.Code == coap.POST {
        dependedList := list.New()
        dependedArray := payload["Depended"].([]interface{})
        for _, s := range dependedArray {
          dependedList.PushBack(s.(string))
        }
        mylrsms.CreateRef(payload["URI"].(string), dependedList,
				payload["CreateTime"].(string), GetResource, Update, Alert)
        //log.Printf("in register resource")
      }

      //update Ref
			if m.Path()[0] == Ref && m.Code == coap.PUT {
			mylrsms.RecieveUpdateFromInside(payload["URI"].(string),
				payload["CreateTime"].(string))
      }

      //create new connected device
      if m.Path()[0] == Dev && m.Code == coap.POST {
				//payload["LRSMSServerAddress"].(string)
				//slog.Printf("Recieved new Device at port: %v", port)
		    mylrsms.NewDevice(payload["LRSMSServerAddress"].(string), localhost+port)
			}

			//get resourceRefs list and its create Time
			if m.Path()[0] == Ref && m.Code == coap.GET {
				//log.Printf("recieved want resourceRefs at port: %v", port)
		    mylrsms.GetRefs(payload["ToAddress"].(string), localhost+port)
			}

      //update connected device Resources and compare createtime
			if m.Path()[0] == Dev && m.Code == coap.PUT {
				//log.Printf("Recieved resourceRefs at port: %v", port)
				mylrsms.UpdateOtherDeviceRes(payload["ToAddress"].(string),
				localhost+port, payload["ResRefs"].(map[string]interface{}))
			}

			//get Resouce Cache
			if m.Path()[0] == Res && m.Code == coap.GET {
				//log.Printf("Recieved want Resource at port: %v", port)
				mylrsms.GetResource(payload["ToAddress"].(string), localhost+port,
				payload["ResourceURI"].(string))
			}

      //update resource
			if m.Path()[0] == Res && m.Code == coap.PUT {
				//log.Printf("Recieved update Resource at port: %v", port)
				mylrsms.UpdateResource(payload["ToAddress"].(string), localhost+port,
				payload["ResourceURI"].(string), payload["CreateTime"].(string),
				payload["ContentString"].(string))
			}

			//print lrsms
			if m.Code == coap.Content {
				mylrsms.Print()
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
//Public Device Call Back Functions
//******************************************************************************

//******************************************************************************
//Public Resource Ref Call Back Functions
//******************************************************************************
func GetDeviceRefs(otherAdd string, hostAdd string){
  //log.Printf("GetDeviceRes called: "+ newDevice)
	mapD := map[string]string{"ToAddress": hostAdd}
  payload, _ := json.Marshal(mapD)
	sendCoAP(otherAdd, Ref, coap.GET,payload)
}

type sentRefsResponse struct{
	ToAddress string
	ResRefs map[string]string
}

func SentRefs(otherAdd string,  hostAdd string, simepleResourceRefs map[string]string){
	sRR := &sentRefsResponse {
    ToAddress:   hostAdd,
    ResRefs: simepleResourceRefs}
	payload, _ := json.Marshal(sRR)
	sendCoAP(otherAdd, Dev, coap.PUT,payload)
}

//******************************************************************************
//Public Resource Call Back Functions
//******************************************************************************
func SentResource(otherAdd string, hostAdd string, resourceURI string,
	createTime string, content string){
	mapD := map[string]string{"ToAddress": hostAdd, "ResourceURI":resourceURI,
	"CreateTime": createTime, "ContentString":content}
	payload, _ := json.Marshal(mapD)
	sendCoAP(otherAdd, Res, coap.PUT, payload)
}

func GetResourceOtherDevice(otherAdd string,  hostAdd string, resourceURI string){
	//log.Printf("get func called: "+ resourceURI)
	mapD := map[string]string{"ToAddress": hostAdd, "ResourceURI":resourceURI}
  payload, _ := json.Marshal(mapD)
	sendCoAP(otherAdd, Res, coap.GET, payload)
	//return nil
}

func GetResource(resourceURI string)[]byte{
	return nil
}

func Update(resourceURI string){
	//log.Printf("update func called: %v", resourceURI)
	resourceURIArray := regexp.MustCompile(`/`).Split(resourceURI, -1)
	path := "/"+resourceURIArray[1]+"/"+resourceURIArray[2]
  Payload, _ := json.Marshal(map[string]string{"Action": "Update"})
	sendCoAP(resourceURIArray[0], path, coap.PUT, Payload)
}

func Alert(resourceURI string){
	//log.Printf("alert func called: "+ resourceURI)
	resourceURIArray := regexp.MustCompile(`/`).Split(resourceURI, -1)
	path := "/"+resourceURIArray[1]+"/"+resourceURIArray[2]
  Payload, _ := json.Marshal(map[string]string{"Action": "Alert"})
	sendCoAP(resourceURIArray[0], path, coap.PUT, Payload)
}

func UpdateCache(resourceCacheURI string, createTime string, content string){
	resourceURIArray := regexp.MustCompile(`/`).Split(resourceCacheURI, -1)
	path := "/"+resourceURIArray[1]+"/"+resourceURIArray[2]
  Payload, _ := json.Marshal(map[string]string{"Action": "UpdateCache", "CreateTime":createTime, "Content":content })
	sendCoAP(resourceURIArray[0], path, coap.PUT, Payload)
}
//******************************************************************************
//Private utility Functions
//******************************************************************************
func sendCoAP(host string, path string, code coap.COAPCode, payload []byte){
  req := coap.Message{
		Type:      coap.Confirmable,
		Code:      code,
		MessageID: uint16(rand.Intn(10000)),
		Payload:   payload,
	}
  req.SetPathString(path)
	c, err := coap.Dial("udp", host)
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
