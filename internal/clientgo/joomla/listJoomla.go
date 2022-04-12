package joomla

import (
	"context"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (jo *Joomla) List(name string) (clientgo.ListNames, error) {
	Names := clientgo.ListNames{}
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("joomla", name)
	servicelist, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	var getErr error
	if err != nil {
		log.Error(err)
		getErr = err
	}

	for _, service := range servicelist.Items {
		Names.Service = append(Names.Service, service.Name)
	}

	deploymentlist, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	for _, deploy := range deploymentlist.Items {
		Names.Deployment = append(Names.Deployment, deploy.Name)
	}

	podlist, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	for _, pod := range podlist.Items {
		Names.Pod = append(Names.Pod, pod.Name)
	}

	pvclist, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}
	for _, pvc := range pvclist.Items {
		Names.Pvc = append(Names.Pvc, pvc.Name)
	}
	if getErr != nil {
		return Names, getErr
	}
	return Names, nil
}
