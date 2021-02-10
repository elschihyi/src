//Local resource state manage system
package lrsms

import (
	"sync"
	"container/list"
  "time"
	//"log"
)
//******************************************************************************
// Types Definition
//******************************************************************************
type LRSMS struct{
	ResourceRefs map[string]*ResourceRef
	Mutex *sync.Mutex
}

//******************************************************************************
//Public Functions
//******************************************************************************
func NewLRSMS ()*LRSMS{
	var newLRSMS LRSMS
  newLRSMS.ResourceRefs = make(map[string]*ResourceRef)
  newLRSMS.Mutex = &sync.Mutex{}
	return &newLRSMS
}

func (lrsms *LRSMS) CreateRef(uri string, depended *list.List,
	 createTime *time.Time, getFunc Get, updateFunc Update, alertFunc Alert){
  //do nothing if resource info already exit
	if _, exist := lrsms.ResourceRefs[uri]; exist {
    return
  }

	new_resource_ref := NewRF(uri, depended, createTime, getFunc, updateFunc,
		 alertFunc)
  lrsms.Mutex.Lock()
  lrsms.ResourceRefs[uri] = new_resource_ref
	for e := depended.Front(); e != nil; e = e.Next() {
		//log.Printf("LRSMS CreateRef depended : %v", e.Value.(string))
		lrsms.ResourceRefs[e.Value.(string)].Dependent.PushBack(uri)
	}
	lrsms.Mutex.Unlock()
	//log.Printf("LRSMS CreateRef: %v", uri)
}

func (lrsms *LRSMS) GetResource(uri string)[]byte{
	return lrsms.ResourceRefs[uri].Getfunc(uri)
}

func (lrsms *LRSMS) DeleteRef(uri string){
	if lrsms.ResourceRefs[uri].Dependent.Len() == 0 {
		lrsms.Mutex.Lock()
		delete(lrsms.ResourceRefs,uri)
		lrsms.Mutex.Unlock()
	}
}

func (lrsms *LRSMS) RecieveUpdateFromInside(uri string, newTime *time.Time){
   lrsms.ResourceRefs[uri].CreateTime = newTime
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
}

func (lrsms *LRSMS) RecieveUpdateFromOutside(uri string){
}

func (lrsms *LRSMS) GetRef(uri string)*ResourceRef{
	return lrsms.ResourceRefs[uri]
}

//******************************************************************************
//Public Functions (Operations)
//******************************************************************************

func (lrsms *LRSMS) CheckUpDate(uri string){
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
	//unflag
	lrsms.Mutex.Lock()
	lrsms.ResourceRefs[uri].Flag = false
	lrsms.Mutex.Unlock()
	//signal all dependent resource check update
	for e := lrsms.ResourceRefs[uri].Dependent.Front(); e != nil; e = e.Next() {
		lrsms.CheckUpDate(e.Value.(string))
	}
}

func (lrsms *LRSMS) FlagPluxPropergate(uri string){
  //flag resourceRef
	lrsms.Mutex.Lock()
	lrsms.ResourceRefs[uri].Flag = true
	lrsms.Mutex.Unlock()
	//Alert Resource
	lrsms.ResourceRefs[uri].Alertfunc(uri)
	//propergate up
	for e := lrsms.ResourceRefs[uri].Dependent.Front(); e != nil; e = e.Next() {
		lrsms.FlagPluxPropergate(e.Value.(string))
	}
}

func (lrsms *LRSMS) Unflag(uri string){
	/*
	lrsms.Mutex.Lock()
	lrsms.ResourceRefs[uri].Flag = false
	lrsms.Mutex.Unlock()
	*/
}
