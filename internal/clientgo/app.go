package clientgo

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type AppInterface interface {
	List(string) (ListNames, error)
	Detail(string) (*appsv1.Deployment, *v1.Service, error)
	Create(string, int32) error
	Update(int32, string, string, string) error
	Delete(string) error
}
