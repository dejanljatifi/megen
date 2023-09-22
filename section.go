package models

//sections for orders
import (
	"encoding/json"
)

type Section struct {
	ID uint `gorm:"primarykey"`

	Name    string
	Process string  // e.g. orders or delfor
	Inputs  []Input `gorm:"foreignKey:SectionID"`
}

type Sections []Section

func (s Sections) Json() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
