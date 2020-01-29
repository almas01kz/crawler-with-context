# How to use it

Open bin folder with all binaries. 

###  Execute on Linux and Mac
[-url] flag stands for url you want to crawl [-d] flag stands for depth
> Sitemap format is like simple folder hierarchy. Sitemap also shows the backlinks. 
```sh
$ chmod +x app
$  ./app -url https://cuvva.com -d 1

└───https://cuvva.com
        ├───https://cuvva.com/api
        ├───https://cuvva.com/car-insurance/learner-driver
        ├───https://cuvva.com/car-insurance/temporary-van-insurance
        ├───https://cuvva.com/car-insurance/subscription
        ├───https://cuvva.com/news
        ├───https://cuvva.com/car-insurance
        ├───https://cuvva.com/get-an-estimate
        ├───https://cuvva.com/about
        ├───https://cuvva.com/careers
        ├───https://cuvva.com/support
        ├───https://cuvva.com/car-insurance/temporary
        └───https://cuvva.com/single-trip-travel-insurance

```

##  Run and Build manually on Linux and Mac

> Go 1.11 and higher recommended.
 
 To run :
```sh
$ go run app.go -url https://kino.kz -d 3
```

 To build :
```sh
$ go get .
$ go build
$ chmod +x app
$  ./app -url https://cuvva.com -d 1
```

###  Run tests on Linux and Mac
```sh
$ go test -v
=== RUN   TestCrawl
--- PASS: TestCrawl (0.00s)
PASS
ok      /crawler-with-context/app      0.002s

```
