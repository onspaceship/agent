package update

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/onspaceship/agent/pkg/client"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckForUpdate(k8s *kubernetes.Clientset) error {
	core := client.NewClient()

	version, err := core.GetVersion()
	if err != nil {
		return err
	}

	return ProcessVersionUpdate(version.Version, k8s)
}

func ProcessVersionUpdate(version string, client *kubernetes.Clientset) error {
	log.WithField("version", version).Info("Updating Agent version")

	name := os.Getenv("POD_NAME")
	namespace := os.Getenv("POD_NAMESPACE")

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to get Agent pod")
		} else {
			return err
		}
	}

	if len(pod.OwnerReferences) != 1 || pod.OwnerReferences[0].Kind != "ReplicaSet" {
		return errors.New("unable to find Agent pod's replica set")
	}

	rs, err := client.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), pod.OwnerReferences[0].Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to get Agent replica set")
		} else {
			return err
		}
	}

	if len(rs.OwnerReferences) != 1 || rs.OwnerReferences[0].Kind != "Deployment" {
		return errors.New("unable to find Agent replica set's deployment")
	}

	deployment, err := client.AppsV1().Deployments(pod.Namespace).Get(context.TODO(), rs.OwnerReferences[0].Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to get Agent deployment")
		} else {
			return err
		}
	}

	imageURI := deployment.Spec.Template.Spec.Containers[0].Image
	imageURI = strings.Split(imageURI, ":")[0]
	imageURI = strings.Split(imageURI, "@sha256")[0]

	deployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", imageURI, version)

	_, err = client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.WithField("deployment", deployment.Name).Info("Agent image updated")
	return nil
}
