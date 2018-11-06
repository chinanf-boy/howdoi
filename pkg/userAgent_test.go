package howdoi

import (
	"testing"
)

func TestGetRandomUA(t *testing.T) {
	one := getRandomUA()
	if len(one) == 0 {
		t.Fail()
	}
}
