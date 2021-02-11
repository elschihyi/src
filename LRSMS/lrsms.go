//Local resource state manage system
package lrsms

import (
	"sync"
	"container/list"
  "time"
	"log"
)
//******************************************************************************
// Types Definition
//******************************************************************************
type GetDeviceRefs func(string, string)
type SentRefs func(string, string, map[string]string)
type GetResource func(string, string, string)
type SentResource func(string, string, string, string, string)
type UpdateCache func(string, string, string)
type SentDeviceRef func(string, string, string, string)

type LRSMS struct{
	ResourceRefs map[string]*ResourceRef
	ConnectedDevRes map[string]*list.List //[host:port]resourceIDList
	Mutex *sync.Mutex
	GetDeviceRefsfunc GetDeviceRefs
	SentRefsfunc SentRefs
	GetResourcefunc GetResource
	SentResourcefunc SentResource
	UpdateCachefunc UpdateCache
	SentDeviceReffunc SentDeviceRef
}

//******************************************************************************
//Public Functions
//******************************************************************************
func NewLRSMS (getDeviceRefsfunc GetDeviceRefs, sentRefsfunc SentRefs,
	getResourcefunc GetResource, sentResourcefunc SentResource,
	updateCachefunc UpdateCache, sentDeviceReffunc SentDeviceRef)*LRSMS{
	var newLRSMS LRSMS
  newLRSMS.ResourceRefs = make(map[string]*ResourceRef)
	newLRSMS.ConnectedDevRes = make (map[string]*list.List)
  newLRSMS.Mutex = &sync.Mutex{}
	newLRSMS.GetDeviceRefsfunc = getDeviceRefsfunc
	newLRSMS.SentRefsfunc = sentRefsfunc
	newLRSMS.GetResourcefunc = getResourcefunc
	newLRSMS.SentResourcefunc = sentResourcefunc
	newLRSMS.UpdateCachefunc = updateCachefunc
	newLRSMS.SentDeviceReffunc = sentDeviceReffunc
	return &newLRSMS
}

//******************************************************************************
//Public Functions Res
//******************************************************************************
func (lrsms *LRSMS) GetResource(otherAdd string, hostAdd string,
	resourceURI string){
  //1. check if  resourceURI in lrsms.resourceRefs
  if _, exists := lrsms.ResourceRefs[resourceURI]; exists{
		//log.Printf("get %v in %v", resourceURI, hostAdd)
		content := lrsms.ResourceRefs[resourceURI].Getfunc(resourceURI)
		lrsms.SentResourcefunc(otherAdd, hostAdd, resourceURI,
			lrsms.ResourceRefs[resourceURI].CreateTime, string(content))
		return
	}
	//2. check if resourceURI is a depended resource of a Cache(?????????????????)
	for k, v := range lrsms.ResourceRefs {
		if lrsms.IsCache(k) && v.Depended.Front().Value.(string) == resourceURI {
			//log.Printf("get %v in %v at cache %v", resourceURI, hostAdd, k)
			content := lrsms.ResourceRefs[k].Getfunc(k)
			lrsms.SentResourcefunc(otherAdd, hostAdd, resourceURI,
				lrsms.ResourceRefs[k].CreateTime, string(content))
			return
		}
	}
}

func (lrsms *LRSMS) UpdateResource(otherAdd string, hostAdd string,
	resourceURI string, createTime string, Content string){
	for k, v := range lrsms.ResourceRefs {
		if lrsms.IsCache(k) && v.Depended.Front().Value.(string) == resourceURI {
			 //log.Printf("UpdateResource %v cache %v", resourceURI, k)
			 lrsms.UpdateCachefunc(k, createTime, Content)
	  }
	}
}

//******************************************************************************
//Public Functions Ref
//******************************************************************************

func (lrsms *LRSMS) CreateRef(uri string, hostAdd string, depended *list.List,
	 createTime string, getFunc Get, updateFunc Update, alertFunc Alert){
  //if resource info already exit
	if _, exist := lrsms.ResourceRefs[uri]; exist {
	  return
  }
	new_resource_ref := NewRF(uri, depended, createTime, getFunc, updateFunc,
		 alertFunc)
  lrsms.Mutex.Lock()
  lrsms.ResourceRefs[uri] = new_resource_ref
	for e := depended.Front(); e != nil; e = e.Next() {
		if _, exist := lrsms.ResourceRefs[e.Value.(string)]; exist {
		  lrsms.ResourceRefs[e.Value.(string)].Dependent.PushBack(uri)
		}
	}
	lrsms.Mutex.Unlock()
	//log.Printf("LRSMS CreateRef: %v", uri)

	//announce public
	for k, _ := range lrsms.ConnectedDevRes{
		if lrsms.IsCache(uri) {
			lrsms.SentDeviceReffunc(k, hostAdd,
				lrsms.ResourceRefs[uri].Depended.Front().Value.(string),
				lrsms.ResourceRefs[uri].CreateTime)
  	} else {
			lrsms.SentDeviceReffunc(k, hostAdd, uri, lrsms.ResourceRefs[uri].CreateTime)
		}
	}
}

func (lrsms *LRSMS) RecieveUpdateFromInside(uri string, hostAdd string,
	newTime string){
	//log.Printf("RecieveUpdateFromInside %v time %v", uri, newTime)
	lrsms.Mutex.Lock()
	//update Ref create time
  lrsms.ResourceRefs[uri].UpdateCreateTime(newTime)
	// set Ref unflag if all its depended resource are unflag
	ResourceUnflag := true
	for e := lrsms.ResourceRefs[uri].Depended.Front(); e != nil; e = e.Next() {
		if _, exist := lrsms.ResourceRefs[e.Value.(string)]; exist{
		  if lrsms.ResourceRefs[e.Value.(string)].Flag {
			  ResourceUnflag = false
			  break
		  }
	  }
	}
	if ResourceUnflag {
		lrsms.ResourceRefs[uri].Unflag()
	}
 	lrsms.Mutex.Unlock()

	//****************************************************************************
	//***dont need to flag self***************************************************
	//lrsms.Mutex.Lock()
	//lrsms.ResourceRefs[uri].Flag=true
	//lrsms.Mutex.Unlock()
	//****************************************************************************
	//signal all dependent to flag
	for e := lrsms.ResourceRefs[uri].Dependent.Front(); e != nil; e = e.Next() {
		lrsms.FlagPluxPropergate(e.Value.(string))
	}
	//****************************************************************************
	//***dont need to update self*************************************************
	//****************************************************************************
	//signal all dependent resource check update
	for e := lrsms.ResourceRefs[uri].Dependent.Front(); e != nil; e = e.Next() {
		lrsms.CheckUpDate(e.Value.(string))
	}

	//announce public
	for k, v := range lrsms.ConnectedDevRes{
		for e := v.Front(); e != nil; e = e.Next(){
		  if	e.Value.(string) == uri{
		 	  lrsms.SentDeviceReffunc(k, hostAdd, uri,
			 	  lrsms.ResourceRefs[uri].CreateTime)
		  }
		}
	}
}

func (lrsms *LRSMS) GetRefs(otherAdd string, hostAdd string){
	simpleResourceRefs := make(map[string]string) //map[resourceID]CreateTime
	i := 0
  for k, v := range lrsms.ResourceRefs {
		if lrsms.IsCache(k){ //if resource is a cache use dependedResource ID as URI
			simpleResourceRefs[v.Depended.Front().Value.(string)] = v.CreateTime
			continue
		}
		simpleResourceRefs[v.URI] = v.CreateTime
		i++
	}
	lrsms.SentRefsfunc(otherAdd, hostAdd, simpleResourceRefs)
}
//******************************************************************************
//Public Functions Dev
//******************************************************************************
func (lrsms *LRSMS) NewDevice(otherAdd string, hostAdd string){
  //create new Device and request that device Resource
	lrsms.Mutex.Lock()
	lrsms.ConnectedDevRes[otherAdd] = list.New()
	lrsms.Mutex.Unlock()
	//log.Printf("add new device %v", newDevice)
	lrsms.GetDeviceRefsfunc(otherAdd, hostAdd)
}

func (lrsms *LRSMS)UpdateOtherDeviceRes(otherAdd string, hostAdd string,
	resourceRefs map[string]interface{}){
	//log.Printf("resourceRefs length %v", len(resourceRefs))
  //1. add other resourceID to ConnectedDevRes[otherAdd]
	lrsms.Mutex.Lock()
	for k, _ := range resourceRefs {
    kInConnectedDevRes :=false
		for e := lrsms.ConnectedDevRes[otherAdd].Front(); e != nil; e = e.Next() {
			if e.Value.(string) == k {
				kInConnectedDevRes = true
				break
			}
		}
		if !kInConnectedDevRes {
    	lrsms.ConnectedDevRes[otherAdd].PushBack(k)
		}
		//log.Printf("push %v",k)
	}
	lrsms.Mutex.Unlock()

	//2. compare create time and get updat for those outdated resources
	for k, v := range resourceRefs { //check all resource in other device
		for k2, v2 := range lrsms.ResourceRefs {
     if lrsms.IsCache(k2) && v2.Depended.Front().Value.(string) == k{
			 //log.Printf("resourceRefs length %v", len(resourceRefs))
			 localRefTime, _ := time.Parse(time.RFC3339, v2.CreateTime)
			 otherRefTime, _ := time.Parse(time.RFC3339, v.(string))
			 if otherRefTime.After(localRefTime){
				 lrsms.GetResourcefunc(otherAdd, hostAdd, k)
			 }
		 }
		}
	}
}

/*
func (lrsms *LRSMS) DeleteRef(uri string){
	if lrsms.ResourceRefs[uri].Dependent.Len() == 0 {
		lrsms.Mutex.Lock()
		delete(lrsms.ResourceRefs,uri)
		lrsms.Mutex.Unlock()
	}
}
*/

//******************************************************************************
//Public Functions
//******************************************************************************

func (lrsms *LRSMS) CheckUpDate(uri string){
	//log.Printf("CheckUpDate %v", uri)
  //check if all depended resource are unflag
	if lrsms.ResourceRefs[uri].Depended.Len() ==0{
		return
	}
	for e := lrsms.ResourceRefs[uri].Depended.Front(); e != nil; e = e.Next() {
		if lrsms.ResourceRefs[e.Value.(string)].Flag {
			return
		}
	}
	//update resource
	lrsms.ResourceRefs[uri].Updatefunc(uri)
}

func (lrsms *LRSMS) FlagPluxPropergate(uri string){
  //flag resourceRef
	lrsms.Mutex.Lock()
	lrsms.ResourceRefs[uri].Flagg()
	lrsms.Mutex.Unlock()
	//Alert Resource
	lrsms.ResourceRefs[uri].Alertfunc(uri)
	//propergate up
	for e := lrsms.ResourceRefs[uri].Dependent.Front(); e != nil; e = e.Next() {
		lrsms.FlagPluxPropergate(e.Value.(string))
	}
}

func (lrsms *LRSMS) Print(){
	for _, v := range lrsms.ResourceRefs {
		v.Print()
	}
}

//******************************************************************************
//Public Untility Functions
//******************************************************************************
func (lrsms *LRSMS) IsCache(resourceID string)bool{
	resouceRef := lrsms.ResourceRefs[resourceID]
	if resouceRef.Depended.Len() == 1 { //one depended
		dependedResourceID := resouceRef.Depended.Front().Value.(string)
		if _, exist := lrsms.ResourceRefs[dependedResourceID]; !exist{ //not in resourceRefs
			return true
		}
	}
	return false
}
