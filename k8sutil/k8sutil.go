// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package k8sutil

import (
	"errors"
	"os"
	"path/filepath"

	"context"
	"fmt"

	istioClient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Helper function to obtain the default kubeConfig location
func GetKubeConfigLocation() (string, error) {

	var kubeConfig string
	kubeConfigEnvVar := os.Getenv("KUBECONFIG")

	if len(kubeConfigEnvVar) > 0 {
		// Find using environment variables
		kubeConfig = kubeConfigEnvVar
	} else if home := homedir.HomeDir(); home != "" {
		// Find in the ~/.kube/ directory
		kubeConfig = filepath.Join(home, ".kube", "config")
	} else {
		// give up
		return kubeConfig, errors.New("unable to find kubeconfig")
	}
	return kubeConfig, nil
}

// GetKubernetesClientset returns the Kubernetes clientset for the cluster set in the environment
func GetKubernetesClientset() (*kubernetes.Clientset, error) {
	// use the current context in the kubeconfig
	var clientset *kubernetes.Clientset
	config, err := GetKubeConfig()
	if err != nil {
		return clientset, err
	}
	clientset, err = kubernetes.NewForConfig(config)
	return clientset, err
}

// Returns kubeconfig from KUBECONFIG env var if set
// Else from default location ~/.kube/config
func GetKubeConfig() (*rest.Config, error) {
	var config *rest.Config
	kubeConfigLoc, err := GetKubeConfigLocation()
	if err != nil {
		return config, err
	}
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigLoc)
	return config, err
}

// GetIstioClientset returns the clientset object for Istio
func GetIstioClientset() (*istioClient.Clientset, error) {
	var cs *istioClient.Clientset
	kubeConfig, err := GetKubeConfig()
	if err != nil {
		return cs, err
	}
	cs, err = istioClient.NewForConfig(kubeConfig)
	return cs, err
}

func GetHostnameFromGateway(namespace string, appConfigName string) (string, error) {
	cs, err := GetIstioClientset()
	if err != nil {
		return "", err
	}
	gateways, err := cs.NetworkingV1alpha3().Gateways(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	// if an optional appConfigName is provided, construct the gateway name from the namespace and
	// appConfigName and look for that specific gateway, otherwise just use the first gateway
	gatewayName := ""
	if len(appConfigName) > 0 {
		gatewayName = fmt.Sprintf("%s-%s-gw", namespace, appConfigName)
	}

	for _, gateway := range gateways.Items {
		if len(gatewayName) > 0 && gatewayName != gateway.ObjectMeta.Name {
			continue
		}
		if len(gateway.Spec.Servers) > 0 && len(gateway.Spec.Servers[0].Hosts) > 0 {
			return gateway.Spec.Servers[0].Hosts[0], nil
		}
	}

	// this can happen if the app gateway has not been created yet, the caller should
	// keep retrying and eventually we should get a gateway with a host
	return "", nil
}
