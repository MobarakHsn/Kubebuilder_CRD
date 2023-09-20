package controller

import (
	crdv1 "github.com/MobarakHsn/kubebuilder_crd/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewService(customRes *crdv1.BookServer, serviceName string) *corev1.Service {
	labels := map[string]string{
		"app":  customRes.Name,
		"kind": "BookServer",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: customRes.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customRes, crdv1.GroupVersion.WithKind("BookServer")),
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
