package controller

import (
	crdv1 "github.com/MobarakHsn/kubebuilder_crd/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewDeployment(customRes *crdv1.Mobarak, deploymentName string) *appsv1.Deployment {
	labels := map[string]string{
		"app":  customRes.Name,
		"kind": "Mobarak",
	}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: customRes.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customRes, crdv1.GroupVersion.WithKind("Mobarak")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: customRes.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: customRes.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  customRes.Name,
							Image: customRes.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: customRes.Spec.Container.Port,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewService(customRes *crdv1.Mobarak, serviceName string) *corev1.Service {
	labels := map[string]string{
		"app":  customRes.Name,
		"kind": "Mobarak",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: customRes.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customRes, crdv1.GroupVersion.WithKind("Mobarak")),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       customRes.Spec.Service.ServicePort,
					TargetPort: intstr.FromInt(int(customRes.Spec.Container.Port)),
					NodePort:   customRes.Spec.Service.ServiceNodePort,
				},
			},
			Selector: labels,
			Type: func() corev1.ServiceType {
				if customRes.Spec.Service.ServiceType == "NodePort" {
					return corev1.ServiceTypeNodePort
				} else {
					return corev1.ServiceTypeClusterIP
				}
			}(),
		},
	}
}
