//Local resource state manage system
package lrsms

import (
	"fmt"
	"errors"
	"sync"
)
//******************************************************************************
// Types Definition
//******************************************************************************
type LRSMS struct{
	resourceRefs map[string]*ResourceRef
	mutex *sync.Mutex
}

//******************************************************************************
//Public Functions
//******************************************************************************
func (lrsms *LRSMS) CreateRef(){
	new_resource_ref := lrsms.NewRF(resource1.URI, resource1.Depended,
    resource1.CreateTime, resource1.Get,  resource1.Update, resource1.Alert)
}
func (lrsms *LRSMS) GetRef(){
}
func (lrsms *LRSMS) Update(){
}
func (lrsms *LRSMS) GerResource(){
}
func (lrsms *LRSMS) DeleteRef(){
}
