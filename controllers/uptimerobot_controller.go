/*
Copyright 2022.

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
	"fmt"
	"github.com/bennsimon/uptimerobot-operator/util/monitorutil"
	"github.com/bennsimon/uptimerobot-tooling/pkg/util/httputil"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strconv"

	network "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type UptimerobotReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	UtilProvider
}

type UtilProvider interface {
	CreateMonitor(host string, labels map[string]string) error
	DeleteMonitor(host string, labels map[string]string) error
}

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;watch;list;

func (r *UptimerobotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nameSpacedName := req.NamespacedName
	ingress := &network.Ingress{}
	err := r.Get(ctx, nameSpacedName, ingress)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	hosts := buildHostSchemeMap(ingress)
	for host, scheme := range hosts {
		hostWithScheme := scheme + "://" + host
		if err := r.UtilProvider.CreateMonitor(hostWithScheme, ingress.Labels); err != nil {
			log.Log.Error(err, fmt.Sprintf("Monitor %s not successfully created/updated", hostWithScheme))
			return ctrl.Result{}, nil
		}
		log.Log.Info(fmt.Sprintf("Monitor %s successfully created/updated", hostWithScheme))
	}

	return ctrl.Result{}, nil
}

func (r *UptimerobotReconciler) CreateMonitor(host string, labels map[string]string) error {
	return monitorutil.CreateMonitor(host, labels)
}

func (r *UptimerobotReconciler) DeleteMonitor(host string, labels map[string]string) error {
	return monitorutil.DeleteMonitor(host, labels)
}

func buildHostSchemeMap(ingress *network.Ingress) map[string]string {
	hosts := map[string]string{}

	for _, rule := range ingress.Spec.Rules {
		hosts[rule.Host] = "http"
	}
	if len(ingress.Spec.TLS) > 0 {
		for _, tls := range ingress.Spec.TLS {
			for _, rule := range tls.Hosts {
				hosts[rule] = "https"
			}
		}
	}
	return hosts
}

func (r *UptimerobotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&network.Ingress{}, builder.WithPredicates(r.FilterEnabledIngress())).
		Complete(r)
}

func (r *UptimerobotReconciler) FilterEnabledIngress() predicate.Predicate {
	return predicate.Funcs{CreateFunc: func(event event.CreateEvent) bool {
		return r.filterCreateEvent(event)
	}, UpdateFunc: func(updateEvent event.UpdateEvent) bool {
		return r.filterUpdateEvent(updateEvent)
	}, DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
		return r.filterDeleteEvent(deleteEvent)
	}, GenericFunc: func(genericEvent event.GenericEvent) bool {
		return r.filterGenericEvent(genericEvent)
	}}
}

func (r *UptimerobotReconciler) filterGenericEvent(genericEvent event.GenericEvent) bool {
	if genericEvent.Object != nil {
		return r.hasEnabledUptimeRobotMonitor(genericEvent.Object.GetLabels())
	}
	return false
}

func (r *UptimerobotReconciler) filterDeleteEvent(deleteEvent event.DeleteEvent) bool {
	if deleteEvent.Object != nil {
		enabled := r.hasEnabledUptimeRobotMonitor(deleteEvent.Object.GetLabels())
		if enabled {
			go r.cleanUpAfterIngressDeletion(deleteEvent.Object.GetLabels())
		}
	}
	return false
}

func (r *UptimerobotReconciler) filterUpdateEvent(updateEvent event.UpdateEvent) bool {
	if updateEvent.ObjectNew != nil {
		return r.hasEnabledUptimeRobotMonitor(updateEvent.ObjectNew.GetLabels())
	}
	return false
}

func (r *UptimerobotReconciler) filterCreateEvent(event event.CreateEvent) bool {
	if event.Object != nil {
		return r.hasEnabledUptimeRobotMonitor(event.Object.GetLabels())
	}
	return false
}

func (r *UptimerobotReconciler) cleanUpAfterIngressDeletion(labels map[string]string) {
	err := r.UtilProvider.DeleteMonitor("", labels)
	if err != nil {
		log.Log.Error(err, fmt.Sprintf("Monitor %s not successfully deleted", labels[monitorutil.GetUptimeRobotMonitorLabelPrefix()+httputil.FriendlyNameField]))
	} else {
		log.Log.Info(fmt.Sprintf("Monitor %s successfully deleted", labels[monitorutil.GetUptimeRobotMonitorLabelPrefix()+httputil.FriendlyNameField]))
	}
}

func (r *UptimerobotReconciler) hasEnabledUptimeRobotMonitor(labelMap map[string]string) bool {
	if val, exists := labelMap[monitorutil.GetUptimeRobotLabelDomain()]; exists {
		isEnabled, err := strconv.ParseBool(val)
		if err != nil {
			return false
		}
		return isEnabled
	}
	return false
}
