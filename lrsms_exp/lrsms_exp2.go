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

func main() {
  //create 2 app: app1 in device 1 has resource 1 and 2 (2 depend on 1), app2 in device 2 has resoruce 1, 2, 3 and 4(4 depend on 3)
  //1. create resoruces
  r1_depended_list := list.New()
  curTime := time.Now()
  resource1 :=  lrsms_util.NewResource("Resource1", r1_depended_list, &curTime,
    []byte("BodyofR1"),nil)

  r2_depended_list := list.New()
  r2_depended_list.PushBack("Resource1")
  curTime = time.Now()
  resource2 :=  lrsms_util.NewResource("Resource2", r2_depended_list, &curTime,
    []byte("BodyofR1BodyofR2"),nil)

  r3_depended_list := list.New()
  curTime = time.Now()
  resource3 :=  lrsms_util.NewResource("Resource3", r3_depended_list, &curTime,
    []byte("BodyofR3"),nil)

  r4_depended_list := list.New()
  r4_depended_list.PushBack("Resource3")
  curTime = time.Now()
  resource4 :=  lrsms_util.NewResource("Resource4", r4_depended_list, &curTime,
    []byte("BodyofR4BodyofR3"),nil)

  r5_depended_list := list.New()
  r5_depended_list.PushBack("Resource3")
  curTime = time.Now()
  resource5 :=  lrsms_util.NewResource("Resource5", r5_depended_list, &curTime,
    []byte("BodyofR5BodyofR3"),nil)

  //2. create apps
  a1ID := "app01"
  app01resources := make(map[string]lrsms_util.Resource)
  app01resources[resource1.URI] = *resource1
  app01resources[resource2.URI] = *resource2

  a2ID := "app02"
  app02resources := make(map[string]lrsms_util.Resource)
  app02resources[resource1.URI] = *resource1
  app02resources[resource2.URI] = *resource2
  app02resources[resource3.URI] = *resource3
  app02resources[resource4.URI] = *resource4

  //3. create devices
  myDevice1 := lrsms_util.StartDevice(":5683",":5684")
  myDevice2 := lrsms_util.StartDevice(":5685",":5686")
  time.Sleep(2 * time.Second) //sleep so all 4 coap server can start running
  myDevice1.AddApp(a1ID, app01resources)
  myDevice2.AddApp(a2ID, app02resources)

  //4. updates.
  resource1.Content = []byte("content update at: "+string(time.Now().String()))
  myDevice1.UpdateResource(a1ID, resource1.URI, []byte("content update at: "+string(time.Now().String())))

  //5. delete
  //fmt.Println("before delete app a1 length is: %d", len(myDevice2.Apps[a2ID]))
  myDevice2.DeleteResource(a2ID, resource4.URI)
  //fmt.Println("after delete app a1 length is: %d", len(myDevice2.Apps[a2ID]))

  //6. add resource
  myDevice2.CreateResource(a2ID, resource5)

  //******************************************************************************
  //Halt
  //******************************************************************************
  //hault til user input e
  for input := "";input!="e";{
    fmt.Println("Enter 'e' to terminate")
    fmt.Scanf("%s", &input)
    //fmt.Println("Hello", name)
  }
}
