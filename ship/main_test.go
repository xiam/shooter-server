package ship

import (
	"testing"
)

func TestNewShip(t *testing.T) {
	ship := NewShip()
	if ship == nil {
		t.Fatalf("Failed!")
	}
}
