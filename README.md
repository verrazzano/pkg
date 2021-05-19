# Kubernetes Golang Utilities for Verrazzano

This repository is a collection of general-purpose, Kubernetes-related Golang utilities used across
[Verrazzano](https://github.com/verrazzano).

## How to use

To import within your project, run:
```
go get github.com/verrazzano/pkg
```

See below for specific examples of using the various packages.

## Packages

### ./diff: Kubernetes Object Comparator

The utilities in this package are intended to be in used in the context of custom Kubernetes controllers, 
when comparing the live state of a Kubernetes object (fetched via API) against a desired state (constructed
via code).  The "live" object, fetched from Kubernetes, will auto-generate many fields that we won't
or can't specify in the "desired" object that we construct.  A common job of a Kubernetes controller is to 
compare a desired and live object, and `Apply()` changes if they are different.  But we don't want these 
auto-generated fields to cause us to apply changes, otherwise we'd be doing so frequently.  A standard Golang 
struct comparison _will_ detect these as differences, and is therefore insufficient for our purposes here.

Our `Diff()` function simply performs a recursive struct comparison, but ignores any fields that were unspecified 
in the `toObject` (the desired object):

```
import (
        "github.com/verrazzano/pkg/diff"
)
```

Example: Detecting changes to a deployment object:

Fetch the live Kubernetes object:
```
cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
log.Fatal(err)
}
kubeClientSet, err := kubernetes.NewForConfig(cfg)
if err != nil {
log.Fatal(err)
}
kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, 30*time.Second)
deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
deploymentLister := deploymentInformer.Lister()
liveDeployment := deploymentsLister.Deployments("mynamespace").Get("mydeployment")
```

Construct your desired object:
```
desiredDeployment := createMyDeploymentObject()
```

Compare:
```
diffs := diff.Diff(liveDeployment, desiredDeployment)
if diffs != "" {
	fmt.Println("Diffs: " + diffs)
}
```