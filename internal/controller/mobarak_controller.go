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
	bookserverapi "github.com/MobarakHsn/kubebuilder_crd/api/v1"
	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BookServerReconciler reconciles a BookServer object
type BookServerReconciler struct {
	client     client.Client
	log        logr.Logger
	ctx        context.Context
	Scheme     *runtime.Scheme
	bookServer *bookserverapi.BookServer
}

//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.github.com,resources=mobaraks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BookServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *BookServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// set up operator logger with resource key
	r.log = ctrl.Log.WithValues("BookServer", req.NamespacedName)
	r.ctx = ctx

	// get bookserver and ensure it exists
	bookServer := &bookserverapi.BookServer{}
	if err := r.client.Get(ctx, req.NamespacedName, bookServer); err != nil {
		r.log.Error(err, fmt.Sprintf("Unable to Get BookServer %s/%s", req.Namespace, req.Name))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	r.bookServer = bookServer
	fmt.Println("BookServer fetched", req.NamespacedName)

	var serviceInstance core.Service
	serviceName := func() string {
		if bookServer.Spec.Service.ServiceName == "" {
			return bookServer.Name
		} else {
			return bookServer.Spec.Service.ServiceName
		}
	}()
	nsname = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}
	if err := r.Get(ctx, nsname, &serviceInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Could not find existing service for ", bookServer.Name, ", creating one...")
			err := r.Create(ctx, NewService(&bookServer, serviceName))
			if err != nil {
				fmt.Printf("Error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", bookServer.Name)
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
	customDeployNameField = ".spec.deploymentName"
	// svcOwnerKey    = ".metadata.controller"
	// apiGVStr       = bookserverapi.GroupVersion.String()
	// ourKind        = "BookServer"
)
