package lrsms

import (
  "container/list"
  //"time"
  "fmt"
  "strconv"
  //"log"
)

//******************************************************************************
// Types Definition
//******************************************************************************

type Get func(string)[]byte
type Update func(string)
type Alert func(string)

type ResourceRef struct{
  URI string //identifier(URI)
  Depended *list.List  //depended resource reference list
  Dependent *list.List //dependent resource reference list
  CreateTime string //ResourceCache Created time
  Flag bool //True: some dependent resource is not ptodate
  Getfunc Get //Get Resouce Cache
  Updatefunc Update //Update Resource
  Alertfunc Alert //Alert Resource
}

//******************************************************************************
//Public Static Functions
//******************************************************************************
func NewRF(uri string, depended *list.List, createTime string,
  getFunc Get, updateFunc Update, alertFunc Alert) *ResourceRef{
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
func (rf *ResourceRef)UpdateCreateTime(newTime string){
  rf.CreateTime = newTime
  //log.Printf("ResourceRef %v changed update time at %v", rf.URI, rf.CreateTime)
}

func (rf *ResourceRef)Flagg(){
  rf.Flag = true
  //log.Printf("ResourceRef %v flag", rf.URI)
}

func (rf *ResourceRef)Unflag(){
  rf.Flag = false
  //log.Printf("ResourceRef %v unflag", rf.URI)
}

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
  fmt.Println("  CreateTime: "+rf.CreateTime)
  //fmt.Println("  Cache: "+string(rf.Getfunc(rf.URI)))
  fmt.Println("  Flag: "+strconv.FormatBool(rf.Flag))
}
