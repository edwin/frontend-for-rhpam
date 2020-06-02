package helper

import "encoding/xml"

type VariableInstance struct {
	XMLName           xml.Name `xml:"variable-instance"`
	Name              string   `xml:"name"`
	OldValue          string   `xml:"old-value"`
	Value             string   `xml:"value"`
	ProcessInstanceId string   `xml:"process-instance-id"`
	ModificationDate  string   `xml:"modification-date"`
}

type VariableInstanceList struct {
	XMLName          xml.Name           `xml:"variable-instance-list"`
	VariableInstance []VariableInstance `xml:"variable-instance"`
}
