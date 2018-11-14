package solrg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/samuel/go-zookeeper/zk"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// NewSolrClient returns a new instance of a Solr Client
func NewSolrClient(zksString string) (*SolrClient, error) {
	sc := SolrClient{}
	err := sc.Connect(zksString)
	return &sc, err
}

func NewDirectSolrClient(solrUrl string) (*SolrClient, error) {
	sc := SolrClient{}
	sc.liveNodes.Nodes = make([]string, 1)
	sc.liveNodes.Nodes[0] = solrUrl
	sc.numNodes = 1
	return &sc, nil
}

type SolrCollectionExistsError struct {
	collectionName string
}

func (e *SolrCollectionExistsError) Error() string {
	return fmt.Sprintf("Collection %s already exists", e.collectionName)
}

// SolrClient Solr Client struct
type SolrClient struct {
	liveNodes     LiveNodes
	lastNodeIndex int
	numNodes      int
	Connection    *zk.Conn
}

// LiveNodes struct to hold slice of live nodes and when the last time live nodes were updated
type LiveNodes struct {
	Nodes      []string
	LastUpdate time.Time
}

// Connect Connects the SolrClient instance
func (sc *SolrClient) Connect(zksString string) error {
	zks := strings.Split(zksString, ",")
	var err error
	sc.Connection, _, err = zk.Connect(zks, time.Second)
	return err
}

// LiveSolrNodes returns a slice of urls to live Solr nodes
func (sc *SolrClient) LiveSolrNodes() (*LiveNodes, error) {
	//only check for new nodes every 5 seconds
	duration := time.Since(sc.liveNodes.LastUpdate)
	if duration.Seconds() > 5 {
		ln, _, err := sc.Connection.Children("/live_nodes")
		if err != nil {
			log.Fatalf("Error getting live_nodes from zk: %s", err)
			return nil, err
		}
		sc.liveNodes.Nodes = make([]string, len(ln))
		//replace "_solr" with "/solr"
		for i, n := range ln {
			sc.liveNodes.Nodes[i] = strings.Replace(n, "_solr", "/solr", 1)
		}
		sc.liveNodes.LastUpdate = time.Now()
		sc.numNodes = len(sc.liveNodes.Nodes)
		sc.lastNodeIndex = 0
	}
	return &sc.liveNodes, nil
}

// Search executes a Solr search
func (sc *SolrClient) Query(collection string, reqHandler string, params *SolrParams, timeout time.Duration) (*SolrSearchResponse, error) {
	url := "http://" + sc.LBNodeAddress() + "/" + collection + "/" + reqHandler

	params.JSONNl = "arrntv"
	v, _ := query.Values(params)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(v.Encode())))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var client = &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error in Solr response, status code = %d, full error:\n%s", resp.StatusCode, body)
	}

	// we have a successful request if we made it this far
	var solrResp SolrSearchResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &solrResp)
	if err != nil {
		return nil, err
	}
	return &solrResp, nil

}

// LBNodeAddress Returns a node address using simple round robin LB of available nodes
func (sc *SolrClient) LBNodeAddress() string {
	// start back at 0
	if sc.numNodes == 0 {
		sc.LiveSolrNodes()
	}
	if sc.numNodes == 1 {
		return sc.liveNodes.Nodes[0]
	}
	//load balance
	if sc.lastNodeIndex == sc.numNodes-1 {
		//reset
		sc.lastNodeIndex = 0
	} else {
		sc.lastNodeIndex++
	}
	return sc.liveNodes.Nodes[sc.lastNodeIndex]

}

// Commit executes a Solr commit command
func (sc *SolrClient) Commit(collectionName string) error {

	//http://localhost:8983/solr/techproducts/update?commit=true
	url := "http://" + sc.LBNodeAddress() + "/" + collectionName + "/update?commit=true"

	var client = &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("Error executing commit command: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error executing commit command, status code = %d, full response:\n%s", resp.StatusCode, body)
	}

	return nil
}

type solrDeleteCommand struct {
	Delete struct {
		Query string `json:"query"`
	} `json:"delete"`
}

// DeleteByQuery deletes documents matching a Solr query
func (sc *SolrClient) DeleteByQuery(collectionName string, query string) error {

	url := "http://" + sc.LBNodeAddress() + "/" + collectionName + "/update"

	cmd := solrDeleteCommand{}
	cmd.Delete.Query = query

	jsn, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("Error marshalling delete query to json: %s", err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsn))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error delete docs by query, status code = %d, full error:\n%s", resp.StatusCode, body)
	}
	return nil
}

func (sc *SolrClient) PostStructs(data []interface{}, targetCollection string) error {
	url := "http://" + sc.LBNodeAddress() + "/" + targetCollection + "/update"
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//fmt.Println(string(dataBytes))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error indexing docs, status code = %d, full error:\n%s", resp.StatusCode, body)
	}
	return nil
}

// PostDocs indexes a SolrDocumentCollection
func (sc *SolrClient) PostDocs(docs *SolrDocumentCollection, targetCollection string) error {
	url := "http://" + sc.LBNodeAddress() + "/" + targetCollection + "/update"
	str, err := docs.SolrJSON()

	jsn := []byte("[" + str + "]")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsn))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Error indexing docs, status code = %d, full error:\n%s", resp.StatusCode, body)
	}
	return nil
}

// CreateCollection creates a Solr collection
func (sc *SolrClient) CreateCollection(name string, numShards int, replicationFactor int, timeout time.Duration) error {
	///admin/collections?action=CREATE&name=name
	url := fmt.Sprintf("http://%s/admin/collections?action=CREATE&name=%s&numShards=%d&replicationFactor=%d", sc.LBNodeAddress(), name, numShards, replicationFactor)

	var client = &http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	respCode := response.StatusCode
	var collectionsAPIResp CollectionsAPIResponse
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(buf, &collectionsAPIResp)

	if respCode == 200 {
		//success
		return nil
	} else if respCode == 400 {
		if strings.HasPrefix(collectionsAPIResp.Exception.Msg, "collection already exists") {
			return &SolrCollectionExistsError{collectionsAPIResp.Exception.Msg}
		}
		return fmt.Errorf("Error in CreateCollection(): %s", collectionsAPIResp.Exception.Msg)
	}
	return nil
}

// DeleteCollection deletes a Solr collection
func (sc *SolrClient) DeleteCollection(name string) error {
	///admin/collections?action=DELETE&name=collection
	url := fmt.Sprintf("http://%s/admin/collections?action=DELETE&name=%s", sc.LBNodeAddress(), name)
	var client = &http.Client{Timeout: time.Second * 10}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	respCode := response.StatusCode
	var collectionsAPIResp CollectionsAPIResponse
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(buf, &collectionsAPIResp)
	if respCode == 200 {
		return nil
	} else if respCode == 400 {
		return fmt.Errorf("Error in DeleteCollection(): %s", collectionsAPIResp.Exception.Msg)
	}
	return nil
}
