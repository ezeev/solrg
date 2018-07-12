package solrg

import (
	"fmt"
	"strings"
)

// SolrDocumentCollection holds a collection of SolrDocuments
type SolrDocumentCollection struct {
	docs map[string]SolrDocument
}

// NewSolrDocumentCollection returns a new instance of a SolrDocumentCollection
func NewSolrDocumentCollection() SolrDocumentCollection {
	sdc := SolrDocumentCollection{}
	sdc.docs = make(map[string]SolrDocument)
	return sdc
}

// AddDoc adds a document to the collection
func (sdc *SolrDocumentCollection) AddDoc(doc SolrDocument) error {
	if doc.ID() == "" {
		return fmt.Errorf("Document is missing an id! Please make sure all docs have an id field")
	}
	sdc.docs[doc.ID()] = doc
	return nil
}

// DeleteDoc removes a doc by id
func (sdc *SolrDocumentCollection) DeleteDoc(id string) {
	delete(sdc.docs, id)
}

// SolrJSON returns a json string representation of the doc collection ready for Solr
func (sdc *SolrDocumentCollection) SolrJSON() (string, error) {

	jsn := ""
	for i, v := range sdc.docs {
		djsn, err := v.SolrJSON()
		if err != nil {
			return "", fmt.Errorf("Error creating json string at position %s, error: %s", i, err)
		}
		jsn = jsn + djsn + ",\n"
	}
	if strings.HasSuffix(jsn, ",\n") {
		jsn = strings.TrimRight(jsn, ",\n")
	}
	return jsn, nil
}

// GetDoc returns a doc by id
func (sdc *SolrDocumentCollection) GetDoc(id string) SolrDocument {
	return sdc.docs[id]
}

// NumDocs returns the number of docs in the collection
func (sdc *SolrDocumentCollection) NumDocs() int {
	return len(sdc.docs)
}
