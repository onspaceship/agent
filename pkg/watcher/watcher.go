package watcher

import (
	"context"
	"time"

	"github.com/apex/log"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Watcher struct {
	Client  *kubernetes.Clientset
	Factory informers.SharedInformerFactory
}

func NewWatcher() *Watcher {
	client := kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie())

	factory := informers.NewSharedInformerFactory(
		client,
		10*time.Minute,
	)
	return &Watcher{
		Client:  client,
		Factory: factory,
	}
}

func (w *Watcher) Start(ctx context.Context) {
	log.Info("Starting the Deployment watcher")

	informer := w.Factory.Apps().V1().Deployments().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: w.onUpdate,
	})

	informer.Run(ctx.Done())

	log.Info("Stopping the Deployment watcher")
}
