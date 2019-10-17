/*
Copyright 2014 Workiva, LLC

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
package palm

import "github.com/Workiva/go-datastructures/common"

type mockKey int

func (mk mockKey) Compare(other common.Comparator) int {
	otherKey := other.(mockKey)

	if mk == otherKey {
		return 0
	}

	if mk > otherKey {
		return 1
	}

	return -1
}
