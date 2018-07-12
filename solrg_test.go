package solrg_test

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/ezeev/solrg"
)

func must(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}

func TestSolrDocCollection(t *testing.T) {

	doc := solrg.NewSolrDocument("1")
	doc.SetField("test", []string{"test1", "test2", "test3"})
	docs := solrg.NewSolrDocumentCollection()
	err := docs.AddDoc(doc)
	must(err)
	t.Log(docs)

	//try getting the doc from the collection
	doc2 := docs.GetDoc("1")
	t.Log(doc2)

	//get one that doesn't exist
	doc3 := docs.GetDoc("2")
	t.Log(doc3)

	//now delete the doc
	docs.DeleteDoc("1")
	if docs.NumDocs() > 0 {
		t.Errorf("NumDocs() should return 0 but returned %d", docs.NumDocs())
	}

}

func TestSolrDocument(t *testing.T) {
	sd := solrg.NewSolrDocument("1")
	vals := []string{"string1", "string2", "string3"}
	sd.SetField("test", vals)

	//now get it
	vals2, err := sd.GetField("test")
	must(err)

	if reflect.DeepEqual(vals, vals2) {
		t.Log("vals == vals2 as expected")
	} else {
		t.Error("vals should equal vals2 but doesn't")
	}

	//try getting a field that doesn't exist
	_, err = sd.GetField("test2")
	if err != nil {
		t.Logf("Received error from GetField() as expected: %s", err)
	} else {
		t.Errorf("We should receive an error when trying to get a field that doesn't exist")
	}

	//test field exists
	if sd.Exists("test") {
		t.Log("test field exists as expected")
	} else {
		t.Error("field test should exist")
	}
	if !sd.Exists("test2") {
		t.Log("test2 does not exist as expected")
	} else {
		t.Error("field test2 should not exist")
	}

}

/*
func TestTempSolrSearchResp(t *testing.T) {

	url := "http://localhost:8983/solr/bb/select?q=ipad&fl=*,score&rows=25&defType=edismax&qf=keywords_txt_en&bq=sku_s:1945531^2.0%20sku_s:2339322^0.51%20sku_s:1945595^0.493%20sku_s:2842056^0.267%20sku_s:2339386^0.244%20sku_s:1945674^0.23%20sku_s:2408224^0.227%20sku_s:2842092^0.175%20sku_s:1918265^0.157%20sku_s:2817582^0.155%20sku_s:1918159^0.146%20sku_s:1918229^0.138%20sku_s:2538172^0.131%20sku_s:2475916^0.131%20sku_s:2809771^0.13%20sku_s:9947181^0.127%20sku_s:2319133^0.124%20sku_s:2701307^0.124%20sku_s:9924603^0.123%20sku_s:2343139^0.121%20sku_s:2678393^0.118%20sku_s:2205043^0.116%20sku_s:2197043^0.113%20sku_s:2319197^0.113%20sku_s:2903297^0.111%20sku_s:2490083^0.109%20sku_s:2884085^0.108%20sku_s:9635348^0.108%20sku_s:2339877^0.108%20sku_s:1151337^0.107%20sku_s:2318055^0.107%20sku_s:2678269^0.107%20sku_s:2874076^0.107%20sku_s:2343263^0.105%20sku_s:2340557^0.104%20sku_s:2390524^0.104%20sku_s:2343563^0.103%20sku_s:1114133^0.103%20sku_s:2339904^0.102%20sku_s:2861158^0.102%20sku_s:2629904^0.102%20sku_s:1609376^0.102%20sku_s:1271742^0.101%20sku_s:9809492^0.101%20sku_s:1286073^0.1%20sku_s:2330145^0.1%20sku_s:3396395^0.1%20sku_s:2330321^0.1%20sku_s:2642076^0.1%20sku_s:2883101^0.1&facet=true&facet.field=class_s&facet.field=on_sale_s&json.nl=arrntv"

	var sresp solrg.SolrSearchResponse

	var client = &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(buf, &sresp)

	//can we iterate through facets?
	t.Logf(sresp.ResponseHeader.Params.Bq)

	for k, v := range sresp.FacetCounts.FacetFields {
		t.Logf("key: %s", k)
		t.Log("values:")
		for _, val := range v {
			t.Logf("\t%s (%d)", val.Name, val.Value)
		}
	}

	//can we iterate through docs?
	t.Log("Showing docs")
	//t.Log(sresp.Response.Docs)
	for _, v := range sresp.Response.Docs {
		//v is the field
		for k, val := range v {
			switch assertVal := val.(type) {
			case float64:
				t.Logf("Numeric field, %s: %f", k, assertVal)
			case string:
				t.Logf("String field %s: %s", k, assertVal)
			case []interface{}:
				t.Logf("An array %s: %s", k, assertVal)
			}
		}

		t.Logf("Getting a string field, type_s: %s", v.String("type_s"))
		f, _ := v.Float64("reg_price_f")
		t.Logf("Getting a float field, reg_price_f: %f", f)

		arr, _ := v.Slice("cat_id_ss")
		t.Logf("Getting an array field, cat_id_ss: %s", arr)

		arr2, _ := v.StringSlice("cat_id_ss")
		t.Logf("Getting a string slice, cat_id_ss: %s", arr2)

	}

}
*/

func TestLBNodes(t *testing.T) {

	sc, err := solrg.NewSolrClient("localhost:9983")
	must(err)

	ln, _ := sc.LiveSolrNodes()

	t.Logf("Nodes: %s", ln)
	t.Logf("Call 1: %s", sc.LBNodeAddress())
	t.Logf("Call 2: %s", sc.LBNodeAddress())
	t.Logf("Call 3: %s", sc.LBNodeAddress())
	t.Logf("Call 4: %s", sc.LBNodeAddress())
	t.Logf("Call 5: %s", sc.LBNodeAddress())
	t.Logf("Call 6: %s", sc.LBNodeAddress())
}

func fakeDocs() solrg.SolrDocumentCollection {
	doc := solrg.NewSolrDocument("1")
	doc.SetField("test_txt", []string{"test1", "test2", "test3"})
	doc.SetField("test_s", []string{"test1"})

	doc2 := solrg.NewSolrDocument("2")
	doc2.SetField("test_txt", []string{"test3", "test4", "test5"})
	doc2.SetField("test_s", []string{"test2"})

	col := solrg.NewSolrDocumentCollection()
	col.AddDoc(doc)
	col.AddDoc(doc2)

	return col
}

func TestSolrDocJson(t *testing.T) {

	doc := solrg.NewSolrDocument("1")
	doc.SetField("test_txt", []string{"test1", "test2", "test3"})
	doc.SetField("test_s", []string{"test1"})

	jsn, _ := doc.SolrJSON()
	t.Logf("Json: %s", jsn)

	//doc collection
	col := solrg.NewSolrDocumentCollection()
	col.AddDoc(doc)

	doc2 := solrg.NewSolrDocument("2")
	doc2.SetField("test_txt", []string{"test3", "test4", "test5"})
	doc2.SetField("test_s", []string{"test2"})

	col.AddDoc(doc2)

	jsn, _ = col.SolrJSON()
	t.Logf("Collection Json string: %s", jsn)
}

func TestIndexDocs(t *testing.T) {

	sc, err := solrg.NewSolrClient("localhost:9983")
	must(err)

	// create test collection
	err = sc.CreateCollection("test", 1, 2, time.Second*180)
	if err != nil {
		t.Logf("An error was returned by collections API: %s", err.Error())
	}

	// give the cluster a few seconds
	time.Sleep(10 * time.Second)

	// index the docs
	docs := fakeDocs()
	//jsn, _ := docs.SolrJSON()

	err = sc.PostDocs(docs, "test")
	must(err)

	// commit the docs
	sc.Commit("test")

	// query all docs and also add a facet
	params := solrg.SolrParams{
		Q:          "*:*",
		Facet:      "true",
		FacetField: []string{"test_s", "test_txt"},
	}
	resp, err := sc.Query("test", "select", params, 10*time.Second)
	must(err)

	numDocs := resp.Response.NumFound
	if numDocs >= 2 {
		t.Logf("The search returned %d docs", numDocs)
		t.Logf("Here is the first doc: %s", resp.Response.Docs[0])
		t.Logf("Here is the second doc: %s", resp.Response.Docs[1])
	} else {
		t.Error("This search should return 2 docs!")
	}

	//were any facets returned?
	if len(resp.FacetCounts.FacetFields) != 2 {
		t.Error("Expected to get 2 facet fields back but did not!")
	}

	t.Logf("We queries w/ 2 facet fields. Here is the first facet %s", resp.FacetCounts.FacetFields["test_s"])
	t.Logf("Here is the second facet: %s", resp.FacetCounts.FacetFields["test_txt2"])

	//query again with 1 facets
	params.FacetField = []string{"test_s"}
	resp, err = sc.Query("test", "select", params, 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	if len(resp.FacetCounts.FacetFields) != 1 {
		t.Error("Expected to get 1 facet fields back but did not!")
	}

	//what are the facets now
	t.Logf("Query with 1 facets returned: %s", resp.FacetCounts.FacetFields["test_s"])

	// now query for a specific doc
	params.Q = "test_s:\"test1\""
	resp, err = sc.Query("test", "select", params, 10*time.Second)
	must(err)
	if resp.Response.NumFound == 1 {
		t.Log("Second search returned 1 doc as expected")
		t.Logf("The value of the test_s field for this doc is: %s", resp.Response.Docs[0].String("test_s"))
	}

	// now delete this doc
	sc.DeleteByQuery("test", "test_s:\"test1\"")
	sc.Commit("test")

	// how many docs are there now?
	params.Q = "*:*"
	resp, err = sc.Query("test", "select", params, 10*time.Second)
	must(err)

	if resp.Response.NumFound < numDocs {
		t.Log("There is one less document as expected")
	} else {
		t.Logf("The delete may have failed, the number of docs did not change after delete")
	}

	//now delete the collection
	err = sc.DeleteCollection("test")
	must(err)

}
func TestCreateDeleteCollection(t *testing.T) {

	sc, err := solrg.NewSolrClient("localhost:9983")
	must(err)

	err = sc.CreateCollection("test", 1, 2, time.Second*180)
	if err != nil {
		t.Error(err)
	}

	//now delete
	err = sc.DeleteCollection("test")
	if err != nil {
		t.Error(err)
	}

	//delete a collection that doesn't exist (we SHOULD get an error)
	err = sc.DeleteCollection("test2")
	if err == nil {
		t.Error("Expected an error but didn't get one!")
	} else {
		t.Logf("Received error as expected: %s", err)
	}

}

func TestZkConnect(t *testing.T) {

	sc, err := solrg.NewSolrClient("localhost:9983")
	must(err)
	t.Log(sc)
	ln, err := sc.LiveSolrNodes()
	must(err)
	t.Log(ln)

	liveNodes, err := sc.LiveSolrNodes()
	must(err)
	for _, n := range liveNodes.Nodes {
		t.Log(n)
	}

	// DISABLING TIMED THE TESTS BELOW - UNCOMMENT AND RUN THESE IF SETTING UP A NEW ENVIRONMENT WHERE
	// YOU WILL BE CHANGING ZK LOGIC
	t.Logf("Live nodes cache was last updated at %s", liveNodes.LastUpdate)
	ts1 := liveNodes.LastUpdate
	time.Sleep(time.Second * 1)
	liveNodes, err = sc.LiveSolrNodes()
	must(err)
	t.Logf("After 1 second, last updated was %s", liveNodes.LastUpdate)
	ts2 := liveNodes.LastUpdate
	time.Sleep(time.Second * 6)
	liveNodes, err = sc.LiveSolrNodes()
	must(err)
	t.Logf("After 7 seconds, last updated was %s", liveNodes.LastUpdate)
	ts3 := liveNodes.LastUpdate
	if ts1 == ts2 {
		t.Log("ts1 = ts2 as expected (1 second delay)")
	}
	if ts1 != ts3 {
		t.Log("ts1 != ts3 as expected (ts was changed after 5 seconds because live nodes was refreshed)")
	}

	// thread lock test 1
	go func() {
		t.Log("Thread test 1 starting")
		finished := false
		start := time.Now()
		for finished == false {
			sc.LiveSolrNodes()
			dur := time.Since(start)
			if dur.Seconds() > 7 {
				finished = true
				t.Log("Thread test 1 ran for 7 seconds")
			}
		}
	}()

	// thread lock test 12
	go func() {
		t.Log("Thread test 2 starting")
		finished := false
		start := time.Now()
		for finished == false {
			sc.LiveSolrNodes()
			dur := time.Since(start)
			if dur.Seconds() > 7 {
				finished = true
				t.Log("Thread test 2 ran for 7 seconds")
			}
		}
	}()
	//sleep this thread until the tests are finished
	time.Sleep(10 * time.Second)
}
