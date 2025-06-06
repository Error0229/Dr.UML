package umldiagram

import (
	"testing"

	"Dr.uml/backend/component"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
	"github.com/stretchr/testify/assert"
)

func TestFuncCommand(t *testing.T) {
	called := false
	undone := false
	cmd := &funcCommand{
		exec: func() duerror.DUError { called = true; return nil },
		undo: func() duerror.DUError { undone = true; return nil },
	}
	assert.NoError(t, cmd.Execute())
	assert.True(t, called)
	assert.NoError(t, cmd.Unexecute())
	assert.True(t, undone)
}

func TestAddGadgetCommand_ExecuteUnexecute(t *testing.T) {
	d, err := CreateEmptyUMLDiagram("cmd_test.uml", ClassDiagram)
	assert.NoError(t, err)
	prev := d.lastModified

	cmd := &addGadgetCommand{
		diagram:    d,
		gadgetType: component.Class,
		point:      utils.Point{X: 10, Y: 10},
		layer:      0,
		color:      drawdata.DefaultGadgetColor,
		header:     "Header",
	}
	err = cmd.Execute()
	assert.NoError(t, err)
	l, _ := d.componentsContainer.Len()
	assert.Equal(t, 1, l)
	assert.True(t, d.lastModified.After(prev) || d.lastModified.Equal(prev))
	assert.NotNil(t, cmd.gadget)
	assert.Len(t, d.associations, 1)

	err = cmd.Unexecute()
	assert.NoError(t, err)
	l, _ = d.componentsContainer.Len()
	assert.Equal(t, 0, l)
}

func TestAddAssociationCommand_ExecuteUnexecute(t *testing.T) {
	d, err := CreateEmptyUMLDiagram("ass_cmd_test.uml", ClassDiagram)
	assert.NoError(t, err)

	// add two gadgets
	_ = d.AddGadget(component.Class, utils.Point{X: 5, Y: 5}, 0, drawdata.DefaultGadgetColor, "")
	_ = d.AddGadget(component.Class, utils.Point{X: 100, Y: 100}, 0, drawdata.DefaultGadgetColor, "")
	gs := d.componentsContainer.GetAll()
	g1 := gs[0].(*component.Gadget)
	g2 := gs[1].(*component.Gadget)
	g1dd := g1.GetDrawData().(drawdata.Gadget)
	g2dd := g2.GetDrawData().(drawdata.Gadget)

	cmd := &addAssociationCommand{
		diagram: d,
		assType: component.Extension,
		start:   utils.Point{X: g1dd.X + 1, Y: g1dd.Y + 1},
		end:     utils.Point{X: g2dd.X + 1, Y: g2dd.Y + 1},
	}
	prev := d.lastModified
	err = cmd.Execute()
	assert.NoError(t, err)
	l, _ := d.componentsContainer.Len()
	assert.Equal(t, 3, l) // 2 gadgets + 1 association
	assert.True(t, d.lastModified.After(prev) || d.lastModified.Equal(prev))
	assert.NotNil(t, cmd.association)

	err = cmd.Unexecute()
	assert.NoError(t, err)
	l, _ = d.componentsContainer.Len()
	assert.Equal(t, 2, l)
}

func TestRemoveComponentsCommand(t *testing.T) {
	d, err := CreateEmptyUMLDiagram("remove_cmd.uml", ClassDiagram)
	assert.NoError(t, err)
	_ = d.AddGadget(component.Class, utils.Point{X: 1, Y: 1}, 0, drawdata.DefaultGadgetColor, "")
	g := d.componentsContainer.GetAll()[0]
	cmd := &removeComponentsCommand{
		diagram:    d,
		components: []component.Component{g},
	}
	prev := d.lastModified
	err = cmd.Execute()
	assert.NoError(t, err)
	l, _ := d.componentsContainer.Len()
	assert.Equal(t, 0, l)
	assert.True(t, d.lastModified.After(prev) || d.lastModified.Equal(prev))
	assert.NoError(t, cmd.Unexecute())
}
