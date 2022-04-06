package wordpress

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wordpress.com/internal"
)

func (wp *WordPress) List() ([]string, []string, []string, []string) {
	clientset := internal.GetConfig()
	servicelist, err := clientset.CoreV1().Services("bibek").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil
	}
	var services []string
	for _, service := range servicelist.Items {
		services = append(services, service.Name)
	}

	deploymentlist, err := clientset.AppsV1().Deployments("bibek").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil
	}
	var deployments []string
	for _, deploy := range deploymentlist.Items {
		deployments = append(deployments, deploy.Name)
	}

	podlist, err := clientset.CoreV1().Pods("bibek").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil
	}
	var pods []string
	for _, pod := range podlist.Items {
		pods = append(pods, pod.Name)
	}

	var pvcs []string
	pvclist, err := clientset.CoreV1().PersistentVolumeClaims("bibek").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil
	}
	for _, pvc := range pvclist.Items {
		pvcs = append(pvcs, pvc.Name)
	}
	return services, deployments, pods, pvcs
}
