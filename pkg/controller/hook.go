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

// Hook is interface for hooks that can be inject into custom controller
type Hook interface {
	// OnAdd runs after the controller finished processing the addObject
	OnAdd(object interface{})
	// OnUpdate runs after the controller finished processing the updatedObject
	OnUpdate(object interface{})
	// OnDelete run after the controller finished processing the deletedObject
	OnDelete(object interface{})
}
