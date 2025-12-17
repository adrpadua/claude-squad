package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectionOverlay represents a selection overlay for choosing from a list of options.
type SelectionOverlay struct {
	Title         string
	Options       []string
	SelectedIndex int
	Submitted     bool
	Canceled      bool
	OnSubmit      func(selected string)
	width, height int
}

// NewSelectionOverlay creates a new selection overlay with the given title and options.
func NewSelectionOverlay(title string, options []string) *SelectionOverlay {
	return &SelectionOverlay{
		Title:         title,
		Options:       options,
		SelectedIndex: 0,
		Submitted:     false,
		Canceled:      false,
	}
}

func (s *SelectionOverlay) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Init initializes the selection overlay model
func (s *SelectionOverlay) Init() tea.Cmd {
	return nil
}

// View renders the model's view
func (s *SelectionOverlay) View() string {
	return s.Render()
}

// HandleKeyPress processes a key press and updates the state accordingly.
// Returns true if the overlay should be closed.
func (s *SelectionOverlay) HandleKeyPress(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyUp:
		if s.SelectedIndex > 0 {
			s.SelectedIndex--
		}
		return false
	case tea.KeyDown:
		if s.SelectedIndex < len(s.Options)-1 {
			s.SelectedIndex++
		}
		return false
	case tea.KeyEnter:
		s.Submitted = true
		if s.OnSubmit != nil {
			s.OnSubmit(s.Options[s.SelectedIndex])
		}
		return true
	case tea.KeyEsc:
		s.Canceled = true
		return true
	default:
		// Handle j/k for vim-style navigation
		if msg.String() == "j" {
			if s.SelectedIndex < len(s.Options)-1 {
				s.SelectedIndex++
			}
			return false
		}
		if msg.String() == "k" {
			if s.SelectedIndex > 0 {
				s.SelectedIndex--
			}
			return false
		}
		return false
	}
}

// GetSelectedOption returns the currently selected option.
func (s *SelectionOverlay) GetSelectedOption() string {
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Options) {
		return s.Options[s.SelectedIndex]
	}
	return ""
}

// IsSubmitted returns whether the selection was submitted.
func (s *SelectionOverlay) IsSubmitted() bool {
	return s.Submitted
}

// IsCanceled returns whether the selection was canceled.
func (s *SelectionOverlay) IsCanceled() bool {
	return s.Canceled
}

// SetOnSubmit sets a callback function for selection submission.
func (s *SelectionOverlay) SetOnSubmit(onSubmit func(selected string)) {
	s.OnSubmit = onSubmit
}

// Render renders the selection overlay.
func (s *SelectionOverlay) Render() string {
	// Create styles
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true).
		MarginBottom(1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	selectedOptionStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MarginTop(1)

	// Build the view
	content := titleStyle.Render(s.Title) + "\n"

	for i, option := range s.Options {
		prefix := "  "
		optStyle := optionStyle
		if i == s.SelectedIndex {
			prefix = "> "
			optStyle = selectedOptionStyle
		}
		content += prefix + optStyle.Render(option) + "\n"
	}

	content += hintStyle.Render("↑↓/jk: navigate • enter: select • esc: cancel")

	return style.Render(content)
}
