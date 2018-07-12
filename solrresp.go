package solrg

import (
	"encoding/json"
	"fmt"
)

// SolrSearchResponse holds information about a Solr search response
type SolrSearchResponse struct {
	ResponseHeader struct {
		ZkConnected bool       `json:"zkConnected"`
		Status      int        `json:"status"`
		QTime       int        `json:"QTime"`
		Params      SolrParams `json:"params"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int                  `json:"numFound"`
		Start    int                  `json:"start"`
		MaxScore float64              `json:"maxScore"`
		Docs     []SolrSearchDocument `json:"docs"`
	} `json:"response"`
	FacetCounts struct {
		FacetQueries struct {
		} `json:"facet_queries"`
		FacetFields map[string][]SolrFacetField `json:"facet_fields"`
		FacetRanges struct {
		} `json:"facet_ranges"`
		FacetIntervals struct {
		} `json:"facet_intervals"`
		FacetHeatmaps struct {
		} `json:"facet_heatmaps"`
	} `json:"facet_counts"`
}

// SolrFacetField holds data for facet fields from a Solr response.
// NOTE: you must use &json.nl=arrntv on your Solr queries for this to work
type SolrFacetField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// SolrSearchDocument holds fields of a returned document and provides helper methods for accessing values
type SolrSearchDocument map[string]interface{}

// HasField returns true if the document has a specified field
func (sd SolrSearchDocument) HasField(fieldName string) bool {
	if _, ok := sd[fieldName]; ok {
		return true
	}
	return false
}

// String returns a string representation of a field
func (sd SolrSearchDocument) String(fieldName string) string {

	// make sure the field exists
	if !sd.HasField(fieldName) {
		return ""
	}

	s := sd[fieldName].(string)
	return s
}

// Float64 returns a float64 field
func (sd SolrSearchDocument) Float64(fieldName string) (float64, error) {
	f, ok := sd[fieldName].(float64)
	if ok {
		return f, nil
	}
	return 0, fmt.Errorf("Unable to assert float64 for field %s", fieldName)
}

// Int64 returns a int64 field or casts a float64 field to an int
func (sd SolrSearchDocument) Int64(fieldName string) (int64, error) {
	f, ok := sd[fieldName].(float64)
	if !ok {
		return 0, fmt.Errorf("Unable to assert float64 for field %s", fieldName)
	}
	return int64(f), nil
}

// Slice returns a slice (array) field
func (sd SolrSearchDocument) Slice(fieldName string) ([]interface{}, error) {
	f, ok := sd[fieldName].([]interface{})
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("Unable to assert []interface{} for field %s", fieldName)
}

// StringSlice returns a string slice (array) field
func (sd SolrSearchDocument) StringSlice(fieldName string) ([]string, error) {
	f, ok := sd[fieldName].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to assert []interface{} for field %s", fieldName)
	}
	var strSlice []string
	for _, v := range f {
		strSlice = append(strSlice, v.(string))
	}
	return strSlice, nil
}

// SolrParams hold information for a Solr request
type SolrParams struct {
	Q          string     `json:"q" url:"q,omitempty"`
	DefType    string     `json:"defType" url:"defType,omitempty"`
	FacetField FacetField `json:"facet.field" url:"facet.field,omitempty"`
	JSONNl     string     `json:"json.nl" url:"json.nl,omitempty"`
	Qf         string     `json:"qf" url:"qf,omitempty"`
	Fl         string     `json:"fl" url:"fl,omitempty"`
	Rows       string     `json:"rows" url:"rows,omitempty"`
	Facet      string     `json:"facet" url:"facet,omitempty"`
	Bq         string     `json:"bq" url:"bq,omitempty"`
}

type FacetField []string

// UnmarshalJSON an override because the Solr response can return a single value
// or a slice depending on how many facet fields are in the request. This gurantees
// that the FacetField part of the response serializes to our static type (a string slice)
func (ff *FacetField) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var v []string
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*ff = v
		return nil
	}
	*ff = []string{s}
	return nil
}
