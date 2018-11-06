package howdoi

import "testing"

func TestGetEnv(t *testing.T) {
	res := getEnv("NOTHINGNOTHINGNOTHING", "ok")

	if res != "ok" {
		t.Fail()
	}
}
