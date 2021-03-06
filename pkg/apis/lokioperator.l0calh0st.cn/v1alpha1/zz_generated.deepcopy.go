//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChunkStoreConfig) DeepCopyInto(out *ChunkStoreConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChunkStoreConfig.
func (in *ChunkStoreConfig) DeepCopy() *ChunkStoreConfig {
	if in == nil {
		return nil
	}
	out := new(ChunkStoreConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonConfig) DeepCopyInto(out *CommonConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonConfig.
func (in *CommonConfig) DeepCopy() *CommonConfig {
	if in == nil {
		return nil
	}
	out := new(CommonConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CompactorConfig) DeepCopyInto(out *CompactorConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CompactorConfig.
func (in *CompactorConfig) DeepCopy() *CompactorConfig {
	if in == nil {
		return nil
	}
	out := new(CompactorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DistributorConfig) DeepCopyInto(out *DistributorConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DistributorConfig.
func (in *DistributorConfig) DeepCopy() *DistributorConfig {
	if in == nil {
		return nil
	}
	out := new(DistributorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendWorkerConfig) DeepCopyInto(out *FrontendWorkerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendWorkerConfig.
func (in *FrontendWorkerConfig) DeepCopy() *FrontendWorkerConfig {
	if in == nil {
		return nil
	}
	out := new(FrontendWorkerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngesterClientConfig) DeepCopyInto(out *IngesterClientConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngesterClientConfig.
func (in *IngesterClientConfig) DeepCopy() *IngesterClientConfig {
	if in == nil {
		return nil
	}
	out := new(IngesterClientConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngesterConfig) DeepCopyInto(out *IngesterConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngesterConfig.
func (in *IngesterConfig) DeepCopy() *IngesterConfig {
	if in == nil {
		return nil
	}
	out := new(IngesterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LimitsConfig) DeepCopyInto(out *LimitsConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimitsConfig.
func (in *LimitsConfig) DeepCopy() *LimitsConfig {
	if in == nil {
		return nil
	}
	out := new(LimitsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Loki) DeepCopyInto(out *Loki) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Loki.
func (in *Loki) DeepCopy() *Loki {
	if in == nil {
		return nil
	}
	out := new(Loki)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Loki) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiList) DeepCopyInto(out *LokiList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Loki, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiList.
func (in *LokiList) DeepCopy() *LokiList {
	if in == nil {
		return nil
	}
	out := new(LokiList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LokiList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiServerConfig) DeepCopyInto(out *LokiServerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiServerConfig.
func (in *LokiServerConfig) DeepCopy() *LokiServerConfig {
	if in == nil {
		return nil
	}
	out := new(LokiServerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiSpec) DeepCopyInto(out *LokiSpec) {
	*out = *in
	if in.DeployMode != nil {
		in, out := &in.DeployMode, &out.DeployMode
		*out = make(map[ModeKind]map[LokiTargetKind]int32, len(*in))
		for key, val := range *in {
			var outVal map[LokiTargetKind]int32
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(map[LokiTargetKind]int32, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
	out.Config = in.Config
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiSpec.
func (in *LokiSpec) DeepCopy() *LokiSpec {
	if in == nil {
		return nil
	}
	out := new(LokiSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStatus) DeepCopyInto(out *LokiStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStatus.
func (in *LokiStatus) DeepCopy() *LokiStatus {
	if in == nil {
		return nil
	}
	out := new(LokiStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Promtail) DeepCopyInto(out *Promtail) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Promtail.
func (in *Promtail) DeepCopy() *Promtail {
	if in == nil {
		return nil
	}
	out := new(Promtail)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Promtail) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailClientConfig) DeepCopyInto(out *PromtailClientConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailClientConfig.
func (in *PromtailClientConfig) DeepCopy() *PromtailClientConfig {
	if in == nil {
		return nil
	}
	out := new(PromtailClientConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailConfig) DeepCopyInto(out *PromtailConfig) {
	*out = *in
	out.Server = in.Server
	out.Clients = in.Clients
	out.Positions = in.Positions
	out.ScrapeConfigs = in.ScrapeConfigs
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailConfig.
func (in *PromtailConfig) DeepCopy() *PromtailConfig {
	if in == nil {
		return nil
	}
	out := new(PromtailConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailList) DeepCopyInto(out *PromtailList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Promtail, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailList.
func (in *PromtailList) DeepCopy() *PromtailList {
	if in == nil {
		return nil
	}
	out := new(PromtailList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PromtailList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailPositionsConfig) DeepCopyInto(out *PromtailPositionsConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailPositionsConfig.
func (in *PromtailPositionsConfig) DeepCopy() *PromtailPositionsConfig {
	if in == nil {
		return nil
	}
	out := new(PromtailPositionsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailServerConfig) DeepCopyInto(out *PromtailServerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailServerConfig.
func (in *PromtailServerConfig) DeepCopy() *PromtailServerConfig {
	if in == nil {
		return nil
	}
	out := new(PromtailServerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailSpec) DeepCopyInto(out *PromtailSpec) {
	*out = *in
	out.Config = in.Config
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailSpec.
func (in *PromtailSpec) DeepCopy() *PromtailSpec {
	if in == nil {
		return nil
	}
	out := new(PromtailSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PromtailStatus) DeepCopyInto(out *PromtailStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromtailStatus.
func (in *PromtailStatus) DeepCopy() *PromtailStatus {
	if in == nil {
		return nil
	}
	out := new(PromtailStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QuerierConfig) DeepCopyInto(out *QuerierConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QuerierConfig.
func (in *QuerierConfig) DeepCopy() *QuerierConfig {
	if in == nil {
		return nil
	}
	out := new(QuerierConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryFrontendConfig) DeepCopyInto(out *QueryFrontendConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryFrontendConfig.
func (in *QueryFrontendConfig) DeepCopy() *QueryFrontendConfig {
	if in == nil {
		return nil
	}
	out := new(QueryFrontendConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QuerySchedulerConfig) DeepCopyInto(out *QuerySchedulerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QuerySchedulerConfig.
func (in *QuerySchedulerConfig) DeepCopy() *QuerySchedulerConfig {
	if in == nil {
		return nil
	}
	out := new(QuerySchedulerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryrangeConfig) DeepCopyInto(out *QueryrangeConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryrangeConfig.
func (in *QueryrangeConfig) DeepCopy() *QueryrangeConfig {
	if in == nil {
		return nil
	}
	out := new(QueryrangeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RulerConfig) DeepCopyInto(out *RulerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RulerConfig.
func (in *RulerConfig) DeepCopy() *RulerConfig {
	if in == nil {
		return nil
	}
	out := new(RulerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeConfig) DeepCopyInto(out *RuntimeConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeConfig.
func (in *RuntimeConfig) DeepCopy() *RuntimeConfig {
	if in == nil {
		return nil
	}
	out := new(RuntimeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchemaConfig) DeepCopyInto(out *SchemaConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchemaConfig.
func (in *SchemaConfig) DeepCopy() *SchemaConfig {
	if in == nil {
		return nil
	}
	out := new(SchemaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScrapeConfigs) DeepCopyInto(out *ScrapeConfigs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScrapeConfigs.
func (in *ScrapeConfigs) DeepCopy() *ScrapeConfigs {
	if in == nil {
		return nil
	}
	out := new(ScrapeConfigs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageConfig) DeepCopyInto(out *StorageConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageConfig.
func (in *StorageConfig) DeepCopy() *StorageConfig {
	if in == nil {
		return nil
	}
	out := new(StorageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TableManagerConfig) DeepCopyInto(out *TableManagerConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TableManagerConfig.
func (in *TableManagerConfig) DeepCopy() *TableManagerConfig {
	if in == nil {
		return nil
	}
	out := new(TableManagerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetConfig) DeepCopyInto(out *TargetConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetConfig.
func (in *TargetConfig) DeepCopy() *TargetConfig {
	if in == nil {
		return nil
	}
	out := new(TargetConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracingConfig) DeepCopyInto(out *TracingConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracingConfig.
func (in *TracingConfig) DeepCopy() *TracingConfig {
	if in == nil {
		return nil
	}
	out := new(TracingConfig)
	in.DeepCopyInto(out)
	return out
}
