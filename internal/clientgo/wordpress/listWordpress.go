package wordpress

import (
	"context"
	"os"

	"k8smanager/internal"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wp *WordPress) List() ([]string, []string, []string, []string, error) {
	clientset := internal.GetConfig()
	servicelist, err := clientset.CoreV1().Services(os.Getenv("NAMESPACE")).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, err
	}
	var services []string
	for _, service := range servicelist.Items {
		services = append(services, service.Name)
	}

	deploymentlist, err := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE")).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, err
	}
	var deployments []string
	for _, deploy := range deploymentlist.Items {
		deployments = append(deployments, deploy.Name)
	}

	podlist, err := clientset.CoreV1().Pods(os.Getenv("NAMESPACE")).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, err
	}
	var pods []string
	for _, pod := range podlist.Items {
		pods = append(pods, pod.Name)
	}

	var pvcs []string
	pvclist, err := clientset.CoreV1().PersistentVolumeClaims(os.Getenv("NAMESPACE")).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, err
	}
	for _, pvc := range pvclist.Items {
		pvcs = append(pvcs, pvc.Name)
	}
	return services, deployments, pods, pvcs, nil
}
