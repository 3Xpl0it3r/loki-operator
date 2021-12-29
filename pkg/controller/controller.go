/*
Copyright 2021 The loki-operator Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

// Controller is generic interface for custom controller, it defines the basic behaviour of custom controller
type Controller interface {
	Start(ctx context.Context) error
	Stop()
	AddHook(hook Hook) error
	RemoveHook(hook Hook) error
}

// this is example, you should remove it in product
type emptyController struct {
}

func (e emptyController) Start(ctx context.Context) error {
	return nil
}
func (e emptyController) Stop() {
}
func (e emptyController) AddHook(hook Hook) error {
	return nil
}
func (e emptyController) RemoveHook(hook Hook) error {
	return nil
}
func NewEmptyController(reg prometheus.Registerer) Controller {
	return &emptyController{}
}
