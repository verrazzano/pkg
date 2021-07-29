// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.
package k8sutil_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verrazzano/pkg/k8sutil"
	"k8s.io/client-go/util/homedir"
)

func TestGetKubeConfigLocation(t *testing.T) {
	asserts := assert.New(t)
	// Test without environment variable
	kubeConfigLoc, err := k8sutil.GetKubeConfigLocation()
	asserts.NoError(err)
	asserts.Equal(kubeConfigLoc, homedir.HomeDir()+"/.kube/config")
	// Test using environment variable
	err = os.Setenv("KUBECONFIG", "/home/xyz/somerandompath")
	asserts.NoError(err)
	kubeConfigLoc, err = k8sutil.GetKubeConfigLocation()
	asserts.NoError(err)
	asserts.Equal("/home/xyz/somerandompath", kubeConfigLoc)
}
