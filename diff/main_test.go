package diff

import (
	"testing"
)

func TestData(t *testing.T) {
	var buf []byte

	diff := NewDiff()

	buf = diff.Serialize()
	if len(buf) > 0 {
		t.Fatalf("Expecting empty serialization.")
	}

	diff.SetData(map[string]interface{}{
		"foo": 123.45,
		"bar": 45,
	})

	buf = diff.Serialize()
	if len(buf) == 0 {
		t.Fatalf("Expecting not empty serialization.")
	}

	buf = diff.Serialize()
	if len(buf) > 0 {
		t.Fatalf("Expecting empty serialization.")
	}

	diff.SetData(map[string]interface{}{
		"foo": 123.45,
		"bar": 45,
	})

	buf = diff.Serialize()
	if len(buf) > 0 {
		t.Fatalf("Expecting empty serialization.")
	}

	diff.SetData(map[string]interface{}{
		"foo": 123.46,
		"bar": 45,
	})

	buf = diff.Serialize()
	if len(buf) == 0 {
		t.Fatalf("Expecting not empty serialization.")
	}

}
