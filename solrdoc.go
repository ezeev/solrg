package solrg

import (
	"encoding/json"
	"fmt"
)

// NewSolrDocument creates a new instance of a SolrDocument
func NewSolrDocument(id string) SolrDocument {
	sd := SolrDocument{}
	sd.fields = make(map[string][]string)
	sd.fields["id"] = []string{id}
	return sd
}

// SolrDocument struct the holds fields and provides methods for manipulating them
type SolrDocument struct {
	fields map[string][]string
}

// SetField sets the value for a field in the SolrDocument
func (sd *SolrDocument) SetField(name string, values []string) {
	sd.fields[name] = values
}

// GetField returns a field from the document if it exists
func (sd *SolrDocument) GetField(name string) ([]string, error) {
	if val, ok := sd.fields[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Document id %s does not contain field named %s", sd.fields["id"][0], name)
}

// Exists returns a bool. True = the field exists. False = the field does not exist.
func (sd *SolrDocument) Exists(name string) bool {
	_, present := sd.fields[name]
	return present
}

// SolrJSON returns the json representation of a solr document for indexing
func (sd *SolrDocument) SolrJSON() (string, error) {
	jsn, err := json.Marshal(sd.fields)
	return string(jsn), err
}

// ID returns the Id of the document
func (sd *SolrDocument) ID() string {
	if !sd.Exists("id") {
		return ""
	}
	return sd.fields["id"][0]
}

// SetID sets the Id of the document
func (sd *SolrDocument) SetID(id string) {
	sd.fields["id"] = []string{id}
}
