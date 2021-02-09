package lrsms_coap

import (
	"log"
	"net"
  "go-coap"
  "lrsms"
  "container/list"
  "encoding/json"
  "time"
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
  log.Fatal(coap.ListenAndServe("udp", port, //":5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			//log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
      log.Printf(port+" Got message path=%q: from %v", m.Path(), a)
      log.Printf("payload is: %v", string(m.Payload))

      if m.Path()[0] == Ref && m.Code == coap.POST {
        var payload map[string]interface{}
        if err := json.Unmarshal(m.Payload, &payload); err != nil {
          panic(err)
        }

        dependedList := list.New()
        dependedArray := payload["Depended"].([]interface{})
        for _, s := range dependedArray {
          dependedList.PushBack(s.(string))
        }
        createTime, _ := time.Parse(time.RFC3339, payload["CreateTime"].(string))
        mylrsms.CreateRef(payload["URI"].(string), dependedList, &createTime, nil, nil, nil)
        //mylrsms.CreateRef("test", list.New(), nil, nil, nil, nil)
        //log.Printf("in register resource")
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
