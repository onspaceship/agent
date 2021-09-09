package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apex/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	AppIdAnnotation      = "onspaceship.com/app-id"
	DeliveryIdAnnotation = "onspaceship.com/delivery-id"
)

type deliveryPayload struct {
	AppId      string `json:"app_id"`
	DeliveryId string `json:"delivery_id"`
	ImageURI   string `json:"image_uri"`
}

func handleDelivery(jsonPayload []byte, socket *socket) {
	var payload deliveryPayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	log.WithField("app_id", payload.AppId).WithField("delivery_id", payload.DeliveryId).Info("Handling delivery")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	deployments, err := socket.client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set{AppIdAnnotation: payload.AppId}.AsSelector().String(),
	})
	if err != nil {
		log.WithError(err).Fatal("Could not get Kubernetes deployment for Agent")
	}

	for _, deployment := range deployments.Items {
		logline := log.WithField("deployment", fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name))
		logline.Info("Updating image")

		for i, container := range deployment.Spec.Template.Spec.Containers {
			container.Image = payload.ImageURI
			deployment.Spec.Template.Spec.Containers[i] = container
		}

		_, err = socket.client.AppsV1().Deployments(deployment.Namespace).Update(ctx, &deployment, metav1.UpdateOptions{})
		if err != nil {
			logline.WithError(err).Info("Could not update image")
			return
		} else {
			logline.Info("Image updated!")
		}
	}
}

func init() {
	handlers["delivery"] = handleDelivery
}
