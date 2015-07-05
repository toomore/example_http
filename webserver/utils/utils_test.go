package utils

import (
	"net/url"
	"testing"
)

func TestValues2Map(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Toomore")
	v.Set("age", "30")
	t.Logf("v: %+v", v)
	t.Log("Values to Map:", Values2Map(v))
}
