package es

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
)

func GetTypedClient() (*elasticsearch.TypedClient, error) {
	return elasticsearch.NewTypedClient(elasticsearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
}

func Add(index string, id string, doc interface{}) (*index.Response, error) {
	es, err := GetTypedClient()
	if err != nil {
		return nil, err
	}
	return es.Core.
		Index(index).
		Id(id).
		Request(doc).
		Do(context.Background())
}

func Delete(index string, id string) (*delete.Response, error) {
	es, err := GetTypedClient()
	if err != nil {
		return nil, err
	}
	return es.Core.
		Delete(index, id).
		Do(context.Background())
}
