package ais_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ilder-as/go-barentswatch-ais/ais"
	"github.com/ilder-as/go-barentswatch-ais/ais/option"
	"github.com/ilder-as/go-barentswatch-ais/responsetype"
	"golang.org/x/oauth2"
)

func oauthSpoofMW(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token") {
			w.Header().Add("Content-Type", "application/json")
			token, _ := json.Marshal(oauth2.Token{AccessToken: "x", TokenType: "Bearer", Expiry: time.Now().Add(time.Hour)})
			w.Write(token)
			return
		}
		handlerFunc(w, r)
	}
}

func server(t *testing.T, handlerFunc http.HandlerFunc) *httptest.Server {
	sv := httptest.NewUnstartedServer(handlerFunc)
	sv.EnableHTTP2 = true
	sv.Start()

	return sv
}

// fixtureServer serves the content of the file as the response payload, Content-Type application/json
func fixtureServer(t *testing.T, filename string) *httptest.Server {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unable to load test fixture \"%s\": %s", filename, err)
	}

	buf := bytes.Buffer{}
	if _, err = io.Copy(&buf, f); err != nil {
		t.Fatalf("unable to copy data from file: %s", err)
	}
	f.Close()

	sv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token") {
			w.Header().Add("Content-Type", "application/json")
			token, _ := json.Marshal(oauth2.Token{AccessToken: "x", TokenType: "Bearer", Expiry: time.Now().Add(time.Hour)})
			w.Write(token)
			return
		}
		io.Copy(w, &buf)
	}))
	sv.EnableHTTP2 = true
	sv.Start()

	return sv
}

// hijackServer serves the content of the provided file byte for byte as the HTTP response. It must therefore
// include all headers etc.
func hijackServer(t *testing.T, filename string) *httptest.Server {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unable to load test fixture \"%s\": %s", filename, err)
	}

	buf := bytes.Buffer{}
	if _, err = io.Copy(&buf, f); err != nil {
		t.Fatalf("unable to copy data from file: %s", err)
	}
	f.Close()

	sv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token") {
			w.Header().Add("Content-Type", "application/json")
			token, _ := json.Marshal(oauth2.Token{AccessToken: "x", TokenType: "Bearer", Expiry: time.Now().Add(time.Hour)})
			w.Write(token)
			return
		}
		c, _, err := http.NewResponseController(w).Hijack()
		if err != nil {
			t.Fatalf("error in response controller: %s", err)
		}
		io.Copy(c, &buf)
		c.Close()
	}))
	sv.EnableHTTP2 = true
	sv.Start()

	return sv
}

func Test_GetAis(t *testing.T) {
	filename := "testdata/get_ais.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.GetAis()
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	numResults := 0
	for a := range ch {
		if a.IsZero() {
			t.Error("found zero record")
		}
		numResults++
	}

	if numResults <= 0 {
		t.Errorf("expected to find results, found %d", numResults)
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func Test_PostAis_NoneIncluded(t *testing.T) {
	filename := "testdata/post_ais_none_included_all.txt"
	sv := hijackServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostAis(ais.FilterInput{})
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()

	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	numResults := 0
	for range ch {
		numResults++
	}

	if numResults > 0 {
		t.Errorf("expected to find results, found %d", numResults)
	}

	if err = stream.Error(); err != nil {
		t.Fatalf("expected error, %s", err)
	}
}

func Test_PostCombined_FullGeojson(t *testing.T) {
	filename := "testdata/combined_full_geojson.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostCombined(ais.CombinedFilterInput{})
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	numResults := 0
	for a := range ch {
		if a.Type != responsetype.FullGeojson {
			t.Error("message type not recognized as FullGeojson")
		} else if a.AsFullGeojson().IsZero() {
			t.Error("message unmarshalled to empty message")
		}
		numResults++
	}

	if numResults <= 0 {
		t.Errorf("expected to find results, found %d", numResults)
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func Test_PostCombined_FullJson(t *testing.T) {
	filename := "testdata/combined_full_json.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostCombined(ais.CombinedFilterInput{})
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	for a := range ch {
		if a.Type != responsetype.FullJson {
			t.Error("message type not recognized as FullJson")
		} else if a.AsFullJson().IsZero() {
			t.Error("message unmarshalled to empty message")
		}
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func Test_PostCombined_SimpleGeoson(t *testing.T) {
	filename := "testdata/combined_simple_geojson.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostCombined(ais.CombinedFilterInput{})
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	for a := range ch {
		if a.Type != responsetype.SimpleGeojson {
			t.Error("message type not recognized as SimpleGeojson")
		} else if a.AsSimpleGeojson().IsZero() {
			t.Error("message unmarshalled to empty message")
		}
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func Test_PostCombined_SimpleJson(t *testing.T) {
	filename := "testdata/combined_simple_json.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostCombined(ais.CombinedFilterInput{})
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	for a := range ch {
		if a.Type != responsetype.SimpleJson {
			t.Error("message type not recognized as SimpleJson")
		} else if a.AsSimpleJson().IsZero() {
			t.Error("message unmarshalled to empty message")
		}
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func Test_PostCombined_Broken(t *testing.T) {
	filename := "testdata/combined_broken.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)
	_ = client

	stream, err := client.PostCombined(ais.CombinedFilterInput{})
	if err != nil {
		t.Fatal(err)
	}
	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	for range ch {
	}

	if err = stream.Error(); err == nil || ais.IsEOF(err) {
		t.Fatalf("expected non-EOF error, got \"%s\"", err)
	}
}

func Test_GetSSEAis(t *testing.T) {
	filename := "testdata/get_sse_ais.txt"
	sv := fixtureServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)

	stream, err := client.GetSSEAis()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := stream.UnmarshalStream()
	if err != nil {
		t.Fatal(err)
	}

	num := 0
	for data := range ch {
		if data.IsZero() {
			t.Error("got zero data")
		}
		num++
	}

	if num <= 0 {
		t.Errorf("expected > 0 results, got %d", num)
	}

	if err = stream.Error(); !ais.IsEOF(err) {
		t.Fatalf("expected EOF, got \"%s\"", err)
	}
}

func TestClient_GetLatestAis(t *testing.T) {
	filename := "testdata/get_latest_ais.txt"
	sv := hijackServer(t, filename)
	defer sv.Close()

	urls := ais.DefaultURLs()
	urls.OAuthBase = sv.URL
	urls.APIBase = sv.URL

	client := ais.NewClient("", "", urls)

	res, err := client.GetLatestAis(option.Since(time.Now()))
	if err != nil {
		t.Fatal(err)
	}

	data, err := res.Unmarshal()
	if err != nil {
		t.Fatalf("error unmarshalling: %s", err)
	}

	if len(data) <= 0 {
		t.Errorf("expected response to be non-empty list, list had %d members", len(data))
	}
}
