package joomla

import (
	"context"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateJoomlaDeployment(wname string) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("joomla", wname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: wname,
			Labels: map[string]string{
				"app": wname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  wname,
					"tier": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  wname,
						"tier": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "joomla",
							Name:  wname,
							Env: []corev1.EnvVar{
								{
									Name:  "JOOMLA_DB_HOST",
									Value: wname + "-mysql",
								},
								{
									Name: "JOOMLA_DB_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: wname + "-mysql-pass",
											},
											Key: "password",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          wname,
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      wname + "-joomla-persistent-storage",
									MountPath: "/var/www/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: wname + "-joomla-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: wname + "-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info("Creating Joomla deployment..")
	result, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Joomla deployment " + result.GetObjectMeta().GetName())
	return nil
}

func int32ptr(i int32) *int32 {
	return &i
}

func CheckIfNamespaceExist(namespace string) error {
	clientset := internal.GetConfig()
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func CreateJoomlaService(wname string, port int32) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("joomla", wname)
	nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	err := CheckIfNamespaceExist(namespace)
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.Background(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			log.Error("Failed to create namespace" + namespace)
			return err
		}
		log.Info("Created Namespace" + namespace)
	}
	servicesClinet := clientset.CoreV1().Services(namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wname,
			Namespace: namespace,
			Labels: map[string]string{
				"app": wname,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: "TCP",
					NodePort: port,
				},
			},
			Selector: map[string]string{
				"app":  wname,
				"tier": "frontend",
			},
			Type: "NodePort",
		},
	}

	log.Info("Creating Joomla service...")
	_, err = servicesClinet.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Joomla Service Created Successfully!")
	return nil
}

func CreateJoomlaPVC(pname string) error {
	clinetset := internal.GetConfig()
	namespace := clientgo.GetNamespace("joomla", pname)
	pvcClinet := clinetset.CoreV1().PersistentVolumeClaims(namespace)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname + "-pv-claim",
			Namespace: namespace,
			Labels: map[string]string{
				"app": pname,
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := pvcClinet.Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Joomla PVC")
	return nil
}
