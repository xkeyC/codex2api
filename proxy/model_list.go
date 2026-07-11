package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/codex2api/security"
)

func NormalizeModelListJSON(value string, supportedModels []string, fieldName string) (string, error) {
	models, err := parseModelListJSON(value, supportedModels, fieldName, true)
	if err != nil {
		return "", err
	}
	if len(models) == 0 {
		return "[]", nil
	}
	body, err := json.Marshal(models)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func modelListJSONContains(value string, model string) bool {
	model = strings.TrimSpace(model)
	if model == "" {
		return false
	}
	models, _ := parseModelListJSON(value, nil, "model_list", false)
	for _, candidate := range models {
		if strings.EqualFold(candidate, model) {
			return true
		}
	}
	return false
}

func parseModelListJSON(value string, supportedModels []string, fieldName string, strict bool) ([]string, error) {
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		fieldName = "model_list"
	}
	value = strings.TrimSpace(value)
	if value == "" || value == "[]" {
		return nil, nil
	}

	var raw []string
	dec := json.NewDecoder(strings.NewReader(value))
	if err := dec.Decode(&raw); err != nil {
		if strict {
			return nil, fmt.Errorf("%s 必须是 JSON 字符串数组: %w", fieldName, err)
		}
		return nil, nil
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if strict {
			if err == nil {
				return nil, fmt.Errorf("%s 只能包含一个 JSON 数组", fieldName)
			}
			return nil, fmt.Errorf("%s JSON 无效: %w", fieldName, err)
		}
		return nil, nil
	}

	models := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for i, rawModel := range raw {
		model := normalizeConfiguredModelName(rawModel, supportedModels)
		if model == "" {
			if strict {
				return nil, fmt.Errorf("%s[%d] 需要非空且合法的 model", fieldName, i)
			}
			continue
		}
		key := strings.ToLower(model)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		models = append(models, model)
	}
	return models, nil
}

func normalizeConfiguredModelName(model string, supportedModels []string) string {
	model = strings.TrimSpace(model)
	if model == "" {
		return ""
	}
	canonical := canonicalizeCodexModel(model, supportedModels)
	if canonical == "" {
		canonical = model
	}
	if err := security.ValidateModelName(canonical); err != nil {
		return ""
	}
	return canonical
}
