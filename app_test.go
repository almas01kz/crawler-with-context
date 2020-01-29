package main

import (
	"context"
	"testing"
)


const testTreeResult = `└───http://golang.org/
	└───http://golang.org/pkg/
		└───http://golang.org/cmd/
`

type testFetcher map[string]*testResult

type testResult struct {
	urls []string
}

func (f *testFetcher) Fetch(ctx context.Context,url string,baseUrl string) ([]string, error) {
	if res, ok := (*f)[url]; ok {
		return res.urls, nil
	}
	return nil,nil
}
// struct taken from GoTOUR
var fetcher = &testFetcher{
	"http://golang.org/": &testResult{
		[]string{
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/": &testResult{
		[]string{
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/cdm/": &testResult{
		[]string{
			"http://golang.org/",
		},
	},
}

func TestCrawl(t *testing.T) {

	baseUrl := "http://golang.org/"
	ctx := context.Background()
	Crawl(ctx,baseUrl, 3,  fetcher)
	res := dfsBuildTree(baseUrl, fetched.parentChldr)
	m:= make(map[string]interface{})
	m[baseUrl] = res
	got := ""
	buildTree(&got,m,"")

	if testTreeResult != got {
		t.Errorf("Expected \n%s. Got \n%s", testTreeResult, got)
	}

}
