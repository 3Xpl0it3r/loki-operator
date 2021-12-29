// Copyright 2021 The Loki-operator Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package options

import (
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
)

type Options struct {
	ListenAddress string
}

var _ options = new(Options)

// NewOptions create an instance option and return
func NewOptions() *Options {
	return &Options{}
}

// Validate validates options
func (o *Options) Validate() []error {
	return nil
}

// Complete fill some default value to options
func (o *Options) Complete() error {
	return nil
}

//
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ListenAddress, "web.listen-addr", "127.0.0.1:8080", "Address on which to expose metrics and web interfaces")
}

func (o *Options) NamedFlagSets() (fs flag.NamedFlagSets) {
	o.AddFlags(fs.FlagSet("loki-operator"))
	// other options addFlags
	return
}
