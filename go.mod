// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

module github.com/verrazzano/pkg

go 1.15

require (
	github.com/google/go-cmp v0.5.0
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.21.1
)

replace (
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
)
