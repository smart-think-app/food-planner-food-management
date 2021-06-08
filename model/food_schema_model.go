package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type FoodSchemaModel struct {
	Name          string      `db:"name"`
	Id            int         `db:"id"`
	Status        int         `db:"status"`
	TypeFood      int         `db:"type_food"`
	CreatedDate   time.Time   `db:"created_date"`
	MaterialLevel interface{} `db:"material_level"`
	Image         string      `db:"image"`
	Mode          string      `db:"mode"`
}

type MaterialLevelSchemaModel struct {
	Protein int `json:"protein"`
	Fiber   int `json:"fiber"`
	Canxi   int `json:"canxi"`
	Fat     int `json:"fat"`
	Starch  int `json:"starch"`
}

func (a MaterialLevelSchemaModel) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *MaterialLevelSchemaModel) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}