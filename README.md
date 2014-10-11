gBench - Apache Bench in Go
======
Mucking about with golang, my first ""helloworld"" learning project:   
A more user friendly Apache Bench CLI. 

## Building
`go build -o gb main.go`

## Usage

`./gb --url=google.ie -c=100 -n=2`

Or, just get prompted for the params.
	
	./gb
	URL to bench?
	<enter URL>
	Number of requests to perform?
	<enter a number>
	Concurrency (i.e. number of requests at one time)?
	<enter a number>
	Pinging  http://google.ie  with  2  requests, doing  2  of these at a time
	1...2...
	Avg time:  147.860332ms
	2 / 2  sucesses,  0 / 2  failures



You may want to bump ulimit to prevent "Error too many files open" when testing a local webserver:  
`ulimit -S -n 2048`

## Differences to Apache Bench

-Implements three possiobe params: URL, concurrency and number. No support (yet) for anything other than HTTP `GET`  
-Get prompted for params if omitted - no more remembering stuffs  
-Unlike Apache Bench, this really doesn't do very much. Not very much at all. 
