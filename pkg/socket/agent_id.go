package socket

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	AgentIdAnnotation = "onspaceship.com/agent-id"
)

func (socket *socket) ensureAgentId() {
	if socket.AgentId != "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client := kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie())

	deployment, err := client.AppsV1().Deployments(socket.Namespace).Get(ctx, "agent-connect", metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Fatal("Could not get Kubernetes deployment for Agent")
	}

	agentId := deployment.Annotations[AgentIdAnnotation]

	if agentId == "" {
		agentId = uuid.NewString()
		log.Infof("Generating a new Agent ID: %v", agentId)

		deployment.Annotations[AgentIdAnnotation] = agentId
		_, err = client.AppsV1().Deployments(socket.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			log.WithError(err).Fatal("Could not store new Agent ID in Kubernetes")
		}
	}

	socket.AgentId = agentId
}
