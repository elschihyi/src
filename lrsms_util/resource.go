package lrsms_util

import (
  "container/list"
  "time"
  //"math/rand"
  "fmt"
  "strconv"
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
  thePool ResourcePool //this is for ourself@@
}

//******************************************************************************
//Public Static Functions
//******************************************************************************
func NewResource(uri string, depended *list.List, createTime *time.Time,
  content []byte, thePool ResourcePool)*Resource{
  var newR Resource
  newR.URI = uri
  newR.Content = content
  newR.Depended =  depended
  newR.CreateTime = createTime
  newR.Flag = false
  newR.thePool = thePool
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
   r.Content =[]byte(r.URI)
   for e := r.Depended.Front(); e != nil; e = e.Next() {
     r.Content = append(r.Content, r.thePool[e.Value.(string)].Content...)
   }
   //update will run ramdom from 0 to 2 sec in our experiment
   //time.Sleep(time.Second*time.Duration(rand.Intn(2)))
   r.Flag = false
}

func(r *Resource) Alert(){
  r.Flag = true
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
