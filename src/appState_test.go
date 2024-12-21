package src

import(
	"testing"
	tea "github.com/charmbracelet/bubbletea"
)

func TestAppStateSimulation(t *testing.T) {
    // Initialize your model
    initialModel := MainModel{
        currentScreen: AppScreen, // Initial state
        width:         0,
        height:        0,
    }
	//---------------------------------------------------------------------------------------------------
   
	windowSizeMsg := tea.WindowSizeMsg{
        Width:  120,
        Height: 40,
    }
    // Update the model with the window size message
    updatedModel, _ := initialModel.Update(windowSizeMsg)
    // Assert the currentScreen state if applicable
    if updatedModel.(MainModel).currentScreen != AppScreen {
		t.Error("failed app state test")
		t.Error("conditions: Width=120, Height=40")
        t.Errorf("expected currentScreen to be AppScreen, got %v", updatedModel.(MainModel).currentScreen)
    }

	//---------------------------------------------------------------------------------------------------

	windowSizeMsg = tea.WindowSizeMsg{
        Width:  90,
        Height: 40,
    }
    // Update the model with the window size message
    updatedModel, _ = initialModel.Update(windowSizeMsg)
    // Assert the currentScreen state if applicable
    if updatedModel.(MainModel).currentScreen != ErrorScreen {
    	t.Error("failed app state test")
		t.Error("conditions: Width=90, Height=40")
        t.Errorf("expected currentScreen to be AppScreen, got %v", updatedModel.(MainModel).currentScreen)
	}

	//---------------------------------------------------------------------------------------------------

	windowSizeMsg = tea.WindowSizeMsg{
        Width:  120,
        Height: 20,
    }
    // Update the model with the window size message
    updatedModel, _ = initialModel.Update(windowSizeMsg)
    // Assert the currentScreen state if applicable
    if updatedModel.(MainModel).currentScreen != ErrorScreen {
    	t.Error("failed app state test")
		t.Error("conditions: Width=120, Height=20")
        t.Errorf("expected currentScreen to be AppScreen, got %v", updatedModel.(MainModel).currentScreen)
	}

	//---------------------------------------------------------------------------------------------------
	
	windowSizeMsg = tea.WindowSizeMsg{
        Width:  90,
        Height: 20,
    }
    // Update the model with the window size message
    updatedModel, _ = initialModel.Update(windowSizeMsg)
    // Assert the currentScreen state if applicable
    if updatedModel.(MainModel).currentScreen != ErrorScreen {
    	t.Error("failed app state test")
		t.Error("conditions: Width=90, Height=20")
        t.Errorf("expected currentScreen to be AppScreen, got %v", updatedModel.(MainModel).currentScreen)
	}
}

