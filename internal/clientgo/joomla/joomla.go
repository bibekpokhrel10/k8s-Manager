package joomla

import (
	"k8smanager/internal/clientgo"
)

type Joomla struct{}

func NewJoomlaApp() clientgo.AppInterface {
	jo := &Joomla{}
	return jo
}
