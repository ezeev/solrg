# Solrg

Solrg is a simple Go client for Apache Solr modeled after [Solrj](https://lucene.apache.org/solr/guide/7_4/using-solrj.html)

## Features

- Built-in load balancing (optional) - Uses ZooKeeper state to discover and route requests
- Simple API for the most commonly used Solr operations.

## Indexing

```go
// Create a solr client
sc, err := solrg.NewSolrClient("localhost:9983")

// Create a collection
err = sc.CreateCollection("test", 1, 2, time.Second*180)

// Create a couple of documents
doc := solrg.NewSolrDocument("1")
doc.SetField("test_txt", []string{"test1", "test2", "test3"})
doc.SetField("test_s", []string{"test1"})

doc2 := solrg.NewSolrDocument("2")
doc2.SetField("test_txt", []string{"test3", "test4", "test5"})
doc2.SetField("test_s", []string{"test2"})

// Put them in a DocumentCollection
col := solrg.NewSolrDocumentCollection()
col.AddDoc(doc)
col.AddDoc(doc2)

// Index them
err = sc.PostDocs(&docs, "test")

// Commit changes
err = sc.Commit("test")
```

## Querying

```go
params := solrg.SolrParams{
    Q:          "*:*",
    Facet:      "true",
    FacetField: []string{"test_s", "test_txt"},
}
resp, err := sc.Query("test", "select", &params, 10*time.Second)
```

The Solr Response is serialized to structs located in [https://github.com/ezeev/solrg/blob/master/solrresp.go](https://github.com/ezeev/solrg/blob/master/solrresp.go)


## Roadmap

- Field collapsing
- More admin ops (schema crud, etc..)
- TBD...