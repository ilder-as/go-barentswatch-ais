package ais_test

import (
	"bytes"
	"encoding/json"
	"github.com/ilder-as/go-barentswatch-ais/ais"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

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
		if !a.IsFullGeojson() {
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
		if !a.IsFullJson() {
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
		if !a.IsSimpleGeojson() {
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
		if !a.IsSimpleJson() {
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
