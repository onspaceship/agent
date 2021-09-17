package watcher

import (
	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (w *Watcher) onUpdate(_, obj interface{}) {
	deployment := obj.(*appsv1.Deployment)

	if deliveryId, ok := deployment.Annotations[config.DeliveryIdAnnotation]; !ok {
		log.WithField("deployment", deployment.Name).Info("No delivery ID found")
	} else {
		status := getDeploymentStatus(deployment.Status)

		log.
			WithField("deployment", deployment.Name).
			WithField("delivery", deliveryId).
			WithField("status", status).
			Info("Updating delivery status with Core")

		w.Core.DeliveryUpdate(deliveryId, status)
	}
}

func getDeploymentStatus(status appsv1.DeploymentStatus) string {
	for i := range status.Conditions {
		if status.Conditions[i].Status == corev1.ConditionTrue {
			return string(status.Conditions[i].Type)
		}
	}

	return "Created"
}
