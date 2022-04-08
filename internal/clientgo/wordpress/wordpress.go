package wordpress

import (
	"k8smanager/internal/clientgo"
)

type WordPress struct{}

func NewWordPress(wp *WordPress) clientgo.AppInterface {
	if wp == nil {
		return nil
	}
	return wp
}

func NewWordpressApp() clientgo.AppInterface {
	w := NewWordPress(&WordPress{})
	return w
}
