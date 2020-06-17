package envoy

import (
	"testing"
)

func TestTagAppenderFilter(t *testing.T) {
	patchStruct := buildPatchStruct("envLabel", "envHeader")
	if len(patchStruct.Fields) != 2 {
		t.Fail()
	}
	code := patchStruct.Fields["typed_config"].GetStructValue().Fields["inline_code"].GetStringValue()
	if len(code) != 692 {
		t.Fail()
	}

}
