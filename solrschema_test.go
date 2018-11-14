package solrg

import (
	"fmt"
	"testing"
)

func TestSolrFieldTypes(t *testing.T) {
	sc, err := NewSolrClient("localhost:9983")
	must(err)

	fmt.Println("hello")
	r, err := sc.FieldTypes("gettingstarted")
	if err != nil {
		t.Error(err)
	}

	//list field types
	for _, v := range r.FieldTypes {
		fmt.Printf("\tname: %s, analyzer: %s, all\n", v.Name, v.Analyzer)
	}

}
