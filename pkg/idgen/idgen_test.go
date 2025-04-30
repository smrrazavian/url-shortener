package idgen

import (
	"testing"
)

func TestGenerateId(t *testing.T) {

	id, err := GenerateID()
	if err != nil {
		t.Fatalf("Failed to generate id: %v", err)
	}
	if len(id) != IdLength {
		t.Errorf("ID length is %d, expected %d", len(id), IdLength)
	}

	const iterations = 1000
	ids := make(map[string]struct{}, iterations)

	for i := 0; i < iterations; i++ {
		id, err := GenerateID()
		if err != nil {
			t.Fatalf("Failed to generate ID at iteration %d: %v", i, err)
		}

		if _, exists := ids[id]; exists {
			t.Fatalf("Duplicate ID generated: %s", id)
		}
		ids[id] = struct{}{}
	}
}
