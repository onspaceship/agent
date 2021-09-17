package watcher

import (
	"github.com/apex/log"
	"github.com/onspaceship/agent/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
)

func (w *Watcher) onUpdate(_, obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	appId := deployment.Labels[config.AppIdLabel]

	log.WithField("deployment", deployment.Name).WithField("appId", appId).Info("updated!")
}
