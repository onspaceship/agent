package watcher

import (
	"context"
	"time"

	"github.com/onspaceship/agent/pkg/client"
	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Watcher struct {
	Client  *kubernetes.Clientset
	Factory informers.SharedInformerFactory
	Core    *client.Client
}

func NewWatcher() *Watcher {
	k8sclient := kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie())

	factory := informers.NewSharedInformerFactoryWithOptions(
		k8sclient,
		10*time.Minute,
		informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.LabelSelector = config.AppIdLabel
		}),
	)

	return &Watcher{
		Client:  k8sclient,
		Factory: factory,
		Core:    client.NewClient(),
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
