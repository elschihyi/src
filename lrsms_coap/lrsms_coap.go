package lrsms_coap

import (
	"log"
	"net"
  "go-coap"
  "lrsms"
  "container/list"
  "encoding/json"
  "time"
	"regexp"
	"math/rand"
	//"github.com/dustin/go-coap"
)

const (
  //modify connected device and resource
  Dev string = "Dev"

  //modify resource info
  Ref string ="Ref"
)

//******************************************************************************
//Public Functions
//******************************************************************************
func CoAPServerStart(port string){
  mylrsms := lrsms.NewLRSMS()
  log.Fatal(coap.ListenAndServe("udp", port,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			//log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
      //log.Printf(port+" Got message path=%q: from %v", m.Path(), a)
      //log.Printf("payload is: %v", string(m.Payload))

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
        createTime, _ := time.Parse(time.RFC3339, payload["CreateTime"].(string))
        mylrsms.CreateRef(payload["URI"].(string), dependedList, &createTime,
					GetResource, Update, Alert)
        //log.Printf("in register resource")
      }

      //update alert Ref
			if m.Path()[0] == Ref && m.Code == coap.PUT {
			  createTime, _ := time.Parse(time.RFC3339, payload["CreateTime"].(string))
				mylrsms.RecieveUpdateFromInside(payload["URI"].(string), &createTime)
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
//Public Call Back Functions
//******************************************************************************
func GetResource(resourceURI string)[]byte{
	log.Printf("get func called: "+ resourceURI)
	return nil
}

func Update(resourceURI string){
	//log.Printf("update func called: "+ resourceURI)
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
