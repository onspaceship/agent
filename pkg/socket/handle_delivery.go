package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/onspaceship/agent/pkg/client"
	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type deliveryPayload struct {
	DeliveryId string `json:"delivery_id"`
	AppId      string `json:"app_id"`
	AppHandle  string `json:"app_handle"`
	TeamHandle string `json:"team_handle"`
	ImageURI   string `json:"image_uri"`
}

func handleDelivery(jsonPayload []byte, socket *socket) {
	var payload deliveryPayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	logline := log.WithField("app_id", payload.AppId).WithField("delivery_id", payload.DeliveryId)
	logline.Info("Handling delivery")

	core := client.NewClient()

	// Run the release job, if present

	ctx := context.Background()
	releaseJobs, err := socket.client.BatchV1().Jobs("").List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set{config.AppIdLabel: payload.AppId}.AsSelector().String(),
	})
	if err != nil {
		logline.WithError(err).Fatal("Could not get deployments for App")
	}

	for _, job := range releaseJobs.Items {
		jobLog := logline.WithField("job", fmt.Sprintf("%s/%s", job.Namespace, job.Name))
		jobLog.Info("Running release job")

		core.DeliveryUpdate(payload.DeliveryId, "releasing")

		// Delete the last job
		jobLog.Info("Deleting prior job")
		bg := metav1.DeletePropagationBackground
		err = socket.client.BatchV1().Jobs(job.Namespace).Delete(ctx, job.Name, metav1.DeleteOptions{PropagationPolicy: &bg})
		if err != nil {
			jobLog.WithError(err).Fatal("Could not delete job")
		}

		// Set metadata and clear imutable fields
		job.ObjectMeta.Annotations[config.DeliveryIdAnnotation] = payload.DeliveryId
		job.ObjectMeta.Annotations[config.AppHandleAnnotation] = payload.AppHandle
		job.ObjectMeta.Annotations[config.TeamHandleAnnotation] = payload.TeamHandle

		job.Spec.Template.ObjectMeta.Labels = nil
		job.Spec.Selector = nil

		job.ObjectMeta.ResourceVersion = ""
		job.ObjectMeta.UID = ""

		for i, container := range job.Spec.Template.Spec.Containers {
			container.Image = payload.ImageURI
			job.Spec.Template.Spec.Containers[i] = container
		}

		// Create the new job
		jobLog.Info("Creating a new job")
		_, err = socket.client.BatchV1().Jobs(job.Namespace).Create(ctx, &job, metav1.CreateOptions{})
		if err != nil {
			jobLog.WithError(err).Fatal("Could not create job")
		}

		timeout := time.After(time.Minute * 10)
		ticker := time.NewTicker(time.Second * 3)

		// Check every 3 seconds until it finishes and timeout after 10 minutes
	CheckIfJobRunning:
		for {
			select {
			case <-ticker.C:
				latestJob, _ := socket.client.BatchV1().Jobs(job.Namespace).Get(ctx, job.Name, metav1.GetOptions{})
				if latestJob.Status.Failed > 0 {
					jobLog.Error("Release job failed!")
					core.DeliveryUpdate(payload.DeliveryId, "error")
					return
				} else if latestJob.Status.Active == 0 {
					jobLog.Info("Release job succeeded!")
					core.DeliveryUpdate(payload.DeliveryId, "deploying")
					break CheckIfJobRunning
				}

				jobLog.Info("Release job still running...")
			case <-timeout:
				ticker.Stop()
				jobLog.Warn("Timeout reached before job completed")
				break CheckIfJobRunning
			}
		}
	}

	// Update the deployment

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	deployments, err := socket.client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set{config.AppIdLabel: payload.AppId}.AsSelector().String(),
	})
	if err != nil {
		logline.WithError(err).Fatal("Could not get deployments for App")
	}

	logline.Infof("Found %d deployments", len(deployments.Items))

	if len(deployments.Items) > 0 {
		core.DeliveryUpdate(payload.DeliveryId, "deploying")
	}

	for _, deployment := range deployments.Items {
		deployLog := logline.WithField("deployment", fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name))
		deployLog.Info("Updating image")

		// Update the container image
		for i, container := range deployment.Spec.Template.Spec.Containers {
			container.Image = payload.ImageURI
			deployment.Spec.Template.Spec.Containers[i] = container
		}

		// Update the metadata
		deployment.ObjectMeta.Annotations[config.DeliveryIdAnnotation] = payload.DeliveryId
		deployment.ObjectMeta.Annotations[config.AppHandleAnnotation] = payload.AppHandle
		deployment.ObjectMeta.Annotations[config.TeamHandleAnnotation] = payload.TeamHandle

		_, err = socket.client.AppsV1().Deployments(deployment.Namespace).Update(ctx, &deployment, metav1.UpdateOptions{})
		if err != nil {
			deployLog.WithError(err).Info("Could not update image")
			return
		} else {
			deployLog.Info("Image updated!")
		}
	}
}

func init() {
	handlers["delivery"] = handleDelivery
}
