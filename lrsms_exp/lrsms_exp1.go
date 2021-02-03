package main

import (
  "lrsms"
  "lrsms_util"
  "container/list"
  "time"
  "sync"
  "fmt"
)

//Create {Resource 3}---Depend on---->{Resource 2}---Depend on ---->{Resource 1}
func main() {
  //create lrsms
  var mylrsms lrsms.LRSMS
  mylrsms.ResourceRefs = make(map[string]*lrsms.ResourceRef)
  mylrsms.Mutex = &sync.Mutex{}

  //create resource_pool
  var myResourcePool lrsms_util.ResourcePool
  myResourcePool = make(map[string]*lrsms_util.Resource)

//******************************************************************************
// 1. Create Resources and their dependency
//******************************************************************************
  fmt.Println("Create Resources and dependency")

  //create resource 1
  r1_depended_list := list.New()
  curTime := time.Now()
  resource1 :=  lrsms_util.NewResource("Resource1", r1_depended_list, &curTime,
    []byte("BodyofR1"),myResourcePool)
  myResourcePool[resource1.URI] = resource1
  //create resource 1 ref in mylrsms
  mylrsms.CreateRef(resource1.URI, resource1.Depended,
    resource1.CreateTime, resource1.Get, resource1.Update, resource1.Alert)

  //create resource 2
  r2_depended_list := list.New()
  r2_depended_list.PushBack("Resource1")
  curTime = time.Now()
  resource2 :=  lrsms_util.NewResource("Resource2", r2_depended_list, &curTime,
    []byte("BodyofR1BodyofR2"),myResourcePool)
  myResourcePool[resource2.URI] = resource2
  //create resource 2 ref in mylrsms
  mylrsms.CreateRef(resource2.URI, resource2.Depended,
    resource2.CreateTime, resource2.Get, resource2.Update, resource2.Alert)

  //create resource 3
  r3_depended_list := list.New()
  r3_depended_list.PushBack("Resource2")
  curTime = time.Now()
  resource3 :=  lrsms_util.NewResource("Resource3", r3_depended_list, &curTime,
    []byte("BodyofR1BodyofR2BodyofR3"),myResourcePool)
  myResourcePool[resource3.URI] = resource3
  //create resource 3 ref in mylrsms
  mylrsms.CreateRef(resource3.URI, resource3.Depended,
    resource3.CreateTime, resource3.Get, resource3.Update, resource3.Alert)

  mylrsms.ResourceRefs["Resource1"].Print()
  mylrsms.ResourceRefs["Resource2"].Print()
  mylrsms.ResourceRefs["Resource3"].Print()

//******************************************************************************
//2. Simple propergation update
//******************************************************************************
  resource1.Content =[]byte("newR1Bodyabc")
  mylrsms.RecieveUpdateFromInside(resource1.URI)
  mylrsms.ResourceRefs["Resource1"].Print()
  mylrsms.ResourceRefs["Resource2"].Print()
  mylrsms.ResourceRefs["Resource3"].Print()
}
