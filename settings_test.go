package main

import (
	"encoding/json"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestValidateSettings(t *testing.T) {
	cases := []struct {
		name                 string
		requiredAnnotations  map[string]string
		forbiddenAnnotations mapset.Set[string]
		isValid              bool
	}{
		{
			"empty",
			map[string]string{},
			mapset.NewSet[string](),
			true,
		},
		{
			"only required annotations",
			map[string]string{
				"cc-center": "marketing",
			},
			mapset.NewSet[string](),
			true,
		},
		{
			"only forbidden annotations",
			map[string]string{},
			mapset.NewSet[string]("priority"),
			true,
		},
		{
			"no contradictions",
			map[string]string{
				"cc-center": "marketing",
			},
			mapset.NewSet[string]("priority"),
			true,
		},
		{
			"contradictions",
			map[string]string{
				"cc-center": "marketing",
			},
			mapset.NewSet[string]("cc-center"),
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			settings := Settings{
				RequiredAnnotations:  tc.requiredAnnotations,
				ForbiddenAnnotations: tc.forbiddenAnnotations,
			}
			settingsJSON, err := json.Marshal(&settings)
			if err != nil {
				t.Errorf("cannot marshal settings: %v", err)
			}

			responseJSON := validateSettings(settingsJSON)
			var response SettingsValidationResponse
			err = json.Unmarshal(responseJSON, &response)
			if err != nil {
				t.Errorf("cannot unmarshal response: %v", err)
			}

			if response.Valid != tc.isValid {
				t.Errorf(
					"didn't get the expected validation outcome, %v was expected, got %v instead",
					tc.isValid, response.Valid)
				if response.Message != nil {
					t.Errorf(
						"validation message: %s",
						*response.Message)
				}
			}
		})
	}
}
