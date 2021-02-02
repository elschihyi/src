package main

import (
  "lrsms"
  "lrsms_util"
  "container/list"
  "time"
)

//Create {Resource 3}---Depend on---->{Resource 2}---Depend on ---->{Resource 1}

func main() {
  //create resource 1
  r1_depended_list := list.New()
  curTime := time.Now()
  resource1 :=  lrsms_util.NewResource("Resource1", r1_depended_list, &curTime,
    []byte("BodyofR1"))
  //create resource 1 reference
	resource_ref_1 := lrsms.NewRF(resource1.URI, resource1.Depended,
    resource1.CreateTime, resource1.Get,  resource1.Update, resource1.Alert)

  //create resource 2
  r2_depended_list := list.New()
  r2_depended_list.PushBack("Resource1")
  resource_ref_1.Dependent.PushBack("Resource2")
  curTime = time.Now()
  resource2 :=  lrsms_util.NewResource("Resource2", r2_depended_list, &curTime,
    []byte("BodyofR2"))
  //create resource 2 reference
  resource_ref_2 := lrsms.NewRF(resource2.URI, resource2.Depended,
    resource2.CreateTime, resource2.Get,  resource2.Update, resource2.Alert)

  //create resource 3
  r3_depended_list := list.New()
  r3_depended_list.PushBack("Resource2")
  resource_ref_2.Dependent.PushBack("Resource3")
  curTime = time.Now()
  resource3 :=  lrsms_util.NewResource("Resource3", r3_depended_list, &curTime,
    []byte("BodyofR3"))
  //create resource 3 reference
  resource_ref_3 := lrsms.NewRF(resource3.URI, resource3.Depended,
    resource3.CreateTime, resource3.Get,  resource3.Update, resource3.Alert)

  resource_ref_1.Print ()
  resource_ref_2.Print ()
  resource_ref_3.Print ()
}
