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

package item

import (
	"github.com/xiam/shooter-server/entity"
)

type Item struct {
	*entity.Entity
}

func NewItem() *Item {
	self := &Item{}
	self.Entity = entity.NewEntity()
	self.Width = 34.0
	self.Height = 33.0
	self.SetDirection(1.0, 0.0)
	self.Life = 1
	return self
}
