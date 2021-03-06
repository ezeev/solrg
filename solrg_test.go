package solrg

import (
	"reflect"
	"testing"
	"time"
)

func TestSolrDocCollection(t *testing.T) {

	doc := NewSolrDocument("1")
	doc.SetField("test", []string{"test1", "test2", "test3"})
	docs := NewSolrDocumentCollection()
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

func TestSolrDirectClient(t *testing.T) {
	sc, err := NewDirectSolrClient("localhost:8983/solr")
	if err != nil {
		t.Error(err)
	}

	err = sc.Commit("gettingstarted")
	if err != nil {
		t.Log(err)
	}

}

func TestSolrPostStruct(t *testing.T) {
	sc, _ := NewSolrClient("localhost:9983")
	sc.DeleteCollection("test")
	err := sc.CreateCollection("test", 1, 2, 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	type TestStructDoc struct {
		Field1 string   `json:"field1_s"`
		Field2 int      `json:"field2_i"`
		Field3 float64  `json:"field3_f"`
		Field4 []string `json:"field4_ss"`
	}

	doc := TestStructDoc{
		Field1: "test",
		Field2: 10,
		Field3: 1.23,
		Field4: []string{"val1", "val2", "val3"},
	}

	docs := make([]interface{}, 1)
	docs[0] = doc

	err = sc.PostStructs(docs, "test")
	if err != nil {
		t.Error(err)
	}

	sc.Commit("test")

	//make sure it is there
	query, err := sc.Query("test", "select", &SolrParams{Q: "*:*"}, time.Second*10)
	if err != nil {
		t.Error(err)
	}
	if query.Response.NumFound == 0 {
		t.Fatalf("Expected at least one doc to be returned but there are %d", query.Response.NumFound)
	}

	err = sc.DeleteCollection("test")
	if err != nil {
		t.Error(err)
	}

}

func TestSolrCollectionAlreadyExists(t *testing.T) {

	sc, _ := NewSolrClient("localhost:9983")
	err := sc.CreateCollection("gettingstarted", 1, 2, 10*time.Second)

	serr, ok := err.(*SolrCollectionExistsError)
	if ok {
		t.Logf("Received SolrCollectionExistsError as expected. msg: %s", serr.Error())
	} else {
		t.Fail()
	}

}

func TestSolrDocument(t *testing.T) {
	sd := NewSolrDocument("1")
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

func TestLBNodes(t *testing.T) {

	sc, err := NewSolrClient("localhost:9983")
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

func fakeDocs() SolrDocumentCollection {
	doc := NewSolrDocument("1")
	doc.SetField("test_txt", []string{"test1", "test2", "test3"})
	doc.SetField("test_s", []string{"test1"})

	doc2 := NewSolrDocument("2")
	doc2.SetField("test_txt", []string{"test3", "test4", "test5"})
	doc2.SetField("test_s", []string{"test2"})

	col := NewSolrDocumentCollection()
	col.AddDoc(doc)
	col.AddDoc(doc2)

	return col
}

func TestSolrDocJson(t *testing.T) {

	doc := NewSolrDocument("1")
	doc.SetField("test_txt", []string{"test1", "test2", "test3"})
	doc.SetField("test_s", []string{"test1"})

	jsn, _ := doc.SolrJSON()
	t.Logf("Json: %s", jsn)

	//doc collection
	col := NewSolrDocumentCollection()
	col.AddDoc(doc)

	doc2 := NewSolrDocument("2")
	doc2.SetField("test_txt", []string{"test3", "test4", "test5"})
	doc2.SetField("test_s", []string{"test2"})

	col.AddDoc(doc2)

	jsn, _ = col.SolrJSON()
	t.Logf("Collection Json string: %s", jsn)
}

func TestIndexDocs(t *testing.T) {

	sc, err := NewSolrClient("localhost:9983")
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

	err = sc.PostDocs(&docs, "test")
	must(err)

	// commit the docs
	sc.Commit("test")

	// query all docs and also add a facet
	params := SolrParams{
		Q:          "*:*",
		Facet:      "true",
		FacetField: []string{"test_s", "test_txt"},
	}
	resp, err := sc.Query("test", "select", &params, 10*time.Second)
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

	t.Logf("We queries w/ 2 facet fields. Here is the first facet %v", resp.FacetCounts.FacetFields["test_s"])
	t.Logf("Here is the second facet: %v", resp.FacetCounts.FacetFields["test_txt2"])

	//query again with 1 facets
	params.FacetField = []string{"test_s"}
	resp, err = sc.Query("test", "select", &params, 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	if len(resp.FacetCounts.FacetFields) != 1 {
		t.Error("Expected to get 1 facet fields back but did not!")
	}

	//what are the facets now
	t.Logf("Query with 1 facets returned: %v", resp.FacetCounts.FacetFields["test_s"])

	// now query for a specific doc
	params.Q = "test_s:\"test1\""
	resp, err = sc.Query("test", "select", &params, 10*time.Second)
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
	resp, err = sc.Query("test", "select", &params, 10*time.Second)
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

	sc, err := NewSolrClient("localhost:9983")
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

	sc, err := NewSolrClient("localhost:9983")
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
