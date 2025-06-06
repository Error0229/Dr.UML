package umldiagram

import (
	"time"

	"Dr.uml/backend/command"
	"Dr.uml/backend/component"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
)

// addGadgetCommand adds a gadget to the diagram.
type addGadgetCommand struct {
	diagram    *UMLDiagram
	gadgetType component.GadgetType
	point      utils.Point
	layer      int
	color      string
	header     string
	gadget     *component.Gadget
}

func (c *addGadgetCommand) Execute() duerror.DUError {
	g, err := component.NewGadget(c.gadgetType, c.point, c.layer, c.color, c.header)
	if err != nil {
		return err
	}
	if err = g.RegisterUpdateParentDraw(c.diagram.updateDrawData); err != nil {
		return err
	}
	if err = c.diagram.componentsContainer.Insert(g); err != nil {
		return err
	}
	c.diagram.associations[g] = [2][]*component.Association{{}, {}}
	c.diagram.lastModified = time.Now()
	c.gadget = g
	return c.diagram.updateDrawData()
}

func (c *addGadgetCommand) Unexecute() duerror.DUError {
	if c.gadget == nil {
		return nil
	}
	return c.diagram.removeGadget(c.gadget)
}

// addAssociationCommand creates an association between two gadgets.
type addAssociationCommand struct {
	diagram     *UMLDiagram
	assType     component.AssociationType
	start, end  utils.Point
	association *component.Association
}

func (c *addAssociationCommand) Execute() duerror.DUError {
	stGad, err := c.diagram.componentsContainer.SearchGadget(c.start)
	if err != nil {
		return err
	}
	if stGad == nil {
		return duerror.NewInvalidArgumentError("start point does not contain a gadget")
	}
	enGad, err := c.diagram.componentsContainer.SearchGadget(c.end)
	if err != nil {
		return err
	}
	if enGad == nil {
		return duerror.NewInvalidArgumentError("end point does not contain a gadget")
	}
	parents := [2]*component.Gadget{stGad, enGad}
	a, err := component.NewAssociation(parents, component.AssociationType(c.assType), c.start, c.end)
	if err != nil {
		return err
	}
	if err = a.RegisterUpdateParentDraw(c.diagram.updateDrawData); err != nil {
		return err
	}
	if err = c.diagram.componentsContainer.Insert(a); err != nil {
		return err
	}
	tmp := c.diagram.associations[stGad]
	tmp[0] = append(tmp[0], a)
	c.diagram.associations[stGad] = tmp
	tmp = c.diagram.associations[enGad]
	tmp[1] = append(tmp[1], a)
	c.diagram.associations[enGad] = tmp
	c.diagram.lastModified = time.Now()
	c.association = a
	return c.diagram.updateDrawData()
}

func (c *addAssociationCommand) Unexecute() duerror.DUError {
	if c.association == nil {
		return nil
	}
	return c.diagram.removeAssociation(c.association)
}

// removeComponentsCommand removes a list of components from the diagram.
type removeComponentsCommand struct {
	diagram    *UMLDiagram
	components []component.Component
}

func (c *removeComponentsCommand) Execute() duerror.DUError {
	for _, comp := range c.components {
		switch comp := comp.(type) {
		case *component.Gadget:
			if err := c.diagram.removeGadget(comp); err != nil {
				return err
			}
		case *component.Association:
			if err := c.diagram.removeAssociation(comp); err != nil {
				return err
			}
		}
	}
	c.diagram.lastModified = time.Now()
	return c.diagram.updateDrawData()
}

func (c *removeComponentsCommand) Unexecute() duerror.DUError {
	// Add back each removed component
	for _, comp := range c.components {
		switch comp := comp.(type) {
		case *component.Gadget:
			// Register update parent draw and insert back into container
			if err := comp.RegisterUpdateParentDraw(c.diagram.updateDrawData); err != nil {
				return err
			}
			if err := c.diagram.componentsContainer.Insert(comp); err != nil {
				return err
			}
			// Recreate empty association arrays for the gadget
			c.diagram.associations[comp] = [2][]*component.Association{{}, {}}

		case *component.Association:
			// Get the parent gadgets and recreate association map entries
			stGad := comp.GetParentStart()
			enGad := comp.GetParentEnd()

			if err := comp.RegisterUpdateParentDraw(c.diagram.updateDrawData); err != nil {
				return err
			}
			if err := c.diagram.componentsContainer.Insert(comp); err != nil {
				return err
			}

			// Add association back to each gadget's association list
			tmp := c.diagram.associations[stGad]
			tmp[0] = append(tmp[0], comp)
			c.diagram.associations[stGad] = tmp

			tmp = c.diagram.associations[enGad]
			tmp[1] = append(tmp[1], comp)
			c.diagram.associations[enGad] = tmp
		}
	}

	c.diagram.lastModified = time.Now()
	return c.diagram.updateDrawData()
}

// helper to wrap a function as a command.
type funcCommand struct {
	exec func() duerror.DUError
	undo func() duerror.DUError
}

func (f *funcCommand) Execute() duerror.DUError {
	if f.exec != nil {
		return f.exec()
	}
	return nil
}

func (f *funcCommand) Unexecute() duerror.DUError {
	if f.undo != nil {
		return f.undo()
	}
	return nil
}

var _ command.Command = (*addGadgetCommand)(nil)
var _ command.Command = (*addAssociationCommand)(nil)
var _ command.Command = (*removeComponentsCommand)(nil)
var _ command.Command = (*funcCommand)(nil)
