package wordpress

import (
	"wordpress.com/internal/clientgo"
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
