package wordpress

import (
	"wordpress.com/internal/clientgo"
)

type WordPress struct {
	Name string `json:"name"`
}

func NewWordPress(wp *WordPress) clientgo.AppInterface {
	if wp == nil {
		return nil
	}
	return wp
}

func NewWordpressApp(wname string) clientgo.AppInterface {

	w := NewWordPress(&WordPress{Name: wname})

	return w
}
