package openai

import "testing"

func TestEvent(t *testing.T) {
	var (
		data = "foo"
	)

	e := SSEEvent{Data: data}

	if got, want := e.Data, "foo"; got != want {
		t.Fatalf(`e.Data = %s, want %s`, got, want)
	}
}
