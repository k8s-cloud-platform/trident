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

package options

import (
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/component-base/logs"
)

const (
	leaderElect                  = true
	leaderElectResourceNamespace = "default"
	leaderElectResourceName      = "manager.trident.kcp.io"
	leaderElectLeaseDuration     = 15 * time.Second
	leaderElectRenewDeadline     = 10 * time.Second
	leaderElectRetryPeriod       = 3 * time.Second

	ConcurrencyTenantSync  = 10
	ConcurrencyClusterSync = 10
)

type Options struct {
	EtcdServers            string
	EtcdSecret             string
	ConcurrencyTenantSync  int
	ConcurrencyClusterSync int

	Log            *logs.Options
	LeaderElection *componentbaseconfig.LeaderElectionConfiguration
}

func NewOptions() *Options {
	return &Options{
		Log: logs.NewOptions(),
		LeaderElection: &componentbaseconfig.LeaderElectionConfiguration{
			ResourceLock: resourcelock.LeasesResourceLock,
		},
	}
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)
	o.Log.AddFlags(flags)

	flags.StringVar(&o.EtcdServers, "etcd-servers", "",
		"Etcd servers, used for tenant apiserver connect to host etcd clusters, use ',' to separate.")
	flags.StringVar(&o.EtcdSecret, "etcd-secret", "",
		"Reference of etcd secret, use [namespace]/[name] or [name](use default namespace).")

	flags.IntVar(&o.ConcurrencyTenantSync, "concurrency-tenant-sync", ConcurrencyTenantSync,
		"Concurrency of tenant controllers to sync.")
	flags.IntVar(&o.ConcurrencyClusterSync, "concurrency-cluster-sync", ConcurrencyClusterSync,
		"Concurrency of cluster controllers to sync.")

	flags.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", leaderElect,
		"Enable leader elect.")
	flags.StringVar(&o.LeaderElection.ResourceNamespace, "leader-elect-resource-namespace", leaderElectResourceNamespace,
		"Namespace of leader elect resource.")
	flags.StringVar(&o.LeaderElection.ResourceName, "leader-elect-resource-name", leaderElectResourceName,
		"Name of leader elect resource.")
	flags.DurationVar(&o.LeaderElection.LeaseDuration.Duration, "leader-elect-lease-duration", leaderElectLeaseDuration,
		"Duration of leader elect lease.")
	flags.DurationVar(&o.LeaderElection.RenewDeadline.Duration, "leader-elect-renew-deadline", leaderElectRenewDeadline,
		"Duration of leader elect renew deadline.")
	flags.DurationVar(&o.LeaderElection.RetryPeriod.Duration, "leader-elect-retry-period", leaderElectRetryPeriod,
		"Duration of leader elect retry period.")
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() field.ErrorList {
	errs := field.ErrorList{}
	newPath := field.NewPath("Options")

	if o.LeaderElection.LeaseDuration.Duration <= 0 {
		errs = append(errs, field.Required(newPath.Child("LeaderElection.LeaseDuration.Duration"), "must bigger than 0"))
	}
	if o.LeaderElection.RenewDeadline.Duration <= 0 {
		errs = append(errs, field.Required(newPath.Child("LeaderElection.RenewDeadline.Duration"), "must bigger than 0"))
	}
	if o.LeaderElection.RetryPeriod.Duration <= 0 {
		errs = append(errs, field.Required(newPath.Child("LeaderElection.RetryPeriod.Duration"), "must bigger than 0"))
	}

	if o.ConcurrencyClusterSync <= 0 {
		errs = append(errs, field.Required(newPath.Child("Concurrey.Cluster.Sync"), "must bigger than 0"))
	}
	if o.ConcurrencyTenantSync <= 0 {
		errs = append(errs, field.Required(newPath.Child("Concurrey.Tenant.Sync"), "must bigger than 0"))
	}

	if o.EtcdServers == "" {
		errs = append(errs, field.Required(newPath.Child("EtcdServers"), "must not empty"))
	}
	if o.EtcdSecret == "" {
		errs = append(errs, field.Required(newPath.Child("EtcdSecret"), "must not empty"))
	}

	return errs
}
