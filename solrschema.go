package solrg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (sc *SolrClient) FieldTypes(collection string) (*FieldTypesResponse, error) {

	url := "http://" + sc.LBNodeAddress() + "/" + collection + "/schema/fieldtypes"
	req, err := http.NewRequest("GET", url, nil)
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var fieldTypeResp FieldTypesResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &fieldTypeResp)
	if err != nil {
		return nil, err
	}
	return &fieldTypeResp, nil
}

type FieldTypesResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	FieldTypes []struct {
		Name          string `json:"name"`
		Class         string `json:"class"`
		IndexAnalyzer struct {
			Tokenizer struct {
				Class string `json:"class"`
			} `json:"tokenizer"`
		} `json:"indexAnalyzer,omitempty"`
		QueryAnalyzer struct {
			Tokenizer struct {
				Class     string `json:"class"`
				Delimiter string `json:"delimiter"`
			} `json:"tokenizer"`
		} `json:"queryAnalyzer,omitempty"`
		SortMissingLast bool `json:"sortMissingLast,omitempty"`
		MultiValued     bool `json:"multiValued,omitempty"`
		Indexed         bool `json:"indexed,omitempty"`
		Stored          bool `json:"stored,omitempty"`
		Analyzer        struct {
			Tokenizer struct {
				Class string `json:"class"`
			} `json:"tokenizer"`
			Filters []struct {
				Class   string `json:"class"`
				Encoder string `json:"encoder"`
			} `json:"filters"`
		} `json:"analyzer,omitempty"`
		DocValues                 bool   `json:"docValues,omitempty"`
		Geo                       string `json:"geo,omitempty"`
		MaxDistErr                string `json:"maxDistErr,omitempty"`
		DistErrPct                string `json:"distErrPct,omitempty"`
		DistanceUnits             string `json:"distanceUnits,omitempty"`
		PositionIncrementGap      string `json:"positionIncrementGap,omitempty"`
		SubFieldSuffix            string `json:"subFieldSuffix,omitempty"`
		Dimension                 string `json:"dimension,omitempty"`
		AutoGeneratePhraseQueries string `json:"autoGeneratePhraseQueries,omitempty"`
	} `json:"fieldTypes"`
}
