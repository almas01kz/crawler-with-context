package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
	"time"
)

const TIMEOUT  = time.Second * 5 // request timeout

type Fetcher interface {
	Fetch(ctx context.Context,url string,baseUrl string) (urls []string, err error)
}

// fetched tracks URLs that have been (or are being) fetched.
var fetched = struct {
	allUlrs     map[string]error
	parentChldr map[string]map[string]struct{}
	backLinks   map[string]map[string]struct{}
	sync.Mutex
}{allUlrs: make(map[string]error), parentChldr: make(map[string]map[string]struct{}), backLinks: make(map[string]map[string]struct{})}

var loading = errors.New("url load in progress")

// Crawl uses fetcher to recursively crawl
func Crawl(ctx context.Context,url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	fetched.Lock()
	if _, ok := fetched.allUlrs[url]; ok {
		fetched.Unlock()
		return
	}
	//Mark the url to be loading to avoid others reloading it at the same time.
	fetched.allUlrs[url] = loading
	fetched.Unlock()

	// Url Fetch concurrently.
	urls, err := fetcher.Fetch(ctx,url,url)

	// And update the status in a synced zone.
	fetched.Lock()

	fetched.allUlrs[url] = err
	m := make(map[string]struct{})
	backLinks := make(map[string]struct{})

	for _, u := range urls {
		if _, ok := fetched.allUlrs[u]; !ok {
			m[u] = struct{}{}
		} else {
			backLinks[u] = struct{}{}
		}
	}
	fetched.parentChldr[url] = m
	fetched.backLinks[url] = backLinks

	fetched.Unlock()

	if err != nil {
		fmt.Printf("<- Error on %v: %v\n", url, err)
		return
	}
	done := make(chan bool)
	for _, u := range urls {
		go func(url string) {
			Crawl(ctx,url, depth-1, fetcher)
			done <- true
		}(u)
	}

	for i := 0; i < len(urls); i++ {
		<-done
	}
}

func dfsBuildTree(name string, parentChild map[string]map[string]struct{}) map[string]interface{} {
	result := make(map[string]interface{})

	if _, ok := parentChild[name]; !ok {
		return make(map[string]interface{})
	}

	if _, ok := fetched.backLinks[name]; !ok {
		return make(map[string]interface{})
	} else {
		m := make(map[string]interface{})
		for k, _ := range fetched.backLinks[name] {
			m[k] = make(map[string]interface{})
		}
		result = m
	}

	for k, _ := range parentChild[name] {
		result[k] = dfsBuildTree(k, parentChild)
	}
	return result
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

type RealFetcher struct{}

func (f *RealFetcher) Fetch(ctx context.Context, url string,baseUrl string) ([]string, error) {

	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request error", err.Error())
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return []string{},nil
	}

	if resp == nil {
		return []string{},nil
	}

	z := html.NewTokenizer(resp.Body)

	var urls []string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return urls, nil
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			if t.Data != "a" {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			//Make sure the url begines in http**
			if strings.Index(url, baseUrl) == 0 {
				urls = append(urls, url)
			} else if strings.Index(url, "http") == -1 && len(url) > 1 && strings.Index(url, "/") == 0 {
				urls = append(urls, baseUrl+url)
			}
		}
	}

	return nil, nil
}

func buildTree(output *string, links map[string]interface{}, indent string) (*string, error) {
	iter := 0
	var nodePrefix, newIndent string
	for url, _ := range links {
		if iter == len(links)-1 {
			nodePrefix = "└───"
			newIndent = indent + "\t"
		} else {
			nodePrefix = "├───"
			newIndent = indent + "│\t"
		}

		*output += indent + nodePrefix + url + "\n"
		output, err := buildTree(output, links[url].(map[string]interface{}), newIndent)
		if err != nil {
			return output, err
		}
		iter++
	}
	return output, nil
}

func main() {
	ctx := context.Background()

	depth := flag.Int("d", 3, "depth")
	baseUrl := flag.String("url", "https://cuvva.com", "url to crawl")
	flag.Parse()

	fmt.Printf("Crawling %s  with depth: %d \n--------------", *baseUrl, *depth)

	Crawl(ctx,*baseUrl, *depth, &RealFetcher{})
	fmt.Println("\n Fetching stats & building hierarchy \n--------------")
	res := dfsBuildTree(*baseUrl, fetched.parentChldr)
	m := make(map[string]interface{})
	m[*baseUrl] = res
	output := ""
	buildTree(&output, m, "")
	fmt.Println(output)
}
