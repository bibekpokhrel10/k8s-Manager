package internal

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSecretKey() {
	clientset := GetConfig()
	secretClinet := clientset.CoreV1().Secrets(os.Getenv("NAMESPACE"))
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql-pass",
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"Name":     []byte(os.Getenv("SECRETNAME")),
			"password": []byte(os.Getenv("PASSWORD")),
		},
	}
	_, err := secretClinet.Get(context.Background(), "mysql-pass", metav1.GetOptions{})
	if err == nil {
		log.Info("Updating Secret Key")
		_, err = secretClinet.Update(context.Background(), secret, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Secret Key Updated Successfully")
	} else {
		fmt.Println("Creating Secret Key")
		_, err = secretClinet.Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Secret Key Created Successfully!")
	}
}
