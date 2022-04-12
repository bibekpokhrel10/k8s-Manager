package wordpress

import (
	"context"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wp *WordPress) Delete(wname string) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", wname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	err := deploymentsClient.Delete(context.Background(), wname, metav1.DeleteOptions{})
	var getErr error
	if err != nil {
		log.Error(err)
		getErr = err
	}
	err = deploymentsClient.Delete(context.Background(), wname+"-mysql", metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	servicesClinet := clientset.CoreV1().Services(namespace)
	err = servicesClinet.Delete(context.Background(), wname+"-mysql", metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}
	err = servicesClinet.Delete(context.Background(), wname, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(namespace)
	err = pvcClient.Delete(context.Background(), wname+"-pv-claim", metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}
	err = pvcClient.Delete(context.Background(), wname+"-mysql-pv-claim", metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	namespaceClient := clientset.CoreV1().Namespaces()
	err = namespaceClient.Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	if getErr != nil {
		return getErr
	}
	return nil
}
