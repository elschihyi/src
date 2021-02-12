package main

import (
	"log"
  "go-coap"
  "math/rand"
  "container/list"
  "lrsms_util"
  "fmt"
  "time"
	//"github.com/dustin/go-coap"
)

const(
  localhost string = "localhost"

  device1AppServPort string = ":5683"
  device1LRSNSServPort string = ":5684" //default LRSMS port
  a1ID string = "app01"
  a2ID string = "app02"

  device2AppServPort string = ":5685"
  device2LRSNSServPort string = ":5686"
  a3ID string = "app03"
)

var ConnectedDevices *list.List
var UnConnectedDevices *list.List

func main() {
  ConnectedDevices = list.New()
  UnConnectedDevices = list.New()

  //****************************************************************************
  //1. create resources
  //Create a device with two app: app01 app02
  //app01 has two resources: Resource1 and Resource2.
  //Resource2 is depending on Resource1
  //app02 has three resources: Resource2(Cache fron app01), Resource3, and Resource4
  //resouce 4 is depending on Resource2 and Resouce3
  //****************************************************************************
  log.Printf("1. Create Resources")

  //create device and apps and resoruces and make it connected Device after
  myDevice1 := lrsms_util.StartDevice(device1AppServPort,device1LRSNSServPort)
  myDevice1Ele := UnConnectedDevices.PushBack(myDevice1)

  //create apps
  myDevice1.AddApp(a1ID)
  myDevice1.AddApp(a2ID)

  //create resources
  r1ID := localhost+device1AppServPort+"/"+a1ID+"/"+"Resource1"
  r1_depended_list := list.New()
  resource1 :=  lrsms_util.NewResource(r1ID, r1_depended_list,
    time.Now().Format(time.RFC3339), []byte("BodyofR1"))
  myDevice1.CreateResource(a1ID, resource1)

  r2ID := localhost+device1AppServPort+"/"+a1ID+"/"+"Resource2"
  r2_depended_list := list.New()
  r2_depended_list.PushBack(resource1.URI)
  resource2 :=  lrsms_util.NewResource(r2ID, r2_depended_list,
    time.Now().Format(time.RFC3339), []byte("BodyofR1BodyofR2"))
  myDevice1.CreateResource(a1ID, resource2)

  r2CacheID := localhost+device1AppServPort+"/"+a2ID+"/"+"Resource2"
  r2Cache_depended_list := list.New()
  r2Cache_depended_list.PushBack(resource2.URI)
  r2Cache :=  lrsms_util.NewResource(r2CacheID, r2Cache_depended_list,
     resource2.CreateTime, []byte("BodyofR1BodyofR2"))
  myDevice1.CreateResource(a2ID, r2Cache)

  r3ID := localhost+device1AppServPort+"/"+a2ID+"/"+"Resource3"
  r3_depended_list := list.New()
  resource3 :=  lrsms_util.NewResource(r3ID, r3_depended_list,
    time.Now().Format(time.RFC3339), []byte("BodyofR3"))
  myDevice1.CreateResource(a2ID, resource3)

  r4ID := localhost+device1AppServPort+"/"+a2ID+"/"+"Resource4"
  r4_depended_list := list.New()
  r4_depended_list.PushBack(resource3.URI)
  r4_depended_list.PushBack(r2Cache.URI)
  resource4 :=  lrsms_util.NewResource(r4ID, r4_depended_list,
    time.Now().Format(time.RFC3339), []byte("BodyofR4BodyofR3"))
  myDevice1.CreateResource(a2ID, resource4)
  //****************************************************************************
  //2. updates
  //****************************************************************************
  log.Printf("")
  log.Printf("2. Update")

  time.Sleep(2 * time.Second)

  myDevice1.UpdateResource(a1ID, resource1.URI)

  //time.Sleep(4 * time.Second)

  myDevice1.UpdateResource(a2ID, resource3.URI)

  //printLRSMS(device1LRSNSServPort)
  time.Sleep(2 * time.Second)

  //****************************************************************************
  //3. sync
  //Create a device with one app: app03
  //app03 has two resources: Resource2 and Resource5.
  //Resource5 is depending on Resource2
  //****************************************************************************
  log.Printf("")
  log.Printf("3. Sync")
  //create another device, apps, resources and connected Device
  myDevice2 := lrsms_util.StartDevice(device2AppServPort,device2LRSNSServPort)
  myDevice2Ele := UnConnectedDevices.PushBack(myDevice2)

  myDevice2.AddApp(a3ID)

  r2CacheID = localhost+device2AppServPort+"/"+a3ID+"/"+"Resource2"
  r2Cache_depended_list = list.New()
  r2Cache_depended_list.PushBack(resource2.URI)
  resource2CreateTime, _ := time.Parse(time.RFC3339, resource2.CreateTime)
  r2CreateTimeMinus20Min := resource2CreateTime.Add(time.Minute * -20)
  r2Cache =  lrsms_util.NewResource(r2CacheID, r2Cache_depended_list,
     r2CreateTimeMinus20Min.Format(time.RFC3339), []byte("BodyofR1BodyofR2"))
  myDevice2.CreateResource(a3ID, r2Cache)

  r5ID := localhost+device2AppServPort+"/"+a3ID+"/"+"Resource5"
  r5_depended_list := list.New()  //localhost:5683/app02/Resource4
  r5_depended_list.PushBack(r2Cache.URI)
  resource5 :=  lrsms_util.NewResource(r5ID, r5_depended_list,
    time.Now().Format(time.RFC3339), []byte("BodyofR5BodyofR2"))
  myDevice2.CreateResource(a3ID, resource5)

  //printLRSMS(device2LRSNSServPort)

  //make device 1 connected
  UnConnectedDevices.Remove(myDevice1Ele)
  ConnectedDevices.PushBack(myDevice1)
  myDevice1.Connect(ConnectedDevices)

  //make device 2 connected and trigger sync
  UnConnectedDevices.Remove(myDevice2Ele)
  ConnectedDevices.PushBack(myDevice2)
  myDevice2.Connect(ConnectedDevices)
  time.Sleep(2 * time.Second)
  //****************************************************************************
  //4. cross device updates
  //****************************************************************************
	log.Printf("")
  log.Printf("4. Cross Device Updates")

  //4.1 create resoruce
	r3CacheID := localhost+device2AppServPort+"/"+a3ID+"/"+"Resource3"
  r3Cache_depended_list := list.New()
  r3Cache_depended_list.PushBack(resource3.URI)
  r3Cache :=  lrsms_util.NewResource(r3CacheID, r3Cache_depended_list,
     resource3.CreateTime, resource3.Content)
  myDevice2.CreateResource(a3ID, r3Cache)
  time.Sleep(2 * time.Second)

  //4.2 update resource
  myDevice1.UpdateResource(a1ID, resource1.URI)
  time.Sleep(2 * time.Second)

	//4.3 delete resource
  //myDevice2.DeleteResource(a3ID, r3CacheID)
	myDevice2.DeleteResource(a3ID, r5ID)
  time.Sleep(2 * time.Second)

	//4.4 update resource
  myDevice1.UpdateResource(a1ID, resource1.URI)
	time.Sleep(2 * time.Second)

	//restore resource 5
  myDevice2.CreateResource(a3ID, resource5)
	time.Sleep(1 * time.Second)
	//****************************************************************************
  //5. Disconnect device 2 update resource 1 than reconnect.  check sync
  //****************************************************************************
	log.Printf("")
  log.Printf("5. Disconnected")

	//5.1. disconnect a device
	ConnectedDevices.Remove(myDevice2Ele)
  UnConnectedDevices.PushBack(myDevice2)
  myDevice2.Disconnect(ConnectedDevices)
  time.Sleep(2 * time.Second)

	//5.2 update a resouce in connected deviceURL
	myDevice1.UpdateResource(a1ID, resource1.URI)
  time.Sleep(3 * time.Second)


	//5.3 reconnect device
	UnConnectedDevices.Remove(myDevice2Ele)
  ConnectedDevices.PushBack(myDevice2)
  myDevice2.Connect(ConnectedDevices)
  time.Sleep(2 * time.Second)

	//5.4 update a resource again
	myDevice1.UpdateResource(a1ID, resource1.URI)
  time.Sleep(2 * time.Second)

	//****************************************************************************
  //Halt
  //****************************************************************************
  //hault til user input e
  for input := "";input!="e";{
    fmt.Println("Enter 'e' to terminate")
    fmt.Scanf("%s", &input)
    //fmt.Println("Hello", name)
  }
}//end main

func connectDevice(device lrsms_util.Device){
}

func disConnectDevice(device lrsms_util.Device){
}

func printLRSMS(port string){
  req := coap.Message{
    Type:      coap.Confirmable,
    Code:      coap.Content,
  	MessageID: uint16(rand.Intn(10000)),
  	Payload:   []byte("{}"),
  }
  req.SetPathString("/Dev")
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
