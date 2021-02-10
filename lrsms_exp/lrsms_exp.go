package main

import (
	//"log"
  //"go-coap"
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

  //create device and apps and resoruces and make it connected Device after
  myDevice1 := lrsms_util.StartDevice(device1AppServPort,device1LRSNSServPort)
  myDevice1Ele := UnConnectedDevices.PushBack(myDevice1)

  //create apps
  myDevice1.AddApp(a1ID)
  myDevice1.AddApp(a2ID)

  //create resources
  r1ID := localhost+device1AppServPort+"/"+a1ID+"/"+"Resource1"
  r1_depended_list := list.New()
  curTime := time.Now()
  resource1 :=  lrsms_util.NewResource(r1ID, r1_depended_list, &curTime,
    []byte("BodyofR1"))
  myDevice1.CreateResource(a1ID, resource1)

  r2ID := localhost+device1AppServPort+"/"+a1ID+"/"+"Resource2"
  r2_depended_list := list.New()
  r2_depended_list.PushBack(resource1.URI)
  curTime = time.Now()
  resource2 :=  lrsms_util.NewResource(r2ID, r2_depended_list, &curTime,
    []byte("BodyofR1BodyofR2"))
  myDevice1.CreateResource(a1ID, resource2)

  r2_depended_list = list.New()
  resource2 =  lrsms_util.NewResource(r2ID, r2_depended_list, &curTime,
    []byte("BodyofR1BodyofR2"))
  myDevice1.CreateResource(a2ID, resource2)


  r3ID := localhost+device1AppServPort+"/"+a2ID+"/"+"Resource3"
  r3_depended_list := list.New()
  curTime = time.Now()
  resource3 :=  lrsms_util.NewResource(r3ID, r3_depended_list, &curTime,
    []byte("BodyofR3"))
  myDevice1.CreateResource(a2ID, resource3)

  r4ID := localhost+device1AppServPort+"/"+a2ID+"/"+"Resource4"
  r4_depended_list := list.New()
  r4_depended_list.PushBack(resource3.URI)
  r4_depended_list.PushBack(resource2.URI)
  curTime = time.Now()
  resource4 :=  lrsms_util.NewResource(r4ID, r4_depended_list, &curTime,
    []byte("BodyofR4BodyofR3"))
  myDevice1.CreateResource(a2ID, resource4)

  //****************************************************************************
  //2. updates
  //****************************************************************************
  resource1.Content = []byte("content update at: "+string(time.Now().String()))
  myDevice1.UpdateResource(a1ID, resource1.URI, resource1.Content)

  resource3.Content = []byte("content update at: "+string(time.Now().String()))
  myDevice1.UpdateResource(a2ID, resource3.URI,   resource3.Content)

  //****************************************************************************
  //3. sync
  //Create a device with one app: app03
  //app03 has two resources: Resource2 and Resource5.
  //Resource5 is depending on Resource2
  //****************************************************************************
  //create another device, apps, resources and connected Device
  myDevice2 := lrsms_util.StartDevice(device2AppServPort,device2LRSNSServPort)
  myDevice2Ele := UnConnectedDevices.PushBack(myDevice2)

  myDevice2.AddApp(a3ID)

  myDevice2.CreateResource(a3ID, resource2)

  r5ID := localhost+device2AppServPort+"/"+a3ID+"/"+"Resource5"
  r5_depended_list := list.New()  //localhost:5683/app02/Resource4
  r5_depended_list.PushBack(resource2.URI)
  curTime = time.Now()
  resource5 :=  lrsms_util.NewResource(r5ID, r5_depended_list, &curTime,
    []byte("BodyofR5BodyofR2"))
  myDevice2.CreateResource(a3ID, resource5)


   //make device 1 connected
   UnConnectedDevices.Remove(myDevice1Ele)
   ConnectedDevices.PushBack(myDevice1)
   myDevice1.Connect(ConnectedDevices)

   //make device 2 connected and trigger sync
   UnConnectedDevices.Remove(myDevice2Ele)
   ConnectedDevices.PushBack(myDevice2)
   myDevice2.Connect(ConnectedDevices)

  //****************************************************************************
  //4. cross device updates
  //****************************************************************************
  resource1.Content = []byte("content update at: "+string(time.Now().String()))
  myDevice1.UpdateResource(a1ID, resource1.URI, resource1.Content)

  //****************************************************************************
  //5. Disconnect device 2 update resource 1 than reconnect.  check sync
  //****************************************************************************


  //****************************************************************************
  //6.  delete resurce
  //****************************************************************************
  //myDevice1.DeleteResource(a2ID, resource4.URI)


  //******************************************************************************
  //Halt
  //******************************************************************************
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
