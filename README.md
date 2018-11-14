# Solrg

Solrg is a simple Go client for Apache Solr modeled after [Solrj](https://lucene.apache.org/solr/guide/7_4/using-solrj.html)

## Features

- Built-in load balancing (optional) - Uses ZooKeeper state to discover and route requests
- Simple API for the most commonly used Solr operations.

## Example Usage

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

Alternatively, if you want to connect directly to a Solr node, you can create a client using this:

```go
sc, err := solrg.NewDirectSolrClient("localhost:8983/solr")
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

For a full list of available request params, see [https://github.com/ezeev/solrg/blob/master/solrparams.go](https://github.com/ezeev/solrg/blob/master/solrparams.go). The current SolrParams struct doesn't cover every available request param by a long shot. I'll be adding more as I need them. PRs welcome.


## Roadmap

- Field collapsing
- More admin ops (schema crud, etc..)
- TBD...