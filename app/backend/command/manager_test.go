package command

import (
	"testing"

	"Dr.uml/backend/utils/duerror"
	"github.com/stretchr/testify/assert"
)

type mockCommand struct {
	execCount int
	undoCount int
	execErr   duerror.DUError
	undoErr   duerror.DUError
}

func (m *mockCommand) Execute() duerror.DUError {
	m.execCount++
	return m.execErr
}

func (m *mockCommand) Unexecute() duerror.DUError {
	m.undoCount++
	return m.undoErr
}

func TestManager_ExecuteCommand(t *testing.T) {
	mgr := NewManager()
	cmd := &mockCommand{}

	err := mgr.ExecuteCommand(cmd)
	assert.NoError(t, err)
	assert.Equal(t, 1, cmd.execCount)
	assert.Len(t, mgr.undoStack, 1)
	assert.Len(t, mgr.redoStack, 0)
}

func TestManager_ExecuteCommand_Error(t *testing.T) {
	mgr := NewManager()
	cmd := &mockCommand{execErr: duerror.NewInvalidArgumentError("boom")}

	err := mgr.ExecuteCommand(cmd)
	assert.Error(t, err)
	assert.Len(t, mgr.undoStack, 0)
}

func TestManager_UndoRedo(t *testing.T) {
	mgr := NewManager()
	cmd := &mockCommand{}
	_ = mgr.ExecuteCommand(cmd)

	err := mgr.Undo()
	assert.NoError(t, err)
	assert.Equal(t, 1, cmd.undoCount)
	assert.Len(t, mgr.undoStack, 0)
	assert.Len(t, mgr.redoStack, 1)

	err = mgr.Redo()
	assert.NoError(t, err)
	assert.Equal(t, 2, cmd.execCount) // executed again
	assert.Len(t, mgr.undoStack, 1)
	assert.Len(t, mgr.redoStack, 0)
}

func TestManager_UndoEmptyRedoEmpty(t *testing.T) {
	mgr := NewManager()
	// Undo with empty stack
	err := mgr.Undo()
	assert.NoError(t, err)
	err = mgr.Redo()
	assert.NoError(t, err)
}

func TestManager_UndoError(t *testing.T) {
	mgr := NewManager()
	cmd := &mockCommand{undoErr: duerror.NewInvalidArgumentError("fail")}
	_ = mgr.ExecuteCommand(cmd)

	err := mgr.Undo()
	assert.Error(t, err)
	assert.Len(t, mgr.redoStack, 0)
}

func TestManager_RedoError(t *testing.T) {
	mgr := NewManager()
	cmd := &mockCommand{undoErr: nil}
	_ = mgr.ExecuteCommand(cmd)
	_ = mgr.Undo()

	cmd.execErr = duerror.NewInvalidArgumentError("redo fail")
	err := mgr.Redo()
	assert.Error(t, err)
	assert.Len(t, mgr.undoStack, 0)
}
