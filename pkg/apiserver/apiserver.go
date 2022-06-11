/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(scheme)
)

func init() {
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// Server contains state for a Kubernetes cluster master/api server.
type Server struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

func (server *Server) Run(ctx context.Context) error {
	//mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
	//	Scheme:                 scheme,
	//	MetricsBindAddress:     "0",
	//	HealthProbeBindAddress: "0",
	//  SyncPeriod:             pointer.Duration(0),
	//})
	//if err != nil {
	//	klog.ErrorS(err, "unable to start apiserver")
	//	os.Exit(1)
	//}

	// TODO add post start hook
	server.GenericAPIServer.AddPostStartHookOrDie("test", func(context genericapiserver.PostStartHookContext) error {
		klog.Info("hook for post start...")
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(ctx.Done())
}
