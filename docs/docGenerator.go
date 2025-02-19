package docs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type EndpointDescription struct {
	Name           string
	Method         string
	QueryParams    []string
	InputType      any
	OutputType     any
	Description    string
	Authentication bool
	ErrorCodes     []int
}

var docHandler *docsGenerator

func init() {
	docHandler = &docsGenerator{}
}

func GetDocsGenerator() DocsHandler {
	return docHandler
}

type DocsHandler interface {
	AddDoc(entpointDescription EndpointDescription)
	GenerateDocs() string
}

type docsGenerator struct {
	Data []EndpointDescription
}

func (dc *docsGenerator) AddDoc(entpointDescription EndpointDescription) {
	dc.Data = append(dc.Data, entpointDescription)
}

func (dc *docsGenerator) GenerateDocs() string {
	var sb strings.Builder
	for index, entry := range dc.Data {

		var InJson string
		var OutJson string
		dataIn, errJSONMarshal := json.Marshal(entry.InputType)
		if errJSONMarshal != nil {
			InJson = "Unable to Parse Input Type"
		}
		InJson = string(dataIn)

		dataOut, errJSONMarshal := json.Marshal(entry.OutputType)
		if errJSONMarshal != nil {
			InJson = "Unable to Parse Input Type"
		}
		sb.WriteString(fmt.Sprintf("Entpoint Nr. %v\n", +index+1))

		OutJson = string(dataOut)
		docEntry := fmt.Sprintf("Name: %v,\nMethod: %v,\nQueryParams: %v,\nAuthentification Required: %v,\nInput Body: %v,\nOutput Body: %v,\nErrorCodes: %v,\nDescription: %v\n",
			entry.Name,
			entry.Method,
			entry.QueryParams,
			entry.Authentication,
			InJson,
			OutJson,
			entry.ErrorCodes,
			entry.Description,
		)
		sb.WriteString(docEntry)
		sb.WriteString("\n")
	}
	return sb.String()
}
