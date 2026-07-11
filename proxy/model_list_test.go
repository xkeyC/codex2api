package proxy

import "testing"

func TestNormalizeModelListJSON(t *testing.T) {
	got, err := NormalizeModelListJSON(`[" gpt5.4 ","gpt-5.4","gpt-5.5"]`, []string{"gpt-5.4", "gpt-5.5"}, "disabled_image_generation_models")
	if err != nil {
		t.Fatalf("NormalizeModelListJSON returned error: %v", err)
	}
	if got != `["gpt-5.4","gpt-5.5"]` {
		t.Fatalf("NormalizeModelListJSON = %s, want canonical unique list", got)
	}
}

func TestNormalizeModelListJSONRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"model":"gpt-5.4"}`,
		`["gpt-5.4",""]`,
		`["gpt-5.4!"]`,
		`["gpt-5.4"] []`,
	} {
		if _, err := NormalizeModelListJSON(input, nil, "disabled_image_generation_models"); err == nil {
			t.Fatalf("NormalizeModelListJSON(%s) expected error", input)
		}
	}
}

func TestModelListJSONContains(t *testing.T) {
	if !modelListJSONContains(`["gpt-5.4","gpt-5.5"]`, "GPT-5.4") {
		t.Fatal("modelListJSONContains should match case-insensitively")
	}
	if modelListJSONContains(`["gpt-5.4"]`, "gpt-5.5") {
		t.Fatal("modelListJSONContains should not match absent model")
	}
	if modelListJSONContains(`not-json`, "gpt-5.4") {
		t.Fatal("modelListJSONContains should ignore malformed JSON")
	}
}
