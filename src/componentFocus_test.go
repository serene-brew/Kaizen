package src

import (
	"testing"

	//"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestComponentFocusSimulation(t *testing.T) {
	m := MainModel{
		tab1: NewTab1Model(),
	}

	assert.Equal(t, inputFocus, m.tab1.focus, "Initial focus should be on input")

	model, _ := m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	if tab1Model, ok := model.(Tab1Model); ok {
		assert.Equal(t, tableFocus, tab1Model.focus, "Focus should be on table")
	} else {
		t.Error("Could not convert model to Tab1Model")
	}

	model, _ = m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'#'}})
	if tab1Model, ok := model.(Tab1Model); ok {
		assert.Equal(t, listOneFocus, tab1Model.focus, "Focus should be on listOne")
	} else {
		t.Error("Could not convert model to Tab1Model")
	}

	model, _ = m.tab1.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	if tab1Model, ok := model.(Tab1Model); ok {
		assert.Equal(t, listTwoFocus, tab1Model.focus, "Focus should be on listTwo")
	} else {
		t.Error("Could not convert model to Tab1Model")
	}

}
