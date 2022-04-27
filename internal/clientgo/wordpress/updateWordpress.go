package wordpress

import (
	"context"
	"strconv"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wp *WordPress) Update(n int32, pname string, aname string, enable string) error {
	switch aname {
	case "deployment":
		err := UpdateDeployment(n, pname)
		return err
	case "service":
		err := UpdateService(n, pname)
		return err
	case "filemanager":
		err := Filemanagerstatus(pname, enable)
		return err
	}
	return nil
}

func UpdateDeployment(n int32, pname string) error {
	log.Info("Updating Deployment")
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

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

func Createfilemanager(pname string) error {
	log.Info("Updating Deployment file")
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	result, err := deploymentsClient.Get(context.Background(), pname, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get data")
		return err
	}

	container := result.Spec.Template.Spec.Containers
	var containers v1.Container
	containers.Image = "hurlenko/filebrowser"
	containers.ImagePullPolicy = v1.PullAlways
	containers.Name = pname + "-file"

	var volumes v1.VolumeMount
	volumes.Name = pname + "-wordpress-persistent-storage"
	volumes.MountPath = "/var/www/html"
	vol := []v1.VolumeMount{}
	vol = append(vol, volumes)
	containers.VolumeMounts = vol

	var ports v1.ContainerPort
	ports.Name = pname + "-file"
	ports.Protocol = v1.ProtocolTCP
	ports.ContainerPort = 80
	port := []v1.ContainerPort{}
	port = append(port, ports)
	containers.Ports = port
	container = append(container, containers)
	result.Spec.Template.Spec.Containers = container

	_, err = deploymentsClient.Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Error("Failed to Create filemanager")
		return err
	}

	log.Info("Created filemanager")
	return nil
}

func DeleteFilemanager(wname string) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", wname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, _ := deploymentsClient.Get(context.Background(), wname, metav1.GetOptions{})
	var container v1.Container

	container.Image = "wordpress:4.8-apache"
	container.Name = wname

	var env1 v1.EnvVar

	env1.Name = "WORDPRESS_DB_HOST"
	env1.Value = wname + "-mysql"

	var env v1.EnvVar
	env.Name = "WORDPRESS_DB_PASSWORD"
	var value v1.EnvVarSource
	var secret v1.SecretKeySelector
	secret.LocalObjectReference.Name = wname + "-mysql-pass"
	secret.Key = "password"
	value.SecretKeyRef = &secret
	env.ValueFrom = &value

	envVar := []v1.EnvVar{}
	envVar = append(envVar, env1)
	envVar = append(envVar, env)
	container.Env = envVar

	var port v1.ContainerPort

	port.Name = wname
	port.Protocol = v1.ProtocolTCP
	port.ContainerPort = 80

	containerPort := []v1.ContainerPort{}
	containerPort = append(containerPort, port)
	container.Ports = containerPort

	var volumemount v1.VolumeMount

	volumemount.Name = wname + "-wordpress-persistent-storage"
	volumemount.MountPath = "/var/www/html"

	vol := []v1.VolumeMount{}
	vol = append(vol, volumemount)
	container.VolumeMounts = vol

	var containers []v1.Container
	containers = append(containers, container)

	result.Spec.Template.Spec.Containers = containers

	_, err := deploymentsClient.Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Error("Failed to Update")
		return err
	}

	log.Info("Updated Deployment to delete filemanager")
	return nil

}

func Filemanagerstatus(wname string, enable string) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", wname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, _ := deploymentsClient.Get(context.Background(), wname, metav1.GetOptions{})
	containers := result.Spec.Template.Spec.Containers
	Enable, err := strconv.ParseBool(enable)

	if err != nil {
		log.Error(err)
		return err
	}
	var status bool
	for _, item := range containers {
		log.Info(item.Name)
		if item.Name == wname+"-file" {
			status = true
			break
		}
		status = false
	}

	if Enable {
		if status {
			return nil
		}
		err := Createfilemanager(wname)
		if err != nil {
			return err
		}
	} else {
		if !status {
			return nil
		}
		err := DeleteFilemanager(wname)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateService(n int32, pname string) error {
	log.Info("Updating Service")
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)
	ServiceClient := clientset.CoreV1().Services(namespace)

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
