package umldiagram

import (
	"slices"
	"time"

	"Dr.uml/backend/command"
	"Dr.uml/backend/component"
	"Dr.uml/backend/components"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
)

type DiagramType int

const (
	ClassDiagram = 1 << iota // 0x01
	UseCaseDiagram
	SequenceDiagram
	supportedType = ClassDiagram
)

var AllDiagramTypes = []struct {
	Value  DiagramType
	TSName string
}{
	{ClassDiagram, "ClassDiagram"},
}

// Other methods
func validateDiagramType(input DiagramType) duerror.DUError {
	if !(input&supportedType == input && input != 0) {
		return duerror.NewInvalidArgumentError("Invalid diagram type")
	}
	return nil
}

type UMLDiagram struct {
	name            string
	diagramType     DiagramType // e.g., "Class", "UseCase", "Sequence"
	lastModified    time.Time
	startPoint      utils.Point // for dragging and linking ass
	backgroundColor string

	componentsContainer components.Container
	componentsSelected  map[component.Component]bool
	associations        map[*component.Gadget]([2][]*component.Association)

	cmdMgr *command.Manager

	updateParentDraw func() duerror.DUError
	drawData         drawdata.Diagram
}

// Constructor
func CreateEmptyUMLDiagram(name string, dt DiagramType) (*UMLDiagram, duerror.DUError) {
	// TODO: also check the file is exist or not
	if err := utils.ValidateFilePath(name); err != nil {
		return nil, err
	}
	if err := validateDiagramType(dt); err != nil {
		return nil, err
	}
	return &UMLDiagram{
		name:                name,
		diagramType:         dt,
		lastModified:        time.Now(),
		startPoint:          utils.Point{X: 0, Y: 0},
		backgroundColor:     drawdata.DefaultDiagramColor, // Default white background
		componentsContainer: components.NewContainerMap(),
		associations:        make(map[*component.Gadget][2][]*component.Association),
		componentsSelected:  make(map[component.Component]bool),
		cmdMgr:              command.NewManager(),
		drawData: drawdata.Diagram{
			Margin:    drawdata.Margin,
			LineWidth: drawdata.LineWidth,
			Color:     drawdata.DefaultDiagramColor,
		},
	}, nil
}

func LoadExistUMLDiagram(name string) (*UMLDiagram, duerror.DUError) {
	// TODO
	return CreateEmptyUMLDiagram(name, ClassDiagram)
}

// Getters
func (ud *UMLDiagram) GetName() string {
	return ud.name
}

func (ud *UMLDiagram) GetDiagramType() DiagramType {
	return ud.diagramType
}

func (ud *UMLDiagram) GetLastModified() time.Time {
	return ud.lastModified
}

// Setters
func (ud *UMLDiagram) SetPointGadget(point utils.Point) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old point for undo
		oldPoint := g.GetPoint()

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetPoint(point) },
			undo: func() duerror.DUError { return g.SetPoint(oldPoint) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

func (ud *UMLDiagram) SetSetLayerGadget(layer int) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old layer for undo
		oldLayer := g.GetLayer()

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetLayer(layer) },
			undo: func() duerror.DUError { return g.SetLayer(oldLayer) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

func (ud *UMLDiagram) SetColorGadget(colorHexStr string) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old color for undo
		oldColor := g.GetColor()

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetColor(colorHexStr) },
			undo: func() duerror.DUError { return g.SetColor(oldColor) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

func (ud *UMLDiagram) SetAttrContentGadget(section int, index int, content string) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old content for undo
		var oldContent string
		if attrs := g.GetAttributes(section); attrs != nil && index >= 0 && index < len(attrs) {
			oldContent = attrs[index].GetContent()
		}

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetAttrContent(section, index, content) },
			undo: func() duerror.DUError { return g.SetAttrContent(section, index, oldContent) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}
func (ud *UMLDiagram) SetAttrSizeGadget(section int, index int, size int) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old size for undo
		var oldSize int
		if attrs := g.GetAttributes(section); attrs != nil && index >= 0 && index < len(attrs) {
			oldSize = attrs[index].GetSize()
		}

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetAttrSize(section, index, size) },
			undo: func() duerror.DUError { return g.SetAttrSize(section, index, oldSize) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}
func (ud *UMLDiagram) SetAttrStyleGadget(section int, index int, style int) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}

	switch g := c.(type) {
	case *component.Gadget:
		// Store the old style for undo
		var oldStyle int
		if attrs := g.GetAttributes(section); attrs != nil && index >= 0 && index < len(attrs) {
			oldStyle = int(attrs[index].GetStyle())
		}

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.SetAttrStyle(section, index, style) },
			undo: func() duerror.DUError { return g.SetAttrStyle(section, index, oldStyle) },
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

// Methods
func (ud *UMLDiagram) AddGadget(gadgetType component.GadgetType, point utils.Point, layer int, colorHexStr string, header string) duerror.DUError {
	cmd := &addGadgetCommand{
		diagram:    ud,
		gadgetType: gadgetType,
		point:      point,
		layer:      layer,
		color:      colorHexStr,
		header:     header,
	}
	return ud.cmdMgr.ExecuteCommand(cmd)
}

func (ud *UMLDiagram) StartAddAssociation(point utils.Point) duerror.DUError {
	if err := ud.validatePoint(point); err != nil {
		return err
	}
	ud.startPoint = point
	return nil
}

func (ud *UMLDiagram) EndAddAssociation(assType component.AssociationType, endPoint utils.Point) duerror.DUError {
	stPoint := ud.startPoint
	ud.startPoint = utils.Point{X: 0, Y: 0}
	if err := ud.validatePoint(endPoint); err != nil {
		return err
	}

	cmd := &addAssociationCommand{
		diagram: ud,
		assType: assType,
		start:   stPoint,
		end:     endPoint,
	}
	return ud.cmdMgr.ExecuteCommand(cmd)
}

func (ud *UMLDiagram) RemoveSelectedComponents() duerror.DUError {
	comps := make([]component.Component, 0, len(ud.componentsSelected))
	for c := range ud.componentsSelected {
		comps = append(comps, c)
	}
	if len(comps) == 0 {
		return nil
	}
	cmd := &removeComponentsCommand{
		diagram:    ud,
		components: comps,
	}
	ud.componentsSelected = make(map[component.Component]bool)
	return ud.cmdMgr.ExecuteCommand(cmd)
}

func (ud *UMLDiagram) SelectComponent(point utils.Point) duerror.DUError {
	c, err := ud.componentsContainer.Search(point)
	if err != nil {
		return err
	}
	if c == nil {
		return nil
	}
	// if is in componentsSelected remove it, else add it
	if _, ok := ud.componentsSelected[c]; ok {
		gadget := c.(*component.Gadget)
		gadget.SetIsSelected(false)
		delete(ud.componentsSelected, c)
	} else {
		gadget := c.(*component.Gadget)
		gadget.SetIsSelected(true)
		ud.componentsSelected[c] = true
	}
	//ud.componentsSelected[c] = true
	return ud.updateDrawData()
}

func (ud *UMLDiagram) UnselectComponent(point utils.Point) duerror.DUError {
	c, err := ud.componentsContainer.Search(point)
	if err != nil {
		return err
	}
	if c == nil {
		return nil
	}
	delete(ud.componentsSelected, c)
	return ud.updateDrawData()
}

func (ud *UMLDiagram) UnselectAllComponents() duerror.DUError {
	ud.componentsSelected = make(map[component.Component]bool)
	return nil
}

func (ud *UMLDiagram) AddAttributeToGadget(section int, content string) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}
	switch g := c.(type) {
	case *component.Gadget:
		// Get the current length to know where the new attribute will be added
		initialLengths := g.GetAttributesLen()
		addedIndex := -1
		if section >= 0 && section < len(initialLengths) {
			addedIndex = initialLengths[section]
		}

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.AddAttribute(section, content) },
			undo: func() duerror.DUError { 
				if addedIndex >= 0 {
					return g.RemoveAttribute(section, addedIndex)
				}
				return nil
			},
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

func (ud *UMLDiagram) RemoveAttributeFromGadget(section int, index int) duerror.DUError {
	c, err := ud.getSelectedComponent()
	if err != nil {
		return err
	}
	switch g := c.(type) {
	case *component.Gadget:
		// Store the content before removing so we can restore it on undo
		var content string
		if attrs := g.GetAttributes(section); attrs != nil && index >= 0 && index < len(attrs) {
			content = attrs[index].GetContent()
		}

		cmd := &funcCommand{
			exec: func() duerror.DUError { return g.RemoveAttribute(section, index) },
			undo: func() duerror.DUError {
				if content != "" {
					return g.AddAttribute(section, content)
				}
				return nil
			},
		}
		return ud.cmdMgr.ExecuteCommand(cmd)
	default:
		return duerror.NewInvalidArgumentError("selected component is not a gadget")
	}
}

// Private methods
func (ud *UMLDiagram) getSelectedComponent() (component.Component, duerror.DUError) {
	if len(ud.componentsSelected) != 1 {
		return nil, duerror.NewInvalidArgumentError("can only operate on one component")
	}
	for c := range ud.componentsSelected {
		return c, nil
	}
	return nil, duerror.NewInvalidArgumentError("no component selected")
}

func (ud *UMLDiagram) removeGadget(gad *component.Gadget) duerror.DUError {
	if _, ok := ud.associations[gad]; ok {
		for _, a := range ud.associations[gad][0] {
			if err := ud.removeAssociation(a); err != nil {
				return err
			}
		}
		for _, a := range ud.associations[gad][1] {
			if err := ud.removeAssociation(a); err != nil {
				return err
			}
		}
		delete(ud.associations, gad)
	}
	delete(ud.componentsSelected, gad)
	return ud.componentsContainer.Remove(gad)
}

func (ud *UMLDiagram) removeAssociation(a *component.Association) duerror.DUError {
	st := a.GetParentStart()
	en := a.GetParentEnd()
	if _, ok := ud.associations[st]; ok {
		stList := ud.associations[st][0]
		index := slices.Index(stList, a)
		if index >= 0 {
			stList = slices.Delete(stList, index, index+1)
		}
		ud.associations[st] = [2][]*component.Association{stList, ud.associations[st][1]}
	}
	if _, ok := ud.associations[en]; ok {
		enList := ud.associations[en][1]
		index := slices.Index(enList, a)
		if index >= 0 {
			enList = slices.Delete(enList, index, index+1)
		}
		ud.associations[en] = [2][]*component.Association{ud.associations[en][0], enList}
	}
	delete(ud.componentsSelected, a)
	return ud.componentsContainer.Remove(a)
}

func (ud *UMLDiagram) validatePoint(point utils.Point) duerror.DUError {
	if point.X < 0 || point.Y < 0 {
		return duerror.NewInvalidArgumentError("point coordinates must be non-negative")
	}
	return nil
}

// draw
func (ud *UMLDiagram) GetDrawData() drawdata.Diagram {
	return ud.drawData
}

func (ud *UMLDiagram) RegisterUpdateParentDraw(update func() duerror.DUError) duerror.DUError {
	if update == nil {
		return duerror.NewInvalidArgumentError("update function cannot be nil")
	}
	ud.updateParentDraw = update
	return nil

}

func (ud *UMLDiagram) updateDrawData() duerror.DUError {
	gs := make([]drawdata.Gadget, 0, len(ud.componentsSelected))
	as := make([]drawdata.Association, 0, len(ud.componentsSelected))
	for _, c := range ud.componentsContainer.GetAll() {
		cDrawData := c.GetDrawData()
		if cDrawData == nil {
			continue
		}
		switch c.(type) {
		case *component.Gadget:
			gs = append(gs, cDrawData.(drawdata.Gadget))
		case *component.Association:
			as = append(as, cDrawData.(drawdata.Association))
		}
	}
	ud.drawData.Gadgets = gs
	ud.drawData.Associations = as
	if ud.updateParentDraw == nil {
		return nil
	}
	return ud.updateParentDraw()
}
