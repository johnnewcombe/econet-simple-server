package fs

import (
	"encoding/json"
	"fmt"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

type Catalogue struct {
	Cycle             byte             `json:"cycle" bson:"cycle"`
	Entries           []CatalogueEntry `json:"entries" bson:"entries"`
	catalogueFilePath string
}

type CatalogueEntry struct {
	Type   string `json:"type" bson:"type"`
	Name   string `json:"name" bson:"name"`
	Access string `json:"access" bson:"access"`
}

func (c *Catalogue) Load(jsonBytes []byte) error {

	if !json.Valid(jsonBytes) {
		return fmt.Errorf("validating catalogue: invalid json")
	}

	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return fmt.Errorf("parsing json: invalid")
	}

	return nil
}

//func (c *Catalogue) Dump() ([]byte, error) {
//
//	var (
//		data []byte
//		err  error
//	)
//
//	if data, err = json.Marshal(c); err != nil {
//		return nil, err
//	}
//
//	return data, nil
//}

func (c *Catalogue) saveToDisk() error {

	var (
		err error
	)

	if len(c.catalogueFilePath) > 0 {
		if err = lib.WriteString(c.catalogueFilePath, c.ToString()); err != nil {
			return err
		}
	}
	// write the userData to disk
	return nil
}

func (c *Catalogue) ToBytes() ([]byte, error) {

	var (
		jsonBytes []byte
		err       error
	)

	if jsonBytes, err = json.Marshal(c); err != nil {
		return []byte{}, err
	}

	return jsonBytes, nil
}

func (c *Catalogue) ToString() string {

	var (
		err error
		b   []byte
	)

	if b, err = c.ToBytes(); err != nil {
		return ""
	}
	return string(b)
}
