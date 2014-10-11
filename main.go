package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "sync"
    "flag"
    "regexp"
    "math"
)

func doRequest(url string)(time.Duration, int){
  start := time.Now()
  response, err := http.Get(url)
  if (err != nil){
    // Bail on critical errors
    fmt.Println("Error reaching ", url, " : ", err)
    os.Exit(1);
  }
  defer response.Body.Close()
  elapsed := time.Since(start)
  return elapsed,response.StatusCode
}

func settingsFromInput()(string, int, int){
  var url string
  var c, n int
  r, _ := regexp.Compile("^http(s)?://")

  
  flag.IntVar(&c, "c", -1, "Concurrency of Requests")
  flag.IntVar(&n, "n", -1, "Number of Requests")
  flag.StringVar(&url, "url", "", "URL to ping")
  flag.Parse()
  
  // Scan from stdin if variables don't exist from flags
  if (url == ""){
    fmt.Println("URL to bench?")
    fmt.Scanf("%s", &url)
  }
  
  if (!r.MatchString(url)){
    url = "http://" + url
  }
  
  if (n == -1){
    fmt.Println("Number of requests to perform?")
    fmt.Scanf("%d", &n)
  }
  
  if (c == -1){
    fmt.Println("Concurrency (i.e. number of requests at one time)?")
    fmt.Scanf("%d", &c)
  }
  
  if (c > n){
    fmt.Println("Error - can't have a concurrency (", c, ") greater than the number of requests (", n, ")")
    os.Exit(1);
  }
  
  return url, c, n
}

func bench(url string, c, n int){
  fmt.Println("Pinging ", url, " with ", n, " requests, doing ", c, " of these at a time")
  
  durations := make(chan time.Duration)
  statuses := make(chan int)

  var wg sync.WaitGroup
  wg.Add(n)
  var concurrentgroups int = int(math.Ceil(float64(n / c)))
  var requestCounter int = 0
  for i := 0; i <=concurrentgroups; i++ {
    for j := 0; j < c && (((i) * c) + j) < n; j++ { // where (((i) * c) + j) is numRequest
      go func() {
        defer wg.Done()
        elapsed, status := doRequest(url)
        // Push latest duration & status code
        requestCounter++
        fmt.Printf("%d...", requestCounter)
        durations <- elapsed
        statuses <- status
      }()
    }
  }
  
  var totalDuration, averageDuration time.Duration
  var successes, failures int = 0, 0
  for i := 0; i < n; i++ {
      dur := <-durations
      status := <-statuses
      totalDuration += dur
      if status == 200{
        successes++;
      }else{
        failures++;
      }
  }
  wg.Wait()
  close(durations)
  
  // Calculate averages
  averageDuration = totalDuration/ time.Duration(n)
  fmt.Println("\nAvg time: ", averageDuration)
  fmt.Println(successes,"/",n," sucesses, ", failures, "/",n," failures")
}

func main() {
    var url string
    var c, n int
    url, c, n = settingsFromInput()
    
    bench(url, c, n)
    
}
