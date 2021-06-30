package server

import (
	"context"
	"github.com/keikoproj/ice-kube/internal/log"
	"github.com/keikoproj/ice-kube/pkg/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"strconv"
	"time"
)

const (
	suspendLabelKey   = "icekube.kubernetes.io/suspend"
	suspendLabelValue = "true"

	//Annotations
	suspendAnnotationKey = "icekube.kubernetes.io/after-mins"
)

type Server struct {
	K8sClient k8s.Interface
}

type Interface interface {
	SuspendWorkLoad(ctx context.Context) error
	CustomScheduleWorkload(ctx context.Context) error
}

//NewServerOrDie provides the client with inbuilt cluster client
func NewServerOrDie() *Server {
	return &Server{
		K8sClient: k8s.NewK8sSelfClientDoOrDie(),
	}
}

//NewServerOrDieFromConfig provides the client for provided rest config
func NewServerOrDieFromConfig(config *rest.Config) *Server {
	return &Server{
		K8sClient: k8s.NewK8sManagedClusterClientDoOrDie(config),
	}
}

//SuspendWorkLoad suspends the work load based on pre-defined annotations
func (s *Server) SuspendWorkLoad(ctx context.Context) error {
	log := log.Logger(ctx)
	log.Debug("Start SuspendWorkLoad")

	//List all the deployments with pre-defined labels
	//Process each deployment
	//  i. compare with annotation to act
	replList, err := s.K8sClient.ListV1ReplicaSets(ctx, suspendLabelKey, suspendLabelValue)
	if err != nil {
		log.WithField("error", err.Error()).Error("Unable to get the deployment list")
		return err
	}

	//print the replicaset names and creation times
	for _, repl := range replList.Items {
		log.WithFields(logrus.Fields{
			"replicaset_name":      repl.Name,
			"deployment_name":      repl.OwnerReferences[0].Name,
			"deployment_namespace": repl.Namespace,
			"replica_count":        repl.Spec.Replicas,
			"creation_time":        repl.CreationTimestamp,
			"after-mins-value":     repl.Annotations[suspendAnnotationKey],
		}).Info("successfully fetched the replica sets")

		//First compare the time stamps
		duration := time.Since(repl.CreationTimestamp.Time).Minutes()
		meltTime, _ := strconv.Atoi(repl.Annotations[suspendAnnotationKey])
		if duration > float64(meltTime) {
			log.WithFields(logrus.Fields{
				"time_since": duration,
				"melt_time":  meltTime,
			}).Info("Threshold exceeded")

			// Lets check the flag here
			// Get the deployment
			depl, err := s.K8sClient.GetV1Deployment(ctx, repl.Namespace, repl.OwnerReferences[0].Name)
			if err != nil {
				log.WithField("error", err.Error()).Error("unable to get the deployments")
				return err
			}
			if err := s.K8sClient.ScaleV1Deployments(ctx, depl, 0, "true"); err != nil {
				log.WithField("error", err.Error()).Error("unable to scale the deployments")
				return err
			}
		}

		//Lets scale the resources

	}
	//  ii. scale down replicas

	return nil
}

//CustomScheduleWorkload function schedules the work load based on annotations
func (s *Server) CustomScheduleWorkload(ctx context.Context) error {

	return nil
}
