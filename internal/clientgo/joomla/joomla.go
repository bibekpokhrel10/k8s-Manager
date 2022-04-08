package joomla

import (
	"k8smanager/internal/clientgo"
)

type Joomla struct{}

func NewJoomla(oc *Joomla) clientgo.AppInterface {
	if oc == nil {
		return nil
	}
	return oc
}

func NewJoomlaApp() clientgo.AppInterface {
	oc := NewJoomla(&Joomla{})
	return oc
}
