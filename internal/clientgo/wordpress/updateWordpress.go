package wordpress

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wordpress.com/internal"
)

func (wp *WordPress) Update(n int32, pname string, aname string) error {
	switch aname {
	case "deployment":
		err := UpdateDeployment(n, pname)
		return err
	case "service":
		err := UpdateService(n, pname)
		return err
	}
	return nil
}

func UpdateDeployment(n int32, pname string) error {
	log.Info("Updating Deployment")
	clientset := internal.GetConfig()
	deploymentsClient := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE"))

	result, err := deploymentsClient.Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get data")
		return err
	}
	result.Spec.Replicas = int32ptr(n)
	_, err = deploymentsClient.Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Error("Failed to update")
		return err
	}

	log.Info("Updated Deployment")
	return nil
}

func UpdateService(n int32, pname string) error {
	log.Info("Updating Service")
	clientset := internal.GetConfig()
	ServiceClient := clientset.CoreV1().Services(os.Getenv("NAMESPACE"))

	result, err := ServiceClient.Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get data")
		return err
	}
	result.Spec.Ports[0].NodePort = n
	_, err = ServiceClient.Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Error("Failed to update")
		return err
	}

	log.Info("Updated Service with port")
	return nil
}