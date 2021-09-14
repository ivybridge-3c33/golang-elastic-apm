package entities

import (
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(20)" json:"name"`
}

// func (e *Test) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(e)
// }

// func (e *Test) UnmarshalBinary(data []byte) error {
// 	return json.Unmarshal(data, e)
// }
