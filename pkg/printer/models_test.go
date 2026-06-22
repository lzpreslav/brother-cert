package printer

import (
	"reflect"
	"testing"
)

// TestModelFieldMaps pins each supported model's form fields to the exact
// values the pre-refactor hardcoded code emitted, so the model-map refactor is
// provably behavior-preserving and future edits can't silently drift a field.
func TestModelFieldMaps(t *testing.T) {
	want := map[string]model{
		"MFC-L2710DW": {
			name: "MFC-L2710DW",
			imp:  importFields{pageID: "390", empty1: "B8ea", empty2: "B8f8", file: "B820", passwd: "B821"},
			activate: activateFields{
				pageID:       "326",
				dropdown:     "B903",
				protocols:    map[string]string{"B86c": "1", "B87e": "1"},
				httpPageMode: "5",
			},
			del: deleteFields{pageID: "383", empty1: "B8ea", empty2: "B8fc"},
		},
	}

	for name, expect := range want {
		got, err := lookupModel(name)
		if err != nil {
			t.Errorf("lookupModel(%q) returned error: %v", name, err)
			continue
		}
		if !reflect.DeepEqual(got, expect) {
			t.Errorf("model %q field map mismatch:\n got  %+v\n want %+v", name, got, expect)
		}
	}
}

func TestLookupModelUnknown(t *testing.T) {
	if _, err := lookupModel("MFC-NOPE"); err == nil {
		t.Fatal("lookupModel(unknown) = nil error, want error")
	}
	if _, err := lookupModel(""); err == nil {
		t.Fatal("lookupModel(\"\") = nil error, want error")
	}
}
