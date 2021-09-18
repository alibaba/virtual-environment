package envoy

import (
	"strconv"
	"testing"
)

func TestTagAppenderFilter(t *testing.T) {
	patchStruct, _ := buildPatchStruct("envHeader")
	if len(patchStruct.Fields) != 2 {
		t.Fatalf("patch fields count should not be " + strconv.Itoa(len(patchStruct.Fields)))
	}
	code := patchStruct.Fields["typed_config"].GetStructValue().Fields["inlineCode"].GetStringValue()
	//println(code)
	if len(code) != 855 {
		t.Fatalf("code len should not be " + strconv.Itoa(len(code)))
	}

}
