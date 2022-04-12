package wordpress

import (
	"context"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wp *WordPress) Detail(pname string) (*appsv1.Deployment, *v1.Service, error) {
	deploymentdetails, err := DeploymentDetail(pname)
	var getErr error
	if err != nil {
		getErr = err
	}
	servicedetails, err := ServiceDetail(pname)
	if err != nil {
		getErr = err
	}
	if getErr != nil {
		return nil, nil, getErr
	}
	return deploymentdetails, servicedetails, nil
}

func DeploymentDetail(pname string) (*appsv1.Deployment, error) {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)
	GetDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetDeployment, nil
}

func ServiceDetail(pname string) (*v1.Service, error) {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)
	GetService, err := clientset.CoreV1().Services(namespace).Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetService, nil
}
