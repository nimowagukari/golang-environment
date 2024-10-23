package main

import (
	"strings"
	"testing"
)

func TestYamabiko(t *testing.T) {
	stdin := "Hello,World.\n"

	want := "Hello,World.\nHello,World.\n"
	got := Yamabiko(strings.NewReader(stdin))
	if got != want {
		t.Errorf("got: %v, but want: %v\n", got, want)
	}
}
