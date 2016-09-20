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

package diff

import (
	"encoding/json"
	"reflect"
)

const maxIgnores = 12

type Diff struct {
	Ignore  map[string]bool
	current map[string]interface{}
	old     map[string]interface{}
	ignored int
	s       string
}

func NewDiff() *Diff {
	self := &Diff{}
	self.Ignore = map[string]bool{}
	self.current = map[string]interface{}{}
	self.old = map[string]interface{}{}
	self.s = `{}`
	return self
}

func (self *Diff) SetData(data *map[string]interface{}) {
	self.current = *data
}

func (self *Diff) MarshalJSON() ([]byte, error) {
	u := map[string]interface{}{}

	for k, _ := range self.current {
		if _, ok := self.old[k]; ok == false {
			// Key does not even exists.
			u[k] = self.current[k]
			self.old[k] = u[k]
		} else {
			// Key exists, is it equal?
			if reflect.DeepEqual(self.old[k], self.current[k]) == false {
				u[k] = self.current[k]
				self.old[k] = u[k]
			}
		}
	}

	ok := false

	for k, _ := range u {
		if _, e := self.Ignore[k]; e == true {
			if self.Ignore[k] == false {
				ok = true
			}
		} else {
			ok = true
		}
	}

	if ok == false {
		self.ignored = self.ignored + 1
		if self.ignored < maxIgnores {
			return nil, nil
		}
	}

	self.ignored = 0

	return json.Marshal(u)
}

func (self *Diff) Serialize() []byte {
	b, _ := json.Marshal(self)
	if b == nil {
		return nil
	}
	if len(b) > 0 {
		s := string(b)
		if s != `{}` {
			return b
		}
		self.s = s
	}
	return nil
}
