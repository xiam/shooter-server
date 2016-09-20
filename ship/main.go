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

package ship

import (
	"github.com/xiam/shooter-server/entity"
)

type Ship struct {
	*entity.Entity
}

func NewShip() *Ship {
	self := &Ship{}
	self.Entity = entity.NewEntity()
	//self.Entity.Diff.Ignore["p"] = true
	self.Entity.Width = 75
	self.Entity.Height = 112
	return self
}
