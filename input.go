package models

// Input for orders
import (
	"bytes"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	ID uint `gorm:"primarykey"`

	EdiNumber string //Numbering like in the edifact spec
	Label     string //readable name e.g. "Supplier Org Code"

	Type InputType

	ParentID         *uint
	Children         []Input `gorm:"foreignkey:ParentID"` // Children are further Inputs that can only be checked if this input is checked.
	ChildrenAreModal bool    // ChildrenAreModal is true when the children section should be displayed in a modal.

	Checked bool // Checked [only applies for Type CheckBox], wheter or not the input is checked

	Value    string // Value [Only applies for Type TextBox], the value of the TextBox
	Disabled bool   // Disabled [Only applies for Type TextBox], if the textbox is editable

	InfoField string // InfoField is a comment on the Input that can be displayed as a modal.

	SectionID uint // The id of the section this input belongs to.

	PrintTemplate string // The template to use when printing this input.

	Process string // Not deeper usage. Just for good overview
}

func (i Input) Child(id uint) Input {
	for _, c := range i.Children {
		if c.ID == id {
			return c
		}
	}
	return Input{}
}

type InputType string

const (
	CheckBoxType InputType = "checkbox"
	TextBoxType  InputType = "text"
)

type templateData struct {
	General General
	I       Input
	Now     time.Time
	OrderID string
}

func (i Input) GetLine(g General) string {
	if i.PrintTemplate == "" {
		return ""
	}
	parsed, err := template.New("input").Parse(i.PrintTemplate)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	var b bytes.Buffer
	err = parsed.Execute(&b,
		templateData{
			General: g,
			I:       i,
			Now:     time.Now(),
			OrderID: uuid.New().String()[:16],
		},
	)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return i.EdiNumber + "\t" + b.String() + "\n"
}
