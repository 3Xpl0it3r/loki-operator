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
	"errors"
)

type Base struct {
	hooks []Hook
}

func NewControllerBase() Base {
	return Base{hooks: make([]Hook, 0)}
}

func (c *Base) GetHooks() []Hook {
	return c.hooks
}

func (c *Base) AddHook(hook Hook) error {
	for _, h := range c.hooks {
		if h == hook {
			return errors.New("Given hook is already installed in the current controller ")
		}
	}
	c.hooks = append(c.hooks, hook)
	return nil
}

func (c *Base) RemoveHook(hook Hook) error {
	for i, h := range c.hooks {
		if h == hook {
			c.hooks = append(c.hooks[:i], c.hooks[i+1:]...)
			return nil
		}
	}
	return errors.New("Given hook is not installed in the current controller ")
}
