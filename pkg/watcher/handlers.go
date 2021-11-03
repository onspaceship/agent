package watcher

import (
	"context"

	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w *Watcher) onUpdate(_, obj interface{}) {
	deployment := obj.(*appsv1.Deployment)

	if deliveryId, ok := deployment.Annotations[config.DeliveryIdAnnotation]; !ok {
		log.WithField("deployment", deployment.Name).Info("No delivery ID found")
	} else {
		// Don't process the wrong deployment revision
		if deployment.Annotations[config.KubernetesRevisionAnnotation] != deployment.Annotations[config.DeliveryRevisionAnnotation] {
			return
		}

		status, err := w.getDeploymentStatus(deployment)
		if err != nil {
			log.WithError(err).Error("Problem getting deployment status")
			return
		}

		log.
			WithField("deployment", deployment.Name).
			WithField("delivery", deliveryId).
			WithField("status", status).
			Info("Updating delivery status with Core")

		w.Core.DeliveryUpdate(deliveryId, status)
	}
}

func (w *Watcher) getDeploymentStatus(deployment *appsv1.Deployment) (string, error) {
	if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
		return "deploying", nil
	}

	replicaSets, err := w.getReplicaSetsForDeployment(deployment)
	currentRev := deployment.Annotations[config.KubernetesRevisionAnnotation]

	for _, replicaSet := range replicaSets.Items {
		if replicaSet.Annotations[config.KubernetesRevisionAnnotation] == currentRev {
			if replicaSet.Status.Replicas != *deployment.Spec.Replicas || replicaSet.Status.AvailableReplicas < replicaSet.Status.Replicas {
				return "deploying", nil
			} else {
				return "complete", nil
			}
		}
	}

	return "error", err
}

func (w *Watcher) getReplicaSetsForDeployment(deployment *appsv1.Deployment) (*appsv1.ReplicaSetList, error) {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	rsList, err := w.Client.AppsV1().ReplicaSets(deployment.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	return rsList, nil
}
