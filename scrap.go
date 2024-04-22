package main

import (
  "fmt"
  "flag"
  "encoding/json"
  "net/http"
  "net/url"
  "sync"
)

// scrap -url=https://localhost:8080/api/foo -params=[{'foo': 'bar'}, {'foo': 'buzz'}] -header={'access-token': 'foo'}

func parseArgToJson(stringValue string, resultMap *[]any) error {
  return json.Unmarshal([]byte(stringValue), resultMap)
}

func logError(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

func request(baseUrl string, param map[string]any, header map[string]any) {
  base, err := url.Parse(baseUrl)
  logError(err)
  client := &http.Client{}

  q := url.Values{}
  for key, value := range param {
    q.Add(key, value.(string))
  }

  base.RawQuery = q.Encode()
  req, _ := http.NewRequest("GET", base.String(), nil)

  for key, value := range header {
    req.Header.Add(key, value.(string))
  }

  res, err := client.Do(req)

  logError(err)

  fmt.Printf("%s - %d\n", base.String(), res.StatusCode)
  return
}

func main() {

  urlPtr := flag.String("url", "", "API url.")
  paramsPtr := flag.String("params", "[{\"foo\": \"bar\"}]", "List of API parameters on JSON format.")
  headerPtr := flag.String("header", "[]", "Request header.")

  flag.Parse()

  var paramsMap []map[string]any
  if err := json.Unmarshal([]byte(*paramsPtr), &paramsMap); err != nil {
    fmt.Println(err)
  }

  var headerMap map[string]any
  if err := json.Unmarshal([]byte(*headerPtr), &headerMap); err != nil {
    fmt.Println(err)
  }

  var wg sync.WaitGroup

  for i, s := range paramsMap {
    wg.Add(1)

    fmt.Printf("Launching request %d/%d\n", i+1, len(paramsMap))
    go func() {
      defer wg.Done()
      request(*urlPtr, s, headerMap)
    }()
  }

  wg.Wait()

  fmt.Println("Finished")
}
