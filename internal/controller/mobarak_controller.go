/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	crdv1 "github.com/MobarakHsn/kubebuilder_crd/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MobarakReconciler reconciles a Mobarak object
type MobarakReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mobarak object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *MobarakReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//log := log.FromContext(ctx)
	//log := r.Log.WithValues("ReqName", req.Name, "ReqNameSpace", req.Namespace)
	//fmt.Println(log)
	// TODO(user): your logic here
	defer fmt.Println("reconciliation done")
	var customRes crdv1.Mobarak
	if err := r.Get(ctx, req.NamespacedName, &customRes); err != nil {
		fmt.Println(err, "Unable to fetch mobarakcrd")
		//log.Error(err, "Unable to fetch mobarakcrd")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("Mobarak fetched", req.NamespacedName)
	var deploymentInstance appsv1.Deployment
	deploymentName := func() string {
		if customRes.Spec.DeploymentName == "" {
			return customRes.Name
		} else {
			return customRes.Spec.DeploymentName
		}
	}()
	nsname := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}
	if err := r.Get(ctx, nsname, &deploymentInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Could not find existing deployment for ", customRes.Name, ", creating one...")
			err := r.Create(ctx, NewDeployment(&customRes, deploymentName))
			if err != nil {
				cnt := int32(0)
				customCopy := customRes.DeepCopy()
				customCopy.Status.AvailableReplicas = &cnt
				if err := r.Update(ctx, customCopy); err != nil {
					fmt.Printf("Error updating Mobarak %s\n", err)
					return ctrl.Result{}, err
				}
				fmt.Printf("Error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", customRes.Name)
			}
		} else {
			fmt.Printf("Error fetching deployment %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if customRes.Spec.Replicas != nil && *customRes.Spec.Replicas != *deploymentInstance.Spec.Replicas {
			fmt.Println(*customRes.Spec.Replicas, *deploymentInstance.Spec.Replicas)
			fmt.Println("Deployment replica miss match.....updating")
			cnt := *deploymentInstance.Spec.Replicas
			deploymentInstance.Spec.Replicas = customRes.Spec.Replicas
			if err := r.Update(ctx, &deploymentInstance); err != nil {
				fmt.Printf("Error updating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				customCopy := customRes.DeepCopy()
				customCopy.Status.AvailableReplicas = &cnt
				if err := r.Update(ctx, customCopy); err != nil {
					fmt.Printf("Error updating Mobarak %s\n", err)
					return ctrl.Result{}, err
				}
			}
			fmt.Println("Deployment updated")
		}
	}
	var serviceInstance corev1.Service
	serviceName := func() string {
		if customRes.Spec.Service.ServiceName == "" {
			return customRes.Name
		} else {
			return customRes.Spec.Service.ServiceName
		}
	}()
	nsname = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}
	if err := r.Get(ctx, nsname, &serviceInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Could not find existing service for ", customRes.Name, ", creating one...")
			err := r.Create(ctx, NewService(&customRes, serviceName))
			if err != nil {
				fmt.Printf("Error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", customRes.Name)
			}
		} else {
			fmt.Printf("Error fetching deployment %s\n", err)
			return ctrl.Result{}, err
		}
	}
	//controllerutil.SetControllerReference()
	return ctrl.Result{}, nil
}

var (
// deployOwnerKey = ".metadata.controller"
// svcOwnerKey    = ".metadata.controller"
// apiGVStr       = crdv1.GroupVersion.String()
// ourKind        = "Mobarak"
)

// SetupWithManager sets up the controller with the Manager.
func (r *MobarakReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.Mobarak{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
