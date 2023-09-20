package controller

import (
	"fmt"
	crdv1 "github.com/MobarakHsn/kubebuilder_crd/api/v1"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *BookServerReconciler) EnsureDeployment() error {
	deployment := &apps.Deployment{}
	deploymentName := bookServer.DeploymentName()

	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: bookServer.Namespace,
		Name:      deploymentName,
	}, deployment); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Could not find existing deployment for ", bookServer.Name, ", creating one...")
			deployment = r.NewDeployment(deploymentName)
			err := r.client.Create(ctx, deployment)
			if err != nil {
				cnt := int32(0)
				customCopy := bookServer.DeepCopy()
				customCopy.Status.AvailableReplicas = &cnt
				if err := r.client.Update(ctx, customCopy); err != nil {
					fmt.Printf("Error updating BookServer %s\n", err)
					return ctrl.Result{}, err
				}
				fmt.Printf("Error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", bookServer.Name)
			}
		} else {
			fmt.Printf("Error fetching deployment %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if bookServer.Spec.Replicas != nil && *bookServer.Spec.Replicas != *deployment.Spec.Replicas {
			fmt.Println(*bookServer.Spec.Replicas, *deployment.Spec.Replicas)
			fmt.Println("Deployment replica miss match.....updating")
			cnt := *deployment.Spec.Replicas
			deployment.Spec.Replicas = bookServer.Spec.Replicas
			if err := r.client.Update(ctx, deployment); err != nil {
				fmt.Printf("Error updating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				customCopy := bookServer.DeepCopy()
				customCopy.Status.AvailableReplicas = &cnt
				if err := r.client.Update(ctx, customCopy); err != nil {
					fmt.Printf("Error updating BookServer %s\n", err)
					return ctrl.Result{}, err
				}
			}
			fmt.Println("Deployment updated")
		}
	}
}

func (r *BookServerReconciler) NewDeployment(deploymentName string) *apps.Deployment {
	labels := map[string]string{
		"app":  r.bookServer.Name,
		"kind": "BookServer",
	}
	return &apps.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: r.bookServer.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(r.bookServer, crdv1.GroupVersion.WithKind("BookServer")),
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: r.bookServer.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: r.bookServer.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  r.bookServer.Name,
							Image: r.bookServer.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: r.bookServer.Spec.Container.Port,
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
