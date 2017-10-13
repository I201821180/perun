package cfspecification

import (
	"encoding/json"
	"io/ioutil"
)

type Specification struct {
	PropertyTypes map[string]PropertyType
	ResourceSpecificationVersion string
	ResourceTypes map[string]Resource
}

type PropertyType struct {
	Documentation string
	Properties map[string]Property
}

type Property struct {
	Documentation string
	DuplicatesAllowed bool
	ItemType string
	PrimitiveItemType string
	PrimitiveType string
	Required bool
	Type string
	UpdateType string
}

type Resource struct {
	Documentation string
	Attributes map[string]Attribute
	Properties map[string]Property
}

type Attribute struct {
	ItemType string
	PrimitiveItemType string
	PrimitiveType string
	Type string
}

func GetSpecification(specificationFilePath string) (specification Specification, err error) {
	specificationFile, err := ioutil.ReadFile(specificationFilePath)
	if err != nil {
		panic(err)
	}

	return ParseSpecificationFile(specificationFile)
}

func ParseSpecificationFile(specificationFile []byte) (specification Specification, err error) {
	err = json.Unmarshal(specificationFile, &specification)
	if err != nil {
		return specification, err
	}

	return specification, nil
}