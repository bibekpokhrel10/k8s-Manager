package wordpress

import (
	"context"
	"os"

	"wordpress.com/internal"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateWordpressDeployment(wname string) error {
	clientset := internal.GetConfig()
	deploymentsClient := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE"))

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
							Image: "wordpress:4.8-apache",
							Name:  wname,
							Env: []corev1.EnvVar{
								{
									Name:  "WORDPRESS_DB_HOST",
									Value: wname + "-mysql",
								},
								{
									Name: "WORDPRESS_DB_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "mysql-pass",
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
									Name:      wname + "-wordpress-persistent-storage",
									MountPath: "/var/www/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: wname + "-wordpress-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: wname + "-wp-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info("Creating Wordpress deployment..")
	result, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Wordpress deployment " + result.GetObjectMeta().GetName())
	return nil
}

func int32ptr(i int32) *int32 {
	return &i
}

func CreateWordpressService(wname string, port int32) error {
	clientset := internal.GetConfig()
	servicesClinet := clientset.CoreV1().Services(os.Getenv("NAMESPACE"))
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wname,
			Namespace: "bibek",
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

	log.Info("Creating Wordpress service...")
	_, err := servicesClinet.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Wordpress Service Created Successfully!")
	return nil
}

func CreateWordpressPVC(pname string) error {
	clinetset := internal.GetConfig()
	pvcClinet := clinetset.CoreV1().PersistentVolumeClaims(os.Getenv("NAMESPACE"))

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname + "-wp-pv-claim",
			Namespace: "bibek",
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
	log.Info("Created Wordpress PVC")
	return nil
}
