package component

import (
	"Dr.uml/backend/component/attribute"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
)

type GadgetType int

const (
	Class               GadgetType = 1 << iota // 0x01
	supportedGadgetType            = Class
)

var AllGadgetTypes = []struct {
	Value  GadgetType
	TSName string
}{
	{Class, "Class"},
}

type Gadget struct {
	gadgetType       GadgetType
	point            utils.Point
	layer            int
	attributes       [][]*attribute.Attribute // Gadget has multiple sections, each section has multiple attributes
	color            string
	IsSelected       bool
	drawData         drawdata.Gadget
	updateParentDraw func() duerror.DUError
}

// Other functions
func validateGadgetType(input GadgetType) duerror.DUError {
	if !(input&supportedGadgetType == input && input != 0) {
		return duerror.NewInvalidArgumentError("gadget type is not supported")
	}
	return nil
}

// Constructor
func NewGadget(gadgetType GadgetType, point utils.Point, layer int, colorHexStr string, header string) (*Gadget, duerror.DUError) {
	if err := validateGadgetType(gadgetType); err != nil {
		return nil, err
	}
	g := Gadget{
		gadgetType: gadgetType,
		point:      point,
		layer:      layer,
		color:      colorHexStr,
	}

	// Init attributes with three sections
	g.attributes = make([][]*attribute.Attribute, 3)

	// The first section contains the header
	g.attributes[0] = make([]*attribute.Attribute, 0, 1)
	if header != "" {
		if err := g.AddAttribute(0, header); err != nil {
			return nil, err
		}
	}

	// The second and third sections are empty
	g.attributes[1] = make([]*attribute.Attribute, 0)
	g.attributes[2] = make([]*attribute.Attribute, 0)

	if err := g.updateDrawData(); err != nil {
		return nil, err
	}
	return &g, nil
}

// Getter
func (g *Gadget) GetPoint() utils.Point {
	return g.point
}

func (g *Gadget) GetLayer() int {
	return g.layer
}

func (g *Gadget) GetColor() string {
	return g.color
}

func (g *Gadget) GetGadgetType() GadgetType {
	return g.gadgetType
}

func (g *Gadget) GetAttributesLen() []int {
	lengths := make([]int, len(g.attributes))
	for i, atts := range g.attributes {
		lengths[i] = len(atts)
	}
	return lengths
}
func (g *Gadget) GetIsSelected() bool {
	return g.IsSelected
}

func (g *Gadget) GetAttributes(section int) []*attribute.Attribute {
	if err := g.validateSection(section); err != nil {
		return nil
	}
	return g.attributes[section]
}

// Setter
func (g *Gadget) SetPoint(point utils.Point) duerror.DUError {
	g.point = point
	g.drawData.X = point.X
	g.drawData.Y = point.Y
	return g.updateParentDraw()
}

func (g *Gadget) SetLayer(layer int) duerror.DUError {
	g.layer = layer
	g.drawData.Layer = layer
	return g.updateParentDraw()
}

func (g *Gadget) SetColor(colorHexStr string) duerror.DUError {
	g.color = colorHexStr
	g.drawData.Color = colorHexStr
	return g.updateParentDraw()
}

func (g *Gadget) SetAttrContent(section int, index int, content string) duerror.DUError {
	if err := g.validateSection(section); err != nil {
		return err
	}
	if err := g.validateIndex(index, section); err != nil {
		return err
	}
	if err := g.attributes[section][index].SetContent(content); err != nil {
		return err
	}
	return g.updateDrawData()
}

func (g *Gadget) SetAttrSize(section int, index int, size int) duerror.DUError {
	if err := g.validateSection(section); err != nil {
		return err
	}
	if err := g.validateIndex(index, section); err != nil {
		return err
	}
	if err := g.attributes[section][index].SetSize(size); err != nil {
		return err
	}
	return g.updateDrawData()
}

func (g *Gadget) SetAttrStyle(section int, index int, style int) duerror.DUError {
	if err := g.validateSection(section); err != nil {
		return err
	}
	if err := g.validateIndex(index, section); err != nil {
		return err
	}
	if err := g.attributes[section][index].SetStyle(attribute.Textstyle(style)); err != nil {
		return err
	}
	return g.updateDrawData()
}
func (g *Gadget) SetIsSelected(isSelected bool) duerror.DUError {
	g.IsSelected = isSelected
	g.drawData.IsSelected = isSelected
	return g.updateDrawData()
}

// Methods
func (g *Gadget) Cover(p utils.Point) (bool, duerror.DUError) {
	tl := g.point                                                                          // top-left
	br := utils.AddPoints(g.point, utils.Point{X: g.drawData.Width, Y: g.drawData.Height}) // bottom-right
	return p.X >= tl.X && p.X <= br.X && p.Y >= tl.Y && p.Y <= br.Y, nil
}

func (g *Gadget) AddAttribute(section int, content string) duerror.DUError {
	if err := g.validateSection(section); err != nil {
		return err
	}
	att, err := attribute.NewAttribute(content)
	if err != nil {
		return err
	}
	if err = att.RegisterUpdateParentDraw(g.updateDrawData); err != nil {
		return err
	}
	g.attributes[section] = append(g.attributes[section], att)
	return g.updateDrawData()
}

func (g *Gadget) RemoveAttribute(section int, index int) duerror.DUError {
	if err := g.validateSection(section); err != nil {
		return err
	}
	if err := g.validateIndex(index, section); err != nil {
		return err
	}
	g.attributes[section] = append(g.attributes[section][:index], g.attributes[section][index+1:]...)
	return g.updateDrawData()
}
func (g *Gadget) validateSection(section int) duerror.DUError {
	if section < 0 || section >= len(g.attributes) {
		return duerror.NewInvalidArgumentError("section out of range")
	}
	return nil
}

func (g *Gadget) validateIndex(index, section int) duerror.DUError {
	if index < 0 || index >= len(g.attributes[section]) {
		return duerror.NewInvalidArgumentError("index out of range")
	}
	return nil
}

// Draw
func (g *Gadget) GetDrawData() any {
	return g.drawData
}

func (g *Gadget) updateDrawData() duerror.DUError {
	height := drawdata.LineWidth
	maxAttWidth := 0
	atts := make([][]drawdata.Attribute, len(g.attributes))
	for i, attsRow := range g.attributes {
		atts[i] = make([]drawdata.Attribute, 0, len(attsRow))
		for _, att := range attsRow {
			attDrawData := att.GetDrawData()
			atts[i] = append(atts[i], attDrawData)
			if attDrawData.Width > maxAttWidth {
				maxAttWidth = attDrawData.Width
			}
			height += drawdata.Margin + attDrawData.Height
		}
		height += drawdata.Margin + drawdata.LineWidth
	}
	width := maxAttWidth + drawdata.Margin*2 + drawdata.LineWidth*2

	g.drawData.GadgetType = int(g.gadgetType)
	g.drawData.X = g.point.X
	g.drawData.Y = g.point.Y
	g.drawData.Layer = g.layer
	g.drawData.Height = height
	g.drawData.Width = width
	g.drawData.Color = g.color
	g.drawData.Attributes = atts

	if g.updateParentDraw == nil {
		return nil
	}
	return g.updateParentDraw()
}

func (g *Gadget) RegisterUpdateParentDraw(update func() duerror.DUError) duerror.DUError {
	if update == nil {
		return duerror.NewInvalidArgumentError("update function is nil")
	}
	g.updateParentDraw = update
	return nil
}
