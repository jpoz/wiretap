# wiretap

A library to watch http traffic

## Use cases:

### Go project

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
