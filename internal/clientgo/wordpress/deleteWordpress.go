package wordpress

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wordpress.com/internal"
)

func (wp *WordPress) Delete(wname string) error {
	clientset := internal.GetConfig()
	deploymentsClient := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE"))
	err := deploymentsClient.Delete(context.Background(), wname, metav1.DeleteOptions{})
	var getErr error
	if err != nil {
		getErr = err
	}
	err = deploymentsClient.Delete(context.Background(), wname+"-mysql", metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}

	servicesClinet := clientset.CoreV1().Services(os.Getenv("NAMESPACE"))
	err = servicesClinet.Delete(context.Background(), wname+"-mysql", metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}
	err = servicesClinet.Delete(context.Background(), wname, metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(os.Getenv("NAMESPACE"))
	err = pvcClient.Delete(context.Background(), wname+"-wp-pv-claim", metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}
	err = pvcClient.Delete(context.Background(), wname+"-mysql-pv-claim", metav1.DeleteOptions{})
	if err != nil {
		getErr = err
	}
	if getErr != nil {
		return getErr
	}
	return nil
}
