package k8s

import (
	"context"
	"fmt"
	"github.com/keikoproj/ice-kube/internal/log"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//ListDeployments functions lists all the deployments with
func (k *Client) ListV1Deployments(ctx context.Context, labelKey string, labelValue string) (*v1.DeploymentList, error) {
	log := log.Logger(ctx)
	log.Debug("Start ListV1Deployments")
	opts := client.MatchingLabels{
		labelKey: labelValue,
	}

	depList := v1.DeploymentList{}

	if err := k.runtimeClient.List(ctx, &depList, opts); err != nil {
		log.WithField("error", err.Error()).Error("Error in listing the deployments")
		return nil, err
	}
	log.WithField("count", len(depList.Items)).Info("successfully listed the v1 deployments")

	return &depList, nil
}

//GetV1Deployment function gets the namespace
func (k *Client) GetV1Deployment(ctx context.Context, ns string, name string) (*v1.Deployment, error) {
	log := log.Logger(ctx)
	log.Debug("Start ListV1Deployments")

	depl := v1.Deployment{}
	if err := k.runtimeClient.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: ns,
	}, &depl); err != nil {
		log.WithField("error", err.Error()).Error("Error in getting the deployment")
		return nil, err
	}

	return &depl, nil
}

//ScaleDownV1Deployments scales down to the requested replica count
func (k *Client) ScaleV1Deployments(ctx context.Context, depl *v1.Deployment, replicas int) error {
	log := log.Logger(ctx)
	log.Debug("Start ScaleDownV1Deployments")

	//Prepare the patch to make replica count to desired
	patchStr := fmt.Sprintf(`{"spec":{"replicas": %d }}`, replicas)
	if err := k.runtimeClient.Patch(context.Background(), &v1.Deployment{
		ObjectMeta: depl.ObjectMeta,
	}, client.RawPatch(types.StrategicMergePatchType, []byte(patchStr))); err != nil {
		log.WithField("error", err.Error()).Error("Error in scaling deployment replicas")
		return err
	}

	log.WithFields(logrus.Fields{
		"desired_count":        replicas,
		"deployment_name":      depl.Name,
		"deployment_namespace": depl.Namespace,
	}).Info("Successfully scaled the v1 deployments")
	return nil
}

//ListDeployments functions lists all the deployments with
func (k *Client) ListV1Beta2Deployments(ctx context.Context, labelKey string, labelValue string) (*v1beta2.DeploymentList, error) {
	log := log.Logger(ctx)
	log.Debug("Start ListV1Beta2Deployments")
	opts := client.MatchingLabels{
		labelKey: labelValue,
	}

	depList := v1beta2.DeploymentList{}

	if err := k.runtimeClient.List(ctx, &depList, opts); err != nil {
		log.WithField("error", err.Error()).Error("Error in listing the deployments")
		return nil, err
	}
	log.WithField("count", len(depList.Items)).Info("successfully listed the v1beta2 deployments")

	return &depList, nil
}

//ScaleDownV1Deployments scales down to the requested replica count
func (k *Client) ScaleV1Beta2Deployments(ctx context.Context, depl *v1beta2.Deployment, replicas int) error {
	log := log.Logger(ctx)
	log.Debug("Start ScaleV1Beta2Deployments")

	//Prepare the patch to make replica count to desired
	patchStr := fmt.Sprintf(`{"spec":{"replicas": %d }}`, replicas)
	if err := k.runtimeClient.Patch(context.Background(), &v1beta2.Deployment{
		ObjectMeta: depl.ObjectMeta,
	}, client.RawPatch(types.StrategicMergePatchType, []byte(patchStr))); err != nil {
		log.WithField("error", err.Error()).Error("Error in scaling deployment replicas")
		return err
	}

	log.WithFields(logrus.Fields{
		"desired_count":        replicas,
		"deployment_name":      depl.Name,
		"deployment_namespace": depl.Namespace,
	}).Info("Successfully scaled the v1Beta2 deployments")

	return nil
}
