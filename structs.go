package solrg

// CollectionsAPIResponse solr collections api response struct
type CollectionsAPIResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Success struct {
		One921681668983Solr struct {
			ResponseHeader struct {
				Status int `json:"status"`
				QTime  int `json:"QTime"`
			} `json:"responseHeader"`
			Core string `json:"core"`
		} `json:"192.168.1.66:8983_solr"`
		One921681667574Solr struct {
			ResponseHeader struct {
				Status int `json:"status"`
				QTime  int `json:"QTime"`
			} `json:"responseHeader"`
			Core string `json:"core"`
		} `json:"192.168.1.66:7574_solr"`
	} `json:"success"`
	Warning                        string `json:"warning"`
	OperationCreateCausedException string `json:"Operation create caused exception:"`
	Exception                      struct {
		Msg     string `json:"msg"`
		RspCode int    `json:"rspCode"`
	} `json:"exception"`
	Error struct {
		Metadata []string `json:"metadata"`
		Msg      string   `json:"msg"`
		Code     int      `json:"code"`
	} `json:"error"`
}
