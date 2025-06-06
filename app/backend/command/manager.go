package command

import "Dr.uml/backend/utils/duerror"

// Command defines the basic behaviour for all commands.
type Command interface {
	Execute() duerror.DUError
	Unexecute() duerror.DUError
}

// Manager keeps track of executed commands and provides undo/redo
// capabilities.
type Manager struct {
	undoStack []Command
	redoStack []Command
}

// NewManager creates a new command manager instance.
func NewManager() *Manager {
	return &Manager{
		undoStack: make([]Command, 0),
		redoStack: make([]Command, 0),
	}
}

// ExecuteCommand runs the given command and stores it for undo.
func (m *Manager) ExecuteCommand(cmd Command) duerror.DUError {
	if cmd == nil {
		return duerror.NewInvalidArgumentError("command is nil")
	}
	if err := cmd.Execute(); err != nil {
		return err
	}
	m.undoStack = append(m.undoStack, cmd)
	m.redoStack = nil
	return nil
}

// Undo reverts the last executed command.
func (m *Manager) Undo() duerror.DUError {
	if len(m.undoStack) == 0 {
		return nil
	}
	cmd := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]
	if err := cmd.Unexecute(); err != nil {
		return err
	}
	m.redoStack = append(m.redoStack, cmd)
	return nil
}

// Redo re-executes the last undone command.
func (m *Manager) Redo() duerror.DUError {
	if len(m.redoStack) == 0 {
		return nil
	}
	cmd := m.redoStack[len(m.redoStack)-1]
	m.redoStack = m.redoStack[:len(m.redoStack)-1]
	if err := cmd.Execute(); err != nil {
		return err
	}
	m.undoStack = append(m.undoStack, cmd)
	return nil
}
