// Copyright 2014 José Carlos Nieto
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

package main

import (
	"fmt"
)

func destroyFn(id string) []byte {
	s := fmt.Sprintf(`[{"fn": "destroy", "id": "%s"}]`, id)
	return []byte(s)
}

func identFn(id string) []byte {
	s := fmt.Sprintf(`[{"fn": "ident", "id": "%s"}]`, id)
	return []byte(s)
}

func scoresFn(data []byte) []byte {
	s := fmt.Sprintf(`[{"fn": "scores", "data": %s}]`, string(data))
	return []byte(s)
}

func updateFn(id string, data []byte) []byte {
	s := fmt.Sprintf(`[{"fn": "update", "id": "%s", "data": %s}]`, id, string(data))
	return []byte(s)
}

func createFn(t string, id string, data []byte) []byte {
	s := fmt.Sprintf(`[{"fn": "create", "kind": "%s", "id": "%s", "data": %s}]`, t, id, string(data))
	return []byte(s)
}
