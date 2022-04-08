package joomla

import (
	"context"
	"os"

	"k8smanager/internal"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (oc *Joomla) Detail(pname string, object string) (*appsv1.Deployment, *v1.Service, error) {
	if object == "deployment" {
		deploymentdetails, err := DeploymentDetail(pname)
		return deploymentdetails, nil, err
	}
	if object == "service" {
		servicedetails, err := ServiceDetail(pname)
		return nil, servicedetails, err
	}
	return nil, nil, nil
}

func DeploymentDetail(pname string) (*appsv1.Deployment, error) {
	clientset := internal.GetConfig()
	GetDeployment, err := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE")).Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetDeployment, nil
}

func ServiceDetail(pname string) (*v1.Service, error) {
	clientset := internal.GetConfig()
	GetService, err := clientset.CoreV1().Services(os.Getenv("NAMESPACE")).Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetService, nil
}
