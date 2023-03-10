package wordpress

import (
	"context"
	"os"

	"k8smanager/internal"
	"k8smanager/internal/clientgo"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateWordpressDeployment(wname string) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", wname)
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
									ClaimName: wname + "-pv-claim",
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

func CheckIfNamespaceExist(namespace string) error {
	clientset := internal.GetConfig()
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func CreateWordpressService(wname string, port int32) error {
	clientset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", wname)
	nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	err := CheckIfNamespaceExist(namespace)
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.Background(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			log.Error("Failed to create namespace" + namespace)
			return err
		}
		log.Info("Created Namespace " + namespace)
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
					Name:     wname,
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
	_, err = servicesClinet.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Wordpress Service Created Successfully!")
	return nil
}

func CreateWordpressPVC(pname string) error {
	clinetset := internal.GetConfig()
	namespace := clientgo.GetNamespace("wordpress", pname)

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
	log.Info("Created Wordpress PVC")
	return nil
}

func CreateHttpProxy(wname string) error {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}
	httpproxy := map[string]interface{}{
		"apiVersion": "projectcontour.io/v1",
		"kind":       "HTTPProxy",
		"metadata": map[string]interface{}{
			"name":      wname,
			"namespace": clientgo.GetNamespace("wordpress", wname),
			"annotations": map[string]string{
				"projectcontour.io/ingress.class": "wordpress-01",
			},
		},
		"spec": map[string]interface{}{

			"routes": []map[string]interface{}{
				{
					"services": []map[string]interface{}{
						{
							"name": wname,
							"port": 80,
						}}},
				{
					"conditions": []map[string]interface{}{{
						"prefix": "/filemanager",
					}},
					"services": []map[string]interface{}{{
						"name": wname,
						"port": 8080,
					}},
				},
			},
		},
	}

	unstructuredObj := &unstructured.Unstructured{Object: httpproxy}
	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	deploymentRes := schema.GroupVersionResource{Group: "projectcontour.io", Version: "v1", Resource: "httpproxies"}
	result, err := dd.Resource(deploymentRes).Namespace(clientgo.GetNamespace("wordpress", wname)).Create(context.Background(), unstructuredObj, v1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Println(result)
	return nil
}
