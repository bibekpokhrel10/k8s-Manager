package joomla

import (
	"context"
	"os"

	"k8smanager/internal"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (oc *Joomla) Delete(wname string) error {
	clientset := internal.GetConfig()
	deploymentsClient := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE"))
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

	servicesClinet := clientset.CoreV1().Services(os.Getenv("NAMESPACE"))
	err = servicesClinet.Delete(context.Background(), wname+"-mysql", metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}
	err = servicesClinet.Delete(context.Background(), wname, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(os.Getenv("NAMESPACE"))
	err = pvcClient.Delete(context.Background(), wname+"-jo-pv-claim", metav1.DeleteOptions{})
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
	err = namespaceClient.Delete(context.Background(), os.Getenv("NAMESPACE"), metav1.DeleteOptions{})
	if err != nil {
		log.Error(err)
		getErr = err
	}

	if getErr != nil {
		return getErr
	}
	return nil
}
