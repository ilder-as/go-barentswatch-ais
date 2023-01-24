# Barentswatch AIS 
## About 
This is a Go implementation of [Barentswatch' AIS API](https://live.ais.barentswatch.net/index.html#/), which supplies
live data from Automatic Identification Systems (AIS) from marine vessels along the Norwegian coast.

This Go API seeks to make it easy to ingest and use these data for further processing. 

This module is developed under the patronage of [Ilder AS](https://ilder.no). Ilder recognizes the importance of 
open software and open data in the modern software landscape. People and industries are indebted to the collective 
efforts that are made available for free (libre) use, and intends to take an active part in making contributions to the
software community.

## Usage
In order to use this API, you must have a valid [Barentswatch account](https://www.barentswatch.no/minside/). 
You must create a new client, and obtain a client id and client secret. 

You initialize the library as follows 

```go
// Replace client id and client secret with your own values
client := ais.NewClient("user@example.com:name", "clientsecret")
```

## Consuming the API
There are two kinds of API endpoint, those which return streaming data, and those which return fixed result sets. 
They can be recognized by the return types. 

### Streams 
A typical streaming endpoint will have a `*StreamResponse` or 
`SSEStreamResponse` return data type, depending on whether the data is streamed via HTTP2 streams or as Server Sent 
Events. A typical signature is 

```go
func (c *Client) GetAis() (*StreamResponse[AisMultiple], error) { // ... }
```

A stream response has an `UnmarshalStream` method which returns channels for data, errors, and a cancellation function
which cancels the stream.

```go
// Replace client id and client secret with your own values
client := ais.NewClient("user@example.com:name", "clientsecret")

stream, err := client.GetAis()
if err != nil {
    panic(err)
}

dataCh, _, _, err := stream.UnmarshalStream(context.Background())
if err != nil {
    panic(err)
}

for aisData := range dataCh {
    fmt.Println(aisData)
}
```

### Queries 
Query responses are those API calls which do not return streams. These calls return simple data types or result sets 
as slices of simple data types, and can be recognized by their `*Response` return type. E.g. 

```go
func (c *Client) GetLatestAis(opts ...latestAisOption) (*Response[[]AisMultiple], error) { // ... }
```

These responses have an `Unmarshal` method which unmarshals the data to a native data type of the correct form.

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


