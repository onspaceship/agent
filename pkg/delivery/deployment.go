package delivery

import (
	"fmt"

	"github.com/onspaceship/agent/pkg/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (d *delivery) updateDeployment() {
	deployments, err := d.kubernetes.AppsV1().Deployments("").List(d.ctx, metav1.ListOptions{
		LabelSelector: labels.Set{config.AppIdLabel: d.appId}.AsSelector().String(),
	})
	if err != nil {
		d.log.WithError(err).Fatal("Could not get deployments for App")
	}

	d.log.Infof("Found %d deployments", len(deployments.Items))

	if len(deployments.Items) > 0 {
		d.core.DeliveryUpdate(d.deliveryId, "deploying")
	}

	for _, deployment := range deployments.Items {
		deployLog := d.log.WithField("deployment", fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name))
		deployLog.Info("Updating image")

		// Update the container image
		for i, container := range deployment.Spec.Template.Spec.Containers {
			container.Image = d.imageURI
			deployment.Spec.Template.Spec.Containers[i] = container
		}

		secrets, err := d.ensureImagePullSecrets(deployment.Namespace)
		if err != nil {
			deployLog.WithError(err).Info("Could not get image pull secrets")
			return
		}

		deployment.Spec.Template.Spec.ImagePullSecrets = secrets

		// Update the metadata
		deployment.ObjectMeta.Annotations[config.DeliveryIdAnnotation] = d.deliveryId
		deployment.ObjectMeta.Annotations[config.AppHandleAnnotation] = d.appHandle
		deployment.ObjectMeta.Annotations[config.TeamHandleAnnotation] = d.teamHandle

		_, err = d.kubernetes.AppsV1().Deployments(deployment.Namespace).Update(d.ctx, &deployment, metav1.UpdateOptions{})
		if err != nil {
			deployLog.WithError(err).Info("Could not update image")
			return
		} else {
			deployLog.Info("Image updated!")
		}
	}
}
