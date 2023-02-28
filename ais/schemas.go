package ais

import (
	"github.com/ilder-as/go-barentswatch-ais/countrycode"
	"github.com/ilder-as/go-barentswatch-ais/modelformat"
	"github.com/ilder-as/go-barentswatch-ais/shiptype"
	geojson "github.com/paulmach/go.geojson"
	"time"
)

// The schema is available at https://live.ais.barentswatch.net/index.html, but note that some discrepancies between
// the documentation and actual behavior might be present. The implementation here is verified against the production
// system.
//
// Conventions/rules of thumb:
// - Nullable types are represented by pointer types.
// - Field names are annotated in full.
// - Native data types are used where it makes sense, as long as they marshal and unmarshal correctly.
// - Intermediary types (such as ShipType) can be introduced if an underlying semantic type exists,
//   even if not declared by the OpenAPI schema, if it makes the API easier to consume.
// - Int is used instead of int32.

// CombinedFilterInput is Ais.LiveApi.Api.AisMessage.Models.CombinedFilterInput
type CombinedFilterInput struct {
	Geometry     *geojson.Geometry         `json:"geometry"`
	Since        *time.Time                `json:"since"`
	MMSI         *int                      `json:"mmsi"`
	ShipTypes    []shiptype.ShipType       `json:"shipTypes"`
	CountryCodes []countrycode.CountryCode `json:"countryCodes"`
	ModelType    ModelType                 `json:"modelType"`
	ModelFormat  modelformat.ModelFormat   `json:"modelFormat"`
	Downsample   bool                      `json:"downsample"`
}

// FilterInput Ais.LiveApi.Api.AisMessage.Models.FilterInput
type FilterInput struct {
	Geometry                     *geojson.Geometry         `json:"geometry"`
	Since                        *time.Time                `json:"since"`
	MMSI                         []int                     `json:"mmsi"`
	ShipTypes                    []shiptype.ShipType       `json:"shipTypes"`
	CountryCodes                 []countrycode.CountryCode `json:"countryCodes"`
	IncludePosition              bool                      `json:"includePosition"`
	IncludeStatic                bool                      `json:"includeStatic"`
	IncludeAton                  bool                      `json:"includeAton"`
	IncludeSafetyRelated         bool                      `json:"includeSafetyRelated"`
	IncludeBinaryBroadcastMetHyd bool                      `json:"includeBinaryBroadcastMetHyd"`
	Downsample                   bool                      `json:"downsample"`
}

type LatestAisFilterInput struct {
	Geometry                     *geojson.Geometry         `json:"geometry"`
	Since                        *time.Time                `json:"since"`
	MMSI                         []int                     `json:"mmsi"`
	ShipTypes                    []shiptype.ShipType       `json:"shipTypes"`
	CountryCodes                 []countrycode.CountryCode `json:"countryCodes"`
	IncludePosition              bool                      `json:"includePosition"`
	IncludeStatic                bool                      `json:"includeStatic"`
	IncludeAton                  bool                      `json:"includeAton"`
	IncludeSafetyRelated         bool                      `json:"includeSafetyRelated"`
	IncludeBinaryBroadcastMetHyd bool                      `json:"includeBinaryBroadcastMetHyd"`
}

// ModelType is Ais.Shared.Models.Enums.ModelType
type ModelType string

const (
	ModelTypeSimple ModelType = "Simple"
	ModelTypeFull   ModelType = "Full"
)

// ApiError is an error supplied by the API.
//
// Note that the error type is undocumented in the Swagger documentation (as of 2023-02-28).
type ApiError struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	TraceId string `json:"traceId"`
}
