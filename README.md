# Barentswatch AIS 
## About 
This is a Go implementation of [Barentswatch' AIS API](https://live.ais.barentswatch.net/index.html#/), which supplies
live data from Automatic Identification Systems (AIS) from marine vessels along the Norwegian coast.

This Go API seeks to make it easy to ingest and use these data for further processing. The API is purposefully designed 
as a thin layer above the HTTP API, to make it easier to keep the implementation synchronized with upstream changes. 
This means that we e.g. sometimes implement endpoints whose supplied data are otherwise identical to each other 
whenever the HTTP API has endpoints with identical data but differences in carrier protocol or exchange format. 

## About Ilder Open Source
This project is proudly supported by [Ilder AS](https://ilder.no), a Norway-based software and IoT company.
Ilder acknowledges the significance of open-source software and open data in its own work and in the software industry 
as a whole. The company recognizes that the contributions made by the open-source community have been instrumental in 
advancing the field, and as such, Ilder is committed to making its own contributions to the software community and 
playing an active role in supporting open-source initiatives.

## Usage
In order to use this API, you must have a valid [Barentswatch account](https://www.barentswatch.no/minside/). 
You must create a new client, and obtain a client id and client secret. 

The package is available as a go module, and can be downloaded as `github.com/ilder-as/go-barentswatch-ais`.

You initialize the library as follows 

```go
// Replace client id and client secret with your own values
client := ais.NewClient("user@example.com:name", "clientsecret")
```

## Consuming the API
There are two kinds of API endpoint, those which return streaming data, and those which return fixed result sets. 
They can be recognized by the return types. 

### Streams 
A streaming endpoint will have a `StreamResponse[T]` return data type. The underlying stream can be a simple HTTP2 stream
or Server Sent Events (SSE). A typical signature is 

```go
func (c *Client) GetAis() (StreamResponse[AisMultiple], error) { // ... }
```

A `StreamResponse` has an `UnmarshalStream` method which returns channels for data, errors, and a cancellation function
which cancels the stream.

```go
// Replace client id and client secret with your own values
client := ais.NewClient("user@example.com:name", "clientsecret")

// A supplied context allows for cancellation of the request, and of the reading of the response stream
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Execute request
stream, err := client.GetAisContext(ctx)
if err != nil {
    panic(err)
}

dataCh, err := stream.UnmarshalStream()
if err != nil {
    panic(err)
}

for aisData := range dataCh {
    fmt.Println(aisData)
}

// If the channel closes, a supplied error will explain why
if err := stream.Error(); err != nil {
	fmt.Println(err)
}
```

### Queries 
Query responses are those API calls which have a `Response[T]` return type. These calls return simple data types or result sets 
as slices of simple data types. E.g. 

```go
func (c *Client) GetLatestAis(opts ...latestAisOption) (Response[[]AisMultiple], error) { // ... }
```

These responses have an `Unmarshal` method which unmarshals the data into the supplied data type `T`, like `[]AisMultiple` in the example above.

```go
// Replace client id and client secret with your own values
client := ais.NewClient("user@example.com:name", "clientsecret")

res, err := client.GetLatestAis()
if err != nil {
    panic(err)
}

latest, err := res.Unmarshal()
if err != nil {
    panic(err)
}

fmt.Println(latest)
```


