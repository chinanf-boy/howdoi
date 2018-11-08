package howdoi

import "testing"

func TestGetEnv(t *testing.T) {
	res := getEnv("NOTHINGNOTHINGNOTHING", "ok")

	if res != "ok" {
		t.Fail()
	}
}

func TestUqineSlice(t *testing.T) {
	old := []string{"h", "o", "o", "l"}
	res := UqineSlice(old)
	if len(res) != 3 {
		t.Fail()
	}
}

func TestColorString(t *testing.T) {

	r := red("rea string")

	g := gree("green string")
	c := cyan("cyan string")
	f := format("format %s", "str")

	if len(r) == 0 || len(g) == 0 || len(c) == 0 || len(f) == 0 {
		t.Fail()
	}
}
