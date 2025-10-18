package src

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppStateSimulation(t *testing.T) {
	initialModel := NewMainModel()
	initialModel.currentScreen = AppScreen

	testCases := []struct {
		name          string
		width         int
		height        int
		expectedState AppState
	}{
		{
			name:          "Valid dimensions",
			width:         120,
			height:        40,
			expectedState: AppScreen,
		},
		{
			name:          "Invalid width",
			width:         90,
			height:        40,
			expectedState: ErrorScreen,
		},
		{
			name:          "Invalid height",
			width:         120,
			height:        20,
			expectedState: ErrorScreen,
		},
		{
			name:          "Both dimensions invalid",
			width:         90,
			height:        20,
			expectedState: ErrorScreen,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tea.WindowSizeMsg{
				Width:  tc.width,
				Height: tc.height,
			}

			updatedModel, _ := initialModel.Update(msg)
			model := updatedModel.(MainModel)

			if model.currentScreen != tc.expectedState {
				t.Errorf("Test case '%s' failed:\nConditions: Width=%d, Height=%d\nExpected state: %v\nGot: %v",
					tc.name, tc.width, tc.height, tc.expectedState, model.currentScreen)
			}
		})
	}
}
