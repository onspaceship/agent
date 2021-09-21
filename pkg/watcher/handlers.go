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
		status := getDeploymentStatus(deployment)

		log.
			WithField("deployment", deployment.Name).
			WithField("delivery", deliveryId).
			WithField("status", status).
			Info("Updating delivery status with Core")

		w.Core.DeliveryUpdate(deliveryId, status)
	}
}

func getDeploymentStatus(deployment *appsv1.Deployment) string {
	if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
		return string(appsv1.DeploymentProgressing)
	}

	for _, cond := range deployment.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			return string(appsv1.DeploymentAvailable)
		}
	}

	for _, cond := range deployment.Status.Conditions {
		if cond.Type == appsv1.DeploymentProgressing && cond.Status == corev1.ConditionTrue {
			return string(appsv1.DeploymentProgressing)
		}
	}

	return string(appsv1.DeploymentReplicaFailure)
}
