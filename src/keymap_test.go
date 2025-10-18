package src

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func TestKeymapBindings(t *testing.T) {
	km := newKeyMap()

	testCases := []struct {
		name        string
		binding     key.Binding
		expectedKey string
		helpKey     string
		helpDesc    string
	}{
		{
			name:        "Input Focus Key",
			binding:     km.Input,
			expectedKey: "!",
			helpKey:     "!",
			helpDesc:    "focus input",
		},
		{
			name:        "List One Focus Key",
			binding:     km.List1,
			expectedKey: "#",
			helpKey:     "#",
			helpDesc:    "focus list one",
		},
		{
			name:        "List Two Focus Key",
			binding:     km.List2,
			expectedKey: "$",
			helpKey:     "$",
			helpDesc:    "focus list two",
		},
		{
			name:        "Table Focus Key",
			binding:     km.Table,
			expectedKey: "@",
			helpKey:     "@",
			helpDesc:    "focus table",
		},
		{
			name:        "InfoBox Focus Key",
			binding:     km.InfoBox,
			expectedKey: "%",
			helpKey:     "%",
			helpDesc:    "focus info box",
		},
		{
			name:        "Enter Key",
			binding:     km.Enter,
			expectedKey: "enter",
			helpKey:     "enter",
			helpDesc:    "Action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the key is properly bound
			assert.Contains(t, tc.binding.Keys(), tc.expectedKey)

			// Check help text
			help := tc.binding.Help()
			assert.Equal(t, tc.helpKey, help.Key)
			assert.Equal(t, tc.helpDesc, help.Desc)
		})
	}
}
