package watcher

import (
	"github.com/apex/log"
	"github.com/onspaceship/agent/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
)

func (w *Watcher) onUpdate(_, obj interface{}) {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return
	}

	appId, ok := deployment.Labels[config.AppIdLabel]
	if !ok {
		return
	}

	log.WithField("deployment", deployment.Name).WithField("appId", appId).Info("updated!")
}
