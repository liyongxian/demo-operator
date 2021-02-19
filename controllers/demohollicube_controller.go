/*
Copyright 2021 hollicube.

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

package controllers

import (
	"context"
	gpaasv1alpha1 "demo-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	_ "strconv"
)

// DemoHollicubeReconciler reconciles a DemoHollicube object
type DemoHollicubeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gpaas.hollicube.io,resources=demohollicubes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gpaas.hollicube.io,resources=demohollicubes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployment,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=service,verbs=get;list;watch;create;update;patch;delete
func (r *DemoHollicubeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("demohollicube", req.NamespacedName)

	// your logic here
	// create instance for DemoHollicube
	demoHollicube := &gpaasv1alpha1.DemoHollicube{}
	// search resource
	// 如果没有实例，就返回空结果，这样外部就不再立即调用Reconcile方法
	if err := r.Get(ctx, req.NamespacedName, demoHollicube); err != nil {
		if err := client.IgnoreNotFound(err); err == nil {
			log.Info("no resource found, so go to lifecycle ending")
			return ctrl.Result{}, nil

		} else {
			log.Error(err, " resource error happen")
			return ctrl.Result{}, err
		}
	}

	// get deployment
	deployment := &appsv1.Deployment{}
	// get error or get nothing
	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		// if not get deployment
		if errors.IsNotFound(err) {
			log.Info("deployment not exists")
			// create service
			if err := createServiceIfNotExists(ctx, r, demoHollicube, req); err != nil {
				log.Error(err, "create service error")
				return ctrl.Result{}, err
			}

			// create deployment
			if err := createDeployment(ctx, r, demoHollicube, req); err != nil {
				log.Error(err, "create deployment error")
				return ctrl.Result{}, err
			}
			// create deployment
			return ctrl.Result{}, nil
		} else {
			log.Error(err, "deployment exists.")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// func  create deployment
func createDeployment(ctx context.Context, r *DemoHollicubeReconciler, demoHollicube *gpaasv1alpha1.DemoHollicube, req ctrl.Request) error {
	log := r.Log.WithValues("func", "createDeployment")
	podLabels := map[string]string{
		"app": req.Name,
	}
	// define deployment
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: demoHollicube.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            req.Name,
							Image:           demoHollicube.Spec.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name:          demoHollicube.Spec.Protocol,
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: demoHollicube.Spec.ContainerPort,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(demoHollicube.Spec.CPURequest),
									"memory": resource.MustParse(demoHollicube.Spec.MEMRequest),
								},
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse(demoHollicube.Spec.CPULimit),
									"memory": resource.MustParse(demoHollicube.Spec.MEMLimit),
								},
							},
						},
					},
				},
			},
		},
	}

	// create ref ,when delete rc,delete deployment
	log.Info("set reference")
	if err := controllerutil.SetControllerReference(demoHollicube, deployment, r.Scheme); err != nil {
		log.Error(err, "SetControllerReference error")
		return err
	}
	// create deployment
	log.Info("start create deployment")
	if err := r.Create(ctx, deployment); err != nil {
		log.Error(err, "create deployment error")
		return err
	}
	log.Info("create deployment success")
	return nil
}

// create service
func createServiceIfNotExists(ctx context.Context, r *DemoHollicubeReconciler, demoHollicube *gpaasv1alpha1.DemoHollicube, req ctrl.Request) error {
	log := r.Log.WithValues("func", "createService")
	service := &corev1.Service{}
	err := r.Get(ctx, req.NamespacedName, service)
	// if get nothing,service is ok,and do nothing
	if err == nil {
		log.Info("service exists")
		return nil
	}

	// if error is not NotFound，return error
	if !errors.IsNotFound(err) {
		log.Error(err, "query service error")
		return err
	}

	// define service
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: demoHollicube.Namespace,
			Name:      demoHollicube.Name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     demoHollicube.Spec.ServicePort,
					NodePort: demoHollicube.Spec.NodePort,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": req.Name,
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	// create ref
	if err := controllerutil.SetControllerReference(demoHollicube, service, r.Scheme); err != nil {
		log.Error(err, "SetControllerReference error")
		return err
	}

	// create service
	log.Info("start create service")
	if err := r.Create(ctx, service); err != nil {
		log.Error(err, "create service error")
		return err
	}
	log.Info("create service success")
	return nil
}

func (r *DemoHollicubeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpaasv1alpha1.DemoHollicube{}).
		Complete(r)
}
