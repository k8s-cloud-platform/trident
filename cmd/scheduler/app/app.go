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

package app

import (
	"context"
	"flag"
	"runtime/debug"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/k8s-cloud-platform/trident/cmd/scheduler/app/options"
	"github.com/k8s-cloud-platform/trident/pkg/apis/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

// NewSchedulerCommand creates a *cobra.Command object with default parameters
func NewSchedulerCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "scheduler",
		Long: `KCP scheduler for trident.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Log.ValidateAndApply(); err != nil {
				return err
			}

			cliflag.PrintFlags(cmd.Flags())
			buildInfo, ok := debug.ReadBuildInfo()
			if ok {
				klog.Infof("build info: \n%s", buildInfo)
			}

			if errs := opts.Validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}

			return run(ctrl.SetupSignalHandler(), opts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)

	return cmd
}

func run(ctx context.Context, opts *options.Options) error {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                        scheme,
		LeaderElection:                opts.LeaderElection.LeaderElect,
		LeaderElectionReleaseOnCancel: true,
		LeaderElectionResourceLock:    opts.LeaderElection.ResourceLock,
		LeaderElectionNamespace:       opts.LeaderElection.ResourceNamespace,
		LeaderElectionID:              opts.LeaderElection.ResourceName,
		LeaseDuration:                 &opts.LeaderElection.LeaseDuration.Duration,
		RenewDeadline:                 &opts.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:                   &opts.LeaderElection.RetryPeriod.Duration,
	})
	if err != nil {
		klog.ErrorS(err, "unable to start manager")
		return err
	}

	klog.Info("starting scheduler")
	if err := mgr.Start(ctx); err != nil {
		klog.ErrorS(err, "unable to run scheduler")
		return err
	}

	// never reach here
	return nil
}
