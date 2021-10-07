package delivery

import (
	"fmt"
	"time"

	"github.com/onspaceship/agent/pkg/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (d *delivery) runReleaseJob() {
	releaseJobs, err := d.kubernetes.BatchV1().Jobs("").List(d.ctx, metav1.ListOptions{
		LabelSelector: labels.Set{config.AppIdLabel: d.appId}.AsSelector().String(),
	})
	if err != nil {
		d.log.WithError(err).Fatal("Could not get deployments for App")
	}

	for _, job := range releaseJobs.Items {
		jobLog := d.log.WithField("job", fmt.Sprintf("%s/%s", job.Namespace, job.Name))
		jobLog.Info("Running release job")

		d.core.DeliveryUpdate(d.deliveryId, "releasing")

		// Delete the last job
		jobLog.Info("Deleting prior job")
		bg := metav1.DeletePropagationBackground
		err = d.kubernetes.BatchV1().Jobs(job.Namespace).Delete(d.ctx, job.Name, metav1.DeleteOptions{PropagationPolicy: &bg})
		if err != nil {
			jobLog.WithError(err).Fatal("Could not delete job")
		}

		// Set metadata and clear imutable fields
		job.ObjectMeta.Annotations[config.DeliveryIdAnnotation] = d.deliveryId
		job.ObjectMeta.Annotations[config.AppHandleAnnotation] = d.appHandle
		job.ObjectMeta.Annotations[config.TeamHandleAnnotation] = d.teamHandle

		job.Spec.Template.ObjectMeta.Labels = nil
		job.Spec.Selector = nil

		job.ObjectMeta.ResourceVersion = ""
		job.ObjectMeta.UID = ""

		for i, container := range job.Spec.Template.Spec.Containers {
			container.Image = d.imageURI
			job.Spec.Template.Spec.Containers[i] = container
		}

		// Create the new job
		jobLog.Info("Creating a new job")
		_, err = d.kubernetes.BatchV1().Jobs(job.Namespace).Create(d.ctx, &job, metav1.CreateOptions{})
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
				latestJob, _ := d.kubernetes.BatchV1().Jobs(job.Namespace).Get(d.ctx, job.Name, metav1.GetOptions{})
				if latestJob.Status.Failed > 0 {
					jobLog.Error("Release job failed!")
					d.core.DeliveryUpdate(d.deliveryId, "error")
					return
				} else if latestJob.Status.Active == 0 {
					jobLog.Info("Release job succeeded!")
					d.core.DeliveryUpdate(d.deliveryId, "deploying")
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
}
