package rag

import "testing"

func TestSplitText(t *testing.T) {
	text := "abcdefghijklmnopqrstuvwxyz"
	chunks := SplitText(text, 10, 2)

	if len(chunks) != 3 {
		t.Fatalf("unexpected chunk count: %d", len(chunks))
	}

	t.Logf("chunks=%v", chunks)
}