package wordpress

import (
	"context"

	"k8smanager/internal"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSecretKey(wname string) {
	clientset := internal.GetConfig()
	secretClinet := clientset.CoreV1().Secrets(wname)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: wname + "-mysql-pass",
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"Name":     []byte(wname + "-mysql-pass"),
			"password": []byte(wname),
		},
	}
	_, err := secretClinet.Get(context.Background(), wname+"-mysql-pass", metav1.GetOptions{})
	if err == nil {
		log.Info("Updating Secret Key")
		_, err = secretClinet.Update(context.Background(), secret, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Secret Key Updated Successfully")
	} else {
		log.Info("Creating Secret Key")
		_, err = secretClinet.Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Secret Key Created Successfully!")
	}
}

func CreateDatabasePvc(pname string) error {
	clinetset := internal.GetConfig()
	pvcClinet := clinetset.CoreV1().PersistentVolumeClaims(pname)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname + "-mysql-pv-claim",
			Namespace: pname,
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
	log.Info("Created Database PVC")
	return nil
}

func CreateDatabaseDeployment(dname string) error {
	clientset := internal.GetConfig()
	deploymentsClient := clientset.AppsV1().Deployments(dname)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dname + "-mysql",
			Labels: map[string]string{
				"app":  dname,
				"tier": "mysql",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  dname,
					"tier": "mysql",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  dname,
						"tier": "mysql",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  dname + "-mysql",
							Image: "mysql:5.6",
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: dname,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          dname + "-mysql",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      dname + "-mysql-persistent-storage",
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: dname + "-mysql-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: dname + "-mysql-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}
	log.Info("Creating Database deployment..")
	result, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Created Database Deployment : " + result.GetObjectMeta().GetName())
	return nil
}

func CreateDatabaseService(dname string) error {
	clientset := internal.GetConfig()
	servicesClinet := clientset.CoreV1().Services(dname)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dname + "-mysql",
			Namespace: dname,

			Labels: map[string]string{
				"app": dname,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 3306,
				},
			},
			Selector: map[string]string{
				"app":  dname,
				"tier": "mysql",
			},
			ClusterIP: "None",
		},
	}
	log.Info("Creating Database service...")
	_, err := servicesClinet.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Database Service Successfully!")
	return nil

}
