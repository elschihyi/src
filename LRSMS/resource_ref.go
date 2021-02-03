package lrsms

import (
  "container/list"
  "time"
  "fmt"
  "strconv"
)

//******************************************************************************
// Types Definition
//******************************************************************************

type Get func()[]byte
type Update func()
type Alert func()

type ResourceRef struct{
  URI string //identifier(URI)
  Depended *list.List  //depended resource reference list
  Dependent *list.List //dependent resource reference list
  CreateTime *time.Time //ResourceCache Created time
  Getfunc Get //Get Resouce Cache
  Updatefunc Update //Update Resource
  Flag bool //True: some dependent resource is not ptodate
  Alertfunc Alert //Alert Resource
}

//******************************************************************************
//Public Static Functions
//******************************************************************************
func NewRF(uri string, depended *list.List, createTime *time.Time, getFunc Get,
  updateFunc Update, alertFunc Alert) *ResourceRef{
  var newRF ResourceRef
  newRF.URI = uri
  newRF.Depended = depended
  newRF.Dependent = list.New()
  newRF.CreateTime = createTime
  newRF.Getfunc = getFunc
  newRF.Updatefunc = updateFunc
  newRF.Alertfunc = alertFunc
  newRF.Flag = false
  return &newRF
}

//******************************************************************************
//Public Functions
//******************************************************************************
func (rf *ResourceRef) Print (){
  fmt.Println(rf.URI)
  fmt.Print("  depended:")
  for e := rf.Depended.Front(); e != nil; e = e.Next() {
    fmt.Print(" "+e.Value.(string))
  }
  fmt.Println("")
  fmt.Print("  dependent:")
  for e := rf.Dependent.Front(); e != nil; e = e.Next() {
    fmt.Print(" "+e.Value.(string))
  }
  fmt.Println("")
  fmt.Println("  CreateTime: "+rf.CreateTime.String())
  fmt.Println("  Cache: "+string(rf.Getfunc()))
  fmt.Println("  Flag: "+strconv.FormatBool(rf.Flag))
}
