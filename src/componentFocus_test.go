package src

import (
	"testing"

	//"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestComponentFocusSimulation(t *testing.T) {
	// Initialize the MainModel with Tab1Model
	m := MainModel{
		tab1: NewTab1Model(),
	}

	// Initial focus should be on the input
	assert.Equal(t, inputFocus, m.tab1.focus, "Initial focus should be on input")

	// Simulate pressing the key to focus table
	newModel, _ := m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	if newModel.focus != tableFocus {
		t.Error("focus is not equal to tableFocus")
		t.Error("expected: 3")
		t.Error("actual: ", newModel.focus)
	}

	// Simulate pressing the key to focus listOne
	newModel, _ = m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'#'}})
	if newModel.focus != listOneFocus {
		t.Error("focus is not equal to listOneFocus")
		t.Error("expected: 0")
		t.Error("actual: ", newModel.focus)
	}

	// Simulate pressing the key to focus listTwo
	newModel, _ = m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	if newModel.focus != listTwoFocus {
		t.Error("focus is not equal to listTwoFocus")
		t.Error("expected: 1")
		t.Error("actual: ", newModel.focus)
	}

}
