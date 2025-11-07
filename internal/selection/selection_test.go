package selection

import "testing"

func TestSelectContent_NilManager(t *testing.T) {
	res, err := SelectContent(nil, 80, 24)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil selection when manager is nil")
	}
}
