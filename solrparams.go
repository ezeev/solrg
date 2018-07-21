package solrg

import "encoding/json"

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
