package lrsms_util

import (
  "container/list"
  "time"
  //"math/rand"
  "fmt"
  "strconv"
  "log"
)

//******************************************************************************
// Types Definition
//******************************************************************************

type Resource struct{
  URI string //identifier(URI)
  Content []byte //Resource Content
  Depended *list.List //depended resource reference list
  CreateTime *time.Time //ResourceCache Created time
  Flag bool //True: some dependent resource is not ptodate
}

//******************************************************************************
//Public Static Functions
//******************************************************************************
func NewResource(uri string, depended *list.List, createTime *time.Time,
  content []byte)*Resource{
  var newR Resource
  newR.URI = uri
  newR.Content = content
  newR.Depended =  depended
  newR.CreateTime = createTime
  newR.Flag = false
  return &newR
}

//******************************************************************************
//public functions
//******************************************************************************
//Get Resource Content
func(r *Resource) Get() []byte{
  return r.Content
}

//make Resource update its info
func(r *Resource) Update(){
   // run the resource update
   r.Content =[]byte(r.URI + " update at: "+time.Now().Format(time.RFC3339))
   //set flag to false
   r.Flag = false
   log.Printf("Resource %v finished update, content is %v", r.URI, string(r.Content))
}

func(r *Resource) Alert(){
  r.Flag = true
  log.Printf("Resource %v recieved update alert ", r.URI)
}

func (r *Resource) Print (){
  fmt.Println(r.URI)
  fmt.Print("  depended:")
  for e := r.Depended.Front(); e != nil; e = e.Next() {
    fmt.Print(" "+e.Value.(string))
  }
  fmt.Println("")
  fmt.Println("  CreateTime: "+r.CreateTime.String())
  fmt.Println("  Cache: "+string(r.Content))
  fmt.Println("  Flag: "+strconv.FormatBool(r.Flag))
}
