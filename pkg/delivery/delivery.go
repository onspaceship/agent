package delivery

import (
	"context"

	"github.com/onspaceship/agent/pkg/client"

	"github.com/apex/log"
	"k8s.io/client-go/kubernetes"
)

type delivery struct {
	ctx context.Context
	log *log.Entry

	core       *client.Client
	kubernetes *kubernetes.Clientset

	deliveryId string
	appId      string
	appHandle  string
	teamHandle string
	imageURI   string
}

func NewDelivery(deliveryId string, appId string, appHandle string, teamHandle string, imageURI string, clientset *kubernetes.Clientset) *delivery {
	return &delivery{
		ctx: context.Background(),
		log: log.WithField("app_id", appId).WithField("delivery_id", deliveryId),

		core:       client.NewClient(),
		kubernetes: clientset,

		deliveryId: deliveryId,
		appId:      appId,
		appHandle:  appHandle,
		teamHandle: teamHandle,
		imageURI:   imageURI,
	}
}

func (d *delivery) Deliver() {
	d.log.Info("Delivering new app version")

	d.runReleaseJob()
	d.updateDeployment()
}
