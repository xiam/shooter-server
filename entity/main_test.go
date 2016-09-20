package entity

import (
	"fmt"
	"testing"
)

func TestSerialize(t *testing.T) {
	var b []byte

	e := NewEntity()
	e.SetPosition(100.0, 100.0)
	e.SetDirection(1.0, -1.0)
	e.SetSpeed(3.1)

	for i := 0; i < 5; i++ {

		b = e.Serialize()
		fmt.Printf("json: %s\n", b)

		b = e.Serialize()
		fmt.Printf("json: %s\n", b)

		b = e.Serialize()
		fmt.Printf("json: %s\n", b)

		b = e.Serialize()
		fmt.Printf("json: %s\n", b)

		e.Tick()
	}
}

func TestPoly(t *testing.T) {
	e := NewEntity()
	e.SetPosition(100.0, 100.0)
	e.SetDirection(-1, 0)
	e.Width = 30
	e.Height = 80
	p := e.Poly()

	for i := 0; i < p.Len; i++ {
		fmt.Printf("poly :%v\n", p.Points[i])
	}

}
