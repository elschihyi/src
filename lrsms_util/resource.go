package lrsms_util

import (
  "container/list"
  "time"
  "fmt"
  "strconv"
  //"log"
)

//******************************************************************************
// Types Definition
//******************************************************************************

type Resource struct{
  URI string //identifier(URI)
  Content []byte //Resource Content
  Depended *list.List //depended resource reference list
  CreateTime string //ResourceCache Created time in RFC3339 format
  Flag bool //True: some dependent resource is not ptodate
}

//******************************************************************************
//Public Static Functions
//******************************************************************************
func NewResource(uri string, depended *list.List,
  createTime string, content []byte)*Resource{
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
  //time.Sleep(2 * time.Second)

  // run the resource update
  r.Content =[]byte(r.URI + " update at: "+time.Now().Format(time.RFC3339))
  r.CreateTime = time.Now().Format(time.RFC3339)
  //set flag to false
  r.Flag = false
  //log.Printf("Resource %v finished update at %v", r.URI, r.CreateTime)
}

func(r *Resource) Alert(){
  r.Flag = true
  //log.Printf("Resource %v recieved update alert at %v", r.URI, time.Now().Format(time.RFC3339))
}

func (r *Resource) Print (){
  fmt.Println(r.URI)
  fmt.Print("  depended:")
  for e := r.Depended.Front(); e != nil; e = e.Next() {
    fmt.Print(" "+e.Value.(string))
  }
  fmt.Println("")
  fmt.Println("  CreateTime: "+r.CreateTime)
  fmt.Println("  Cache: "+string(r.Content))
  fmt.Println("  Flag: "+strconv.FormatBool(r.Flag))
}
