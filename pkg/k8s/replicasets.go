package k8s

import (
	"context"
	"github.com/keikoproj/ice-kube/internal/log"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func (k *Client) GetRecentReplicaSet(ctx context.Context, owner string, ns string) (*v1.ReplicaSet, error) {
	log := log.Logger(ctx)
	log.Debug("Start GetRecentReplicaset")

	opts := client.InNamespace(ns)

	repList := v1.ReplicaSetList{}

	if err := k.runtimeClient.List(ctx, &repList, opts); err != nil {
		log.WithField("error", err.Error()).Error("Error in listing the v1 replica sets")
		return nil, err
	}
	log.WithField("count", len(repList.Items)).Info("successfully listed the v1 replica sets")

	recent := -1
	for i, rep := range repList.Items {
		if rep.OwnerReferences[0].Name == owner {
			log.WithFields(logrus.Fields{
				"revision":   rep.Annotations["deployment.kubernetes.io/revision"],
				"time_stamp": rep.CreationTimestamp,
				"i":          i,
			}).Info("replica sets list")
			revision, _ := strconv.Atoi(rep.Annotations["deployment.kubernetes.io/revision"])
			if recent < revision {
				recent = i
			}
		}
	}

	log.WithField("recent_one", recent).Info("recent found")

	return &repList.Items[recent], nil
}

type Replica struct {
	Revision int
	Index    int
}

//ListV1ReplicaSets functions lists all the deployments with
func (k *Client) ListV1ReplicaSets(ctx context.Context, labelKey string, labelValue string) (*v1.ReplicaSetList, error) {
	log := log.Logger(ctx)
	log.Debug("Start ListV1ReplicaSets")
	opts := client.MatchingLabels{
		labelKey: labelValue,
	}

	repList := v1.ReplicaSetList{}

	if err := k.runtimeClient.List(ctx, &repList, opts); err != nil {
		log.WithField("error", err.Error()).Error("Error in listing the replica sets")
		return nil, err
	}
	log.WithField("count", len(repList.Items)).Info("successfully listed the v1 replica sets")
	unique := make(map[string]Replica)
	for i, repl := range repList.Items {
		revision, _ := strconv.Atoi(repl.Annotations["deployment.kubernetes.io/revision"])
		for _, own := range repl.OwnerReferences {
			if repl.OwnerReferences[0].Kind == "Deployment" && *repl.Spec.Replicas != 0 {
				if val, ok := unique[own.Name]; ok {
					if val.Revision < revision {
						// Revision comparision for existing
						unique[own.Name] = Replica{
							Index:    i,
							Revision: revision,
						}
					}
				} else {
					// First time use case
					unique[own.Name] = Replica{
						Index:    i,
						Revision: revision,
					}
				}
			}

		}
	}

	// Prepare the final list
	var items []v1.ReplicaSet
	for _, v := range unique {
		items = append(items, repList.Items[v.Index])
	}
	repList.Items = items
	return &repList, nil
}
