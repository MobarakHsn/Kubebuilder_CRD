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
	"github.com/go-logr/logr"

	bookserverapi "github.com/MobarakHsn/kubebuilder-crd/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BookServerReconciler reconciles a BookServer object
type BookServerReconciler struct {
	Client     client.Client
	Log        logr.Logger
	ctx        context.Context
	Scheme     *runtime.Scheme
	bookServer *bookserverapi.BookServer
}

//+kubebuilder:rbac:groups=bookcrd.github.com,resources=bookservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bookcrd.github.com,resources=bookservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bookcrd.github.com,resources=bookservers/finalizers,verbs=update

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
	r.Log = ctrl.Log.WithValues("BookServer", req.NamespacedName)
	r.ctx = ctx
	// get bookserver and ensure it exists
	bookServer := &bookserverapi.BookServer{}
	if err := r.Client.Get(ctx, req.NamespacedName, bookServer); err != nil {
		r.Log.Error(err, fmt.Sprintf("Unable to Get BookServer %s/%s", req.Namespace, req.Name))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if bookServer.Spec.DeploymentName != bookServer.DeploymentName() {
		bookServer.Spec.DeploymentName = bookServer.DeploymentName()
		if err := r.Client.Update(ctx, bookServer); err != nil {
			return ctrl.Result{}, err
		}
		r.Log.Info("Book server updated")
	}

	r.Log.Info(fmt.Sprintf("BookServer fetched %v", req.NamespacedName))
	r.bookServer = bookServer
	if err := r.EnsureDeployment(); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.EnsureService(); err != nil {
		return ctrl.Result{}, err
	}
	return r.reconcile()
}
