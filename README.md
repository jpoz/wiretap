# wiretap

A library to watch http traffic

## Use cases:

### Inside a Go project

```go
func main() {
	// Works like a normal http.Client but writes Request/Response to "./cache"
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}
	client := &http.Client{Transport: tr}

	resp, _ := client.Get("http://jsonip.com/")

	fmt.Printf("Returned %+s\n", resp.Status)

	fmt.Printf("Wrote to ./cache/jsonpi.com/GET/{TIMESTAMP}/response.txt")
}
```

### Proxy for a remote resource

```go
func director(r *http.Request) {
	r.URL.Scheme = "http"
	r.URL.Host = "github.com"
}

func main() {
	// Create wiretap Transport
	tr := &wiretap.Transport{
		Storage: disk.Storage{"./cache"},
	}

	// HTTP client with wiretap Transport
	client := &http.Client{Transport: tr}

	// Proxy
	proxy := wiretap.Proxy{
		Client:   client,
		Director: director,
	}

	http.HandleFunc("/", proxy.Handle)

	fmt.Println("localhost:8000 -> http://github.com/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```
