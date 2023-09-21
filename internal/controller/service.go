package controller

import (
	"fmt"
	bookserverapi "github.com/MobarakHsn/kubebuilder-crd/api/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *BookServerReconciler) EnsureService() error {
	service := &core.Service{}
	if err := r.Client.Get(r.ctx, types.NamespacedName{
		Namespace: r.bookServer.Namespace,
		Name:      r.bookServer.ServiceName(),
	}, service); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Could not find existing service for ", r.bookServer.Name, ", creating one...")
			err := r.Client.Create(r.ctx, r.NewService())
			if err != nil {
				fmt.Printf("Error while creating service %s\n", err)
				return err
			} else {
				fmt.Printf("%s Service Created...\n", r.bookServer.Name)
			}
		} else {
			fmt.Printf("Error getting service %s\n", err)
			return err
		}
	}

	return nil
}

func (r *BookServerReconciler) NewService() *core.Service {
	labels := map[string]string{
		"app":  r.bookServer.Name,
		"kind": "BookServer",
	}
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      r.bookServer.ServiceName(),
			Namespace: r.bookServer.Namespace,
			OwnerReferences: []meta.OwnerReference{
				*meta.NewControllerRef(r.bookServer, bookserverapi.GroupVersion.WithKind("BookServer")),
			},
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Protocol:   "TCP",
					Port:       r.bookServer.Spec.Service.ServicePort,
					TargetPort: intstr.FromInt(int(r.bookServer.Spec.Container.Port)),
					NodePort:   r.bookServer.Spec.Service.ServiceNodePort,
				},
			},
			Selector: labels,
			Type: func() core.ServiceType {
				if r.bookServer.Spec.Service.ServiceType == "NodePort" {
					return core.ServiceTypeNodePort
				} else {
					return core.ServiceTypeClusterIP
				}
			}(),
		},
	}
}
