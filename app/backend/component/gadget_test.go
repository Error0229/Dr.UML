package component

import (
	"fmt"
	"testing"

	"Dr.uml/backend/component/attribute"
	"Dr.uml/backend/drawdata"
	"Dr.uml/backend/utils"
	"Dr.uml/backend/utils/duerror"
	"github.com/stretchr/testify/assert"
)

// testing purpose
var gadgetDefaultAtts = map[GadgetType][][]string{
	Class: {
		{"UMLProject"},
		{"id: String", "name: String", "lastModified: Date"},
		{"GetAvailableDiagrams(): List<String>", "GetLastOpenedDiagrams(): List<String>", "SelectDiagram(diagramName: String): DUError", "CreateDiagram(diagramName: String): DUError"},
	},
}

// test util
func newEmptyGadget(gadgetType GadgetType, point utils.Point) *Gadget {
	g := &Gadget{
		gadgetType: gadgetType,
		point:      point,
		layer:      0,
		color:      drawdata.DefaultGadgetColor,
	}

	// Initialize attributes using gadgetDefaultAtts
	g.attributes = make([][]*attribute.Attribute, len(gadgetDefaultAtts[gadgetType]))
	for i, contents := range gadgetDefaultAtts[gadgetType] {
		g.attributes[i] = make([]*attribute.Attribute, 0, len(contents))
		for _, content := range contents {
			att, err := attribute.NewAttribute(content)
			if err != nil {
				panic(err)
			}
			g.attributes[i] = append(g.attributes[i], att)
		}
	}

	// Initialize drawData
	if err := g.updateDrawData(); err != nil {
		panic(err)
	}

	return g
}

type mockParent struct {
	Times int
}

func (m *mockParent) UpdateParentDraw() duerror.DUError {
	m.Times++
	return nil
}

// Constructor
func TestNewGadget(t *testing.T) {
	// success
	g, err := NewGadget(Class, utils.Point{X: 1, Y: 1}, 0, drawdata.DefaultGadgetColor, "")
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, Class, g.GetGadgetType())

	// invalid gadget type
	g, err = NewGadget(-1, utils.Point{X: 1, Y: 1}, 0, drawdata.DefaultGadgetColor, "")
	assert.Error(t, err)

	// some errors are hard to test :(
}

// Getter
func TestGetPoint(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	assert.Equal(t, utils.Point{X: 1, Y: 1}, g.GetPoint())
}

func TestGetLayer(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	assert.Equal(t, 0, g.GetLayer())
}

func TestGetColor(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	assert.Equal(t, drawdata.DefaultGadgetColor, g.GetColor())
}

func TestGetGadgetType(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	assert.Equal(t, Class, g.GetGadgetType())
}

func TestGetAttributesLen(t *testing.T) {
	// for Class type
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	assert.Equal(t, []int{1, 3, 4}, g.GetAttributesLen())
}

func TestGetAttributes(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})

	// Test valid section
	attrs := g.GetAttributes(1) // Section 1 has attributes by default
	assert.NotNil(t, attrs)
	assert.Equal(t, 3, len(attrs)) // Should have 3 attributes by default in section 1

	// Add a new attribute and verify
	assert.NoError(t, g.AddAttribute(1, "test attr"))
	attrs = g.GetAttributes(1)
	assert.NotNil(t, attrs)
	assert.Equal(t, 4, len(attrs))
	assert.Equal(t, "test attr", attrs[3].GetContent()) // New attribute should be at the end

	// Test invalid section
	attrs = g.GetAttributes(-1)
	assert.Nil(t, attrs)
	attrs = g.GetAttributes(len(g.attributes))
	assert.Nil(t, attrs)
}

// Setter
func TestSetPoint(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)

	assert.NoError(t, g.SetPoint(utils.Point{X: 2, Y: 2}))
	assert.Equal(t, utils.Point{X: 2, Y: 2}, g.GetPoint())

	assert.NoError(t, err)
	assert.NoError(t, g.SetPoint(utils.Point{X: 3, Y: 3}))
	assert.Equal(t, utils.Point{X: 3, Y: 3}, g.GetPoint())
	assert.Equal(t, 2, mp.Times)
}

func TestSetLayer(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, g.SetLayer(1))
	assert.Equal(t, 1, g.GetLayer())

	assert.NoError(t, err)
	assert.NoError(t, g.SetLayer(2))
	assert.Equal(t, 2, g.GetLayer())
	assert.Equal(t, 2, mp.Times)
}

func TestSetColor(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, g.SetColor("#FF0000"))
	assert.Equal(t, "#FF0000", g.GetColor())

	assert.NoError(t, err)
	assert.NoError(t, g.SetColor("#00FF00"))
	assert.Equal(t, "#00FF00", g.GetColor())
	assert.Equal(t, 2, mp.Times)
}

func TestSetAttrContent(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)

	// Add an attribute to test
	assert.NoError(t, g.AddAttribute(0, "test"))

	// Test valid case
	assert.NoError(t, g.SetAttrContent(0, 0, "updated"))

	// Test invalid section
	assert.Error(t, g.SetAttrContent(-1, 0, "invalid"))
	assert.Error(t, g.SetAttrContent(len(g.attributes), 0, "invalid"))

	// Test invalid index
	assert.Error(t, g.SetAttrContent(0, -1, "invalid"))
	assert.Error(t, g.SetAttrContent(0, len(g.attributes[0]), "invalid"))

	// Test with parent draw update
	assert.NoError(t, g.SetAttrContent(0, 0, "updated again"))
	assert.Equal(t, 3, mp.Times)
}
func TestSetAttrSize(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)

	// Add an attribute to test
	assert.NoError(t, g.AddAttribute(0, "test"))

	// Test valid case
	assert.NoError(t, g.SetAttrSize(0, 0, 16))

	// Test invalid section
	assert.Error(t, g.SetAttrSize(-1, 0, 16))
	assert.Error(t, g.SetAttrSize(len(g.attributes), 0, 16))

	// Test invalid index
	assert.Error(t, g.SetAttrSize(0, -1, 16))
	assert.Error(t, g.SetAttrSize(0, len(g.attributes[0]), 16))

	// Test with parent draw update
	assert.NoError(t, g.SetAttrSize(0, 0, 18))
	assert.Equal(t, 3, mp.Times)
}
func TestSetAttrStyle(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)
	// Add an attribute to test
	assert.NoError(t, g.AddAttribute(0, "test"))

	// Test valid cases
	assert.NoError(t, g.SetAttrStyle(0, 0, int(attribute.Bold)))
	assert.NoError(t, g.SetAttrStyle(0, 0, int(attribute.Italic)))
	assert.NoError(t, g.SetAttrStyle(0, 0, int(attribute.Underline)))
	assert.NoError(t, g.SetAttrStyle(0, 0, int(attribute.Bold|attribute.Italic)))

	// Test invalid section
	assert.Error(t, g.SetAttrStyle(-1, 0, int(attribute.Bold)))
	assert.Error(t, g.SetAttrStyle(len(g.attributes), 0, int(attribute.Bold)))

	// Test invalid index
	assert.Error(t, g.SetAttrStyle(0, -1, int(attribute.Bold)))
	assert.Error(t, g.SetAttrStyle(0, len(g.attributes[0]), int(attribute.Bold)))

	// Test with parent draw update
	assert.NoError(t, g.SetAttrStyle(0, 0, int(attribute.Bold|attribute.Underline)))
	assert.Equal(t, 6, mp.Times)
}

// Methods
func TestCover(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	width := g.GetDrawData().(drawdata.Gadget).Width
	height := g.GetDrawData().(drawdata.Gadget).Height

	// corners
	corners := []utils.Point{
		{X: 1, Y: 1},                  // top-left
		{X: 1 + width, Y: 1},          // top-right
		{X: 1, Y: 1 + height},         // bottom-left
		{X: 1 + width, Y: 1 + height}, // bottom-right
	}
	for _, corner := range corners {
		val, err := g.Cover(corner)
		assert.NoError(t, err)
		assert.True(t, val)
	}

	// inside
	val, err := g.Cover(utils.Point{X: 1 + width/2, Y: 1 + height/2})
	assert.NoError(t, err)
	assert.True(t, val)

	// outside
	outsides := []utils.Point{
		{X: 1, Y: 0},
		{X: 1, Y: 1 + height + 1},
		{X: 0, Y: 1},
		{X: 1 + width + 1, Y: 1},
	}
	for _, outside := range outsides {
		val, err := g.Cover(outside)
		assert.NoError(t, err)
		assert.False(t, val)
	}
}

func TestAddAttribute(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)
	// Get initial attribute lengths
	initialLengths := g.GetAttributesLen()

	// Test adding attributes to different sections
	assert.NoError(t, g.AddAttribute(0, "test0"))
	assert.NoError(t, g.AddAttribute(1, "test1"))
	assert.NoError(t, g.AddAttribute(2, "test2"))

	// Verify lengths increased
	newLengths := g.GetAttributesLen()
	for i := 0; i < len(initialLengths); i++ {
		assert.Equal(t, initialLengths[i]+1, newLengths[i])
	}

	// Test invalid section
	assert.Error(t, g.AddAttribute(-1, "invalid"))
	assert.Error(t, g.AddAttribute(len(g.attributes), "invalid"))

	// Test with parent draw update
	assert.NoError(t, g.AddAttribute(0, "test with parent"))
	assert.Equal(t, 4, mp.Times)

	// Test with invalid content (this depends on attribute.NewAttribute implementation)
	// We can't easily test this without knowing what makes content invalid
	// But we can at least call it with empty string to see if it works
	assert.NoError(t, g.AddAttribute(0, ""))
}

func TestRemoveAttribute(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)

	// Add attributes to test removal
	assert.NoError(t, g.AddAttribute(0, "test0"))
	assert.NoError(t, g.AddAttribute(0, "test1"))
	assert.NoError(t, g.AddAttribute(1, "test2"))

	// Get lengths before removal
	beforeLengths := g.GetAttributesLen()

	// Test removing attributes
	assert.NoError(t, g.RemoveAttribute(0, 0))

	// Verify length decreased
	afterLengths := g.GetAttributesLen()
	assert.Equal(t, beforeLengths[0]-1, afterLengths[0])
	assert.Equal(t, beforeLengths[1], afterLengths[1]) // Unchanged

	// Test invalid section
	assert.Error(t, g.RemoveAttribute(-1, 0))
	assert.Error(t, g.RemoveAttribute(len(g.attributes), 0))

	// Test invalid index
	assert.Error(t, g.RemoveAttribute(0, -1))
	assert.Error(t, g.RemoveAttribute(0, len(g.attributes[0])))

	// Test with parent draw update
	assert.NoError(t, g.RemoveAttribute(1, 0)) // Remove the attribute we added to section 1
	assert.Equal(t, 5, mp.Times)
}

func TestGetDrawData(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	err := g.RegisterUpdateParentDraw(mp.UpdateParentDraw)
	assert.NoError(t, err)
	// Get draw data
	data := g.GetDrawData()

	// Verify it's the correct type
	gadgetData, ok := data.(drawdata.Gadget)
	assert.True(t, ok)

	// Verify basic properties
	assert.Equal(t, int(Class), gadgetData.GadgetType)
	assert.Equal(t, 1, gadgetData.X)
	assert.Equal(t, 1, gadgetData.Y)
	assert.Equal(t, 0, gadgetData.Layer)
	assert.Equal(t, drawdata.DefaultGadgetColor, gadgetData.Color)

	// Verify attributes
	assert.Equal(t, len(g.attributes), len(gadgetData.Attributes))
	for i, attrSection := range g.attributes {
		assert.Equal(t, len(attrSection), len(gadgetData.Attributes[i]))
	}

	// Test after modification
	assert.NoError(t, g.SetPoint(utils.Point{X: 2, Y: 2}))
	assert.NoError(t, g.SetLayer(1))
	assert.NoError(t, g.SetColor("#FF0000"))

	// Get updated draw data
	updatedData := g.GetDrawData()
	updatedGadgetData, ok := updatedData.(drawdata.Gadget)
	assert.True(t, ok)

	// Verify updated properties
	assert.Equal(t, 2, updatedGadgetData.X)
	assert.Equal(t, 2, updatedGadgetData.Y)
	assert.Equal(t, 1, updatedGadgetData.Layer)
	assert.Equal(t, "#FF0000", updatedGadgetData.Color)
}

func TestRegisterUpdateParentDraw(t *testing.T) {
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})
	mp := mockParent{}
	assert.NoError(t, g.RegisterUpdateParentDraw(mp.UpdateParentDraw))
	assert.Equal(t, 0, mp.Times)

	assert.NoError(t, g.AddAttribute(0, "test"))
	assert.Equal(t, 1, mp.Times)

	// nil function
	assert.Error(t, g.RegisterUpdateParentDraw(nil))
}

func TestValidateSection(t *testing.T) {
	// Create a gadget for testing validation methods
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})

	tests := []struct {
		name     string
		section  int
		hasError bool
	}{
		{"ValidSection", 0, false},
		{"ValidSectionMiddle", 1, false},
		{"ValidSectionLast", 2, false},
		{"NegativeSection", -1, true},
		{"SectionEqualToNumSections", len(g.attributes), true},
		{"SectionGreaterThanNumSections", len(g.attributes) + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := g.validateSection(tt.section)
			if tt.hasError {
				assert.Error(t, err)
				assert.IsType(t, duerror.NewInvalidArgumentError(""), err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateIndex(t *testing.T) {
	// Create a gadget for testing validation methods
	g := newEmptyGadget(Class, utils.Point{X: 1, Y: 1})

	// Test for each section
	for section, attrLen := range g.GetAttributesLen() {
		sectionName := fmt.Sprintf("Section%d", section)

		t.Run(sectionName, func(t *testing.T) {
			tests := []struct {
				name     string
				index    int
				hasError bool
			}{
				{"ValidFirstIndex", 0, false},
				{"ValidMiddleIndex", attrLen / 2, false},
				{"ValidLastIndex", attrLen - 1, false},
				{"NegativeIndex", -1, true},
				{"IndexEqualToCount", attrLen, true},
				{"IndexGreaterThanCount", attrLen + 1, true},
			}

			for _, tt := range tests {
				testName := fmt.Sprintf("%s_%s", sectionName, tt.name)
				t.Run(testName, func(t *testing.T) {
					err := g.validateIndex(tt.index, section)
					if tt.hasError {
						assert.Error(t, err)
						assert.IsType(t, duerror.NewInvalidArgumentError(""), err)
					} else {
						assert.NoError(t, err)
					}
				})
			}
		})
	}
}
