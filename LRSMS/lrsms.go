//Local resource state manage system
package lrsms

import (
	"sync"
	"container/list"
  "time"
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
func (lrsms *LRSMS) CreateRef(uri string, depended *list.List,
	 createTime *time.Time, getFunc Get, updateFunc Update, alertFunc Alert){
	new_resource_ref := NewRF(uri, depended, createTime, getFunc,
	 updateFunc, alertFunc)
  lrsms.Mutex.Lock()
  lrsms.ResourceRefs[uri] = new_resource_ref
	for e := depended.Front(); e != nil; e = e.Next() {
		lrsms.ResourceRefs[e.Value.(string)].Dependent.PushBack(uri)
	}
	lrsms.Mutex.Unlock()
}
func (lrsms *LRSMS) GetRef(uri string)*ResourceRef{
	return lrsms.ResourceRefs[uri]
}
func (lrsms *LRSMS) GetResource(uri string)[]byte{
	return lrsms.ResourceRefs[uri].Getfunc()
}
func (lrsms *LRSMS) DeleteRef(uri string){
	if lrsms.ResourceRefs[uri].Dependent.Len() == 0 {
		lrsms.Mutex.Lock()
		delete(lrsms.ResourceRefs,uri)
		lrsms.Mutex.Unlock()
	}
}
func (lrsms *LRSMS) Update(){
}
