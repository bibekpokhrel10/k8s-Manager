package wordpress

import (
	"k8smanager/internal/clientgo"
)

type WordPress struct{}

func NewWordpressApp() clientgo.AppInterface {
	w := &WordPress{}
	return w
}
