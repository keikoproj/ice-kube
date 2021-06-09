package k8s

import (
	"context"
	"fmt"
	"k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	runtimeClient client.Client
}

//NewK8sSelfClientDoOrDie gets the new k8s go client
func NewK8sSelfClientDoOrDie() *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("THIS IS LOCAL")
		// Do i need to panic here?
		//How do i test this from local?
		//Lets get it from local config file
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	}

	dClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}

	k8sCl := &Client{
		runtimeClient: dClient,
	}
	return k8sCl
}

//NewK8sManagedClusterClientDoOrDie creates a client for managed cluster or config passed
func NewK8sManagedClusterClientDoOrDie(config *rest.Config) *Client {

	//https://godoc.org/sigs.k8s.io/controller-runtime/pkg/client#New
	dClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}

	k8sCl := &Client{
		runtimeClient: dClient,
	}

	return k8sCl
}

type Interface interface {
	ListV1ReplicaSets(ctx context.Context, labelKey string, labelValue string) (*v1.ReplicaSetList, error)
	GetV1Deployment(ctx context.Context, ns string, name string) (*v1.Deployment, error)
	ListV1Deployments(ctx context.Context, labelKey string, labelValue string) (*v1.DeploymentList, error)
	ListV1Beta2Deployments(ctx context.Context, labelKey string, labelValue string) (*v1beta2.DeploymentList, error)
	ScaleV1Deployments(ctx context.Context, depl *v1.Deployment, replicas int) error
	ScaleV1Beta2Deployments(ctx context.Context, depl *v1beta2.Deployment, replicas int) error
	GetRecentReplicaSet(ctx context.Context, owner string, ns string) (*v1.ReplicaSet, error)
	ListStateFulSets(ctx context.Context, labelKey string, labelValue string)
}
