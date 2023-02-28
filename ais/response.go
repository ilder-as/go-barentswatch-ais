package ais

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ilder-as/go-barentswatch-ais/responsetype"
	"net/http"
	"reflect"
	"regexp"
	"time"
)

type StreamType int

const (
	None StreamType = iota
	Simple
	SSE
)

// eof is a sentinel error which indicates end of stream ("end of file")
var eof = errors.New("EOF")

// IsEOF returns true iff the supplied error is an EOF error (i.e. signals the end of a stream).
func IsEOF(err error) bool {
	return errors.Is(err, eof)
}

type CancelFunc func()

// Response is an API response which contains the given T data type.
type Response[T any] struct {
	*http.Response
}

// Unmarshal unmarshals the reponse body into a new object of type T.
func (r Response[T]) Unmarshal() (T, error) {
	var obj T
	defer r.Body.Close()
	return obj, json.NewDecoder(r.Body).Decode(&obj)
}

// StreamResponse is an API response whose body is a continuous stream of data.
//
// A StreamResponse can be consumed using the `UnmarshalStream` method.
type StreamResponse[T any] struct {
	*http.Response
	streamType StreamType
	err        error
	ctx        context.Context
}

// Error returns the underlying error or reason when a stream ends.
func (r *StreamResponse[T]) Error() error {
	return r.err
}

// unmarshalDefault unmarshals streaming data sent as individual json objects
func (r *StreamResponse[T]) unmarshalDefault() (<-chan T, error) {
	scan := bufio.NewScanner(r.Body)
	if scan == nil {
		return nil, errors.New("received nil scanner")
	}

	out := make(chan T)

	go func() {
		defer r.Body.Close()
		defer close(out)

		for {
			select {
			case <-r.ctx.Done():
				r.err = r.ctx.Err()
				return
			default:
				success := scan.Scan()
				if !success {
					if err := scan.Err(); err != nil {
						r.err = err
					} else {
						r.err = eof
					}
					return
				}

				var res T
				if err := json.Unmarshal(scan.Bytes(), &res); err != nil {
					r.err = err
					return
				} else {
					out <- res
				}
			}
		}
	}()

	return out, nil
}

// unmarshalSSE unmarshals streaming data sent as Server Sent Events (SSE)
func (r *StreamResponse[T]) unmarshalSSE() (<-chan T, error) {
	scan := bufio.NewScanner(r.Body)
	if scan == nil {
		return nil, errors.New("received nil scanner")
	}

	out := make(chan T)

	go func() {
		defer r.Body.Close()
		defer close(out)

		for {
			select {
			case <-r.ctx.Done():
				r.err = r.ctx.Err()
				return
			default:
				success := scan.Scan()
				if !success {
					if err := scan.Err(); err != nil {
						r.err = err
					} else {
						r.err = eof
					}
					return
				}

				var res T
				var err error
				res, err = unmarshalSSEData[T](scan.Bytes())
				if errors.Is(err, empty) || errors.Is(err, noMatch) {
					continue
				} else if err != nil {
					r.err = err
					return
				} else {
					out <- res
				}
			}
		}
	}()

	return out, nil
}

// UnmarshalStream unmarshals a stream of serialized data into the underlying data structure.
//
// The returned channel returns the next object unmarshalled. The channel only closes when it encounters an error,
// or when the stream closes. Use StreamResponse.Error to check the reason for the closed stream. If UnmarshalStream
// encounters an error, the underlying connection is closed. To continue consuming data, another api call must be made
// to get a new StreamResponse.
func (r *StreamResponse[T]) UnmarshalStream() (<-chan T, error) {
	switch r.streamType {
	case Simple:
		return r.unmarshalDefault()
	case SSE:
		return r.unmarshalSSE()
	default:
		return nil, errors.New("unknown stream type")
	}
}

var (
	noMatch = errors.New("no match")
	empty   = errors.New("empty")
)

// unmarshalSSEData unmarshals an SSE data stream
func unmarshalSSEData[T any](raw []byte) (T, error) {
	var res T

	if len(raw) == 0 {
		return res, empty
	}

	re := regexp.MustCompile(`^data:\s+(\{.*\})\s*$`)
	match := re.FindSubmatch(raw)
	if len(match) <= 1 {
		return res, noMatch
	}

	err := json.Unmarshal(match[1], &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// AisMultiple holds a union of the multiple response types that an AIS data request can return.
// Use the Type property to inspect which type the message is.
type AisMultiple struct {
	Type responsetype.Ais
	Position
	Aton
	Staticdata
}

func (a *AisMultiple) UnmarshalJSON(data []byte) error {
	typ := struct {
		Type responsetype.Ais `json:"type"`
	}{}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	a.Type = typ.Type

	switch a.Type {
	case responsetype.Position:
		return json.Unmarshal(data, &a.Position)
	case responsetype.Aton:
		return json.Unmarshal(data, &a.Aton)
	case responsetype.Staticdata:
		return json.Unmarshal(data, &a.Staticdata)
	default:
		return fmt.Errorf("unknown type: %s", a.Type)
	}
}

// IsZero is true iff the receiver is a default-valued AisMultiple struct.
func (a AisMultiple) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// AsPosition returns the underlying Position response data if the response is of the correct type,
// and a zero (default-valued) Position struct otherwise.
func (a AisMultiple) AsPosition() Position {
	return a.Position
}

// AsAton returns the underlying Aton response data if the response is of the correct type,
// and a zero (default-valued) Aton struct otherwise.
func (a AisMultiple) AsAton() Aton {
	return a.Aton
}

// AsStaticdata returns the underlying Staticdata response data if the response is of the correct type,
// and a zero (default-valued) Staticdata struct otherwise.
func (a AisMultiple) AsStaticdata() Staticdata {
	return a.Staticdata
}

// Position
//
// Verified in private communication Jan 24 2023.
type Position struct {
	MessageType        int       `json:"messageType"`
	Mmsi               int       `json:"mmsi"`
	Msgtime            time.Time `json:"msgtime"`
	Altitude           *int      `json:"altitude"`
	Longitude          *float64  `json:"longitude"`
	Latitude           *float64  `json:"latitude"`
	CourseOverGround   *float64  `json:"courseOverGround"`
	AisClass           string    `json:"aisClass"`
	NavigationalStatus int       `json:"navigationalStatus"`
	RateOfTurn         *float64  `json:"rateOfTurn"`
	SpeedOverGround    *float64  `json:"speedOverGround"`
	TrueHeading        *int      `json:"trueHeading"`
}

// IsZero is true iff the receiver is a default-valued Position struct.
func (a Position) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// Aton
//
// Verified in private communication Jan 24 2023.
type Aton struct {
	MessageType                  int       `json:"messageType"`
	Mmsi                         int       `json:"mmsi"`
	Msgtime                      time.Time `json:"msgtime"`
	Longitude                    *float64  `json:"longitude"`
	Latitude                     *float64  `json:"latitude"`
	Name                         string    `json:"name"`
	DimensionA                   *int      `json:"dimensionA"`
	DimensionB                   *int      `json:"dimensionB"`
	DimensionC                   *int      `json:"dimensionC"`
	DimensionD                   *int      `json:"dimensionD"`
	TypeOfAidsToNavigation       int       `json:"typeOfAidsToNavigation"`
	TypeOfElectronicFixingDevice int       `json:"typeOfElectronicFixingDevice"`
}

// IsZero is true iff the receiver is a default-valued Aton struct.
func (a Aton) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// Staticdata
//
// Verified in private communication Jan 24 2023.
type Staticdata struct {
	MessageType              int       `json:"messageType"`
	Mmsi                     int       `json:"mmsi"`
	Msgtime                  time.Time `json:"msgtime"`
	Name                     string    `json:"name"`
	DimensionA               *int      `json:"dimensionA"`
	DimensionB               *int      `json:"dimensionB"`
	DimensionC               *int      `json:"dimensionC"`
	DimensionD               *int      `json:"dimensionD"`
	ImoNumber                *int      `json:"imoNumber"`
	CallSign                 string    `json:"callSign"`
	Destination              string    `json:"destination"`
	Eta                      string    `json:"eta"`
	Draught                  *int      `json:"draught"`
	ShipLength               *int      `json:"shipLength"`
	ShipWidth                *int      `json:"shipWidth"`
	ShipType                 *int      `json:"shipType"`
	PositionFixingDeviceType int       `json:"positionFixingDeviceType"`
	ReportClass              string    `json:"reportClass"`
}

// IsZero is true iff the receiver is a default-valued Staticdata struct.
func (a Staticdata) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// CombinedMultiple is a response which can be either CombinedSimpleJson, CombinedFullJson, CombinedSimpleGeojson or
// CombinedFullGeojson. Which one it is depends on what was requested by the user, and must be checked on use.
type CombinedMultiple struct {
	Type responsetype.Combined `json:"type"`
	CombinedSimpleJson
	CombinedFullJson
	CombinedSimpleGeojson
	CombinedFullGeojson
}

// UnmarshalJSON unmarshals the supplied JSON data into a CombinedMultiple.
func (c *CombinedMultiple) UnmarshalJSON(data []byte) error {
	keys := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &keys); err != nil {
		return err
	}

	isGeojson := false
	isFull := false

	// The "properties" key exists iff message is Geojson
	_, isGeojson = keys["properties"]

	if !isGeojson {
		// "eta" exists only on the full type
		_, isFull = keys["eta"]
	} else {
		propKeys := make(map[string]json.RawMessage)
		if err := json.Unmarshal(keys["properties"], &propKeys); err != nil {
			return err
		}
		// "eta" exists only on the full type
		_, isFull = propKeys["eta"]
	}

	if isGeojson {
		if isFull {
			c.Type = responsetype.FullGeojson
			return json.Unmarshal(data, &c.CombinedFullGeojson)
		} else {
			c.Type = responsetype.SimpleGeojson
			return json.Unmarshal(data, &c.CombinedSimpleGeojson)
		}
	} else {
		if isFull {
			c.Type = responsetype.FullJson
			return json.Unmarshal(data, &c.CombinedFullJson)
		} else {
			c.Type = responsetype.SimpleJson
			return json.Unmarshal(data, &c.CombinedSimpleJson)
		}
	}
}

// IsZero is true iff the receiver is a default-valued CombinedMultiple struct.
func (c CombinedMultiple) IsZero() bool {
	return reflect.ValueOf(c).IsZero()
}

// AsSimpleJson returns the underlying CombinedSimpleJson response data if the response is of the correct type,
// and a zero (default-valued) CombinedSimpleJson struct otherwise.
func (c CombinedMultiple) AsSimpleJson() CombinedSimpleJson {
	return c.CombinedSimpleJson
}

// AsFullJson returns the underlying CombinedFullJson response data if the response is of the correct type,
// and a zero (default-valued) CombinedFullJson struct otherwise.
func (c CombinedMultiple) AsFullJson() CombinedFullJson {
	return c.CombinedFullJson
}

// AsSimpleGeojson returns the underlying CombinedSimpleGeojson response data if the response is of the correct type,
// and a zero (default-valued) CombinedSimpleGeojson struct otherwise.
func (c CombinedMultiple) AsSimpleGeojson() CombinedSimpleGeojson {
	return c.CombinedSimpleGeojson
}

// AsFullGeojson returns the underlying CombinedFullGeojson response data if the response is of the correct type,
// and a zero (default-valued) CombinedFullGeojson struct otherwise.
func (c CombinedMultiple) AsFullGeojson() CombinedFullGeojson {
	return c.CombinedFullGeojson
}

// CombinedFullJson is a response to Combined when requesting ModelType "Full" and ModelFormat "Json"
type CombinedFullJson struct {
	CourseOverGround         *float64  `json:"courseOverGround"`
	Latitude                 *float64  `json:"latitude"`
	Longitude                *float64  `json:"longitude"`
	Name                     string    `json:"name"`
	RateOfTurn               *float64  `json:"rateOfTurn"`
	ShipType                 *int      `json:"shipType"`
	SpeedOverGround          *float64  `json:"speedOverGround"`
	TrueHeading              *int      `json:"trueHeading"`
	Mmsi                     int       `json:"mmsi"`
	Msgtime                  time.Time `json:"msgtime"`
	Altitude                 *int      `json:"altitude"`
	NavigationalStatus       int       `json:"navigationalStatus"`
	ImoNumber                *int      `json:"imoNumber"`
	CallSign                 string    `json:"callSign"`
	Destination              string    `json:"destination"`
	Eta                      string    `json:"eta"`
	Draught                  *int      `json:"draught"`
	ShipLength               *int      `json:"shipLength"`
	ShipWidth                *int      `json:"shipWidth"`
	DimensionA               *int      `json:"dimensionA"`
	DimensionB               *int      `json:"dimensionB"`
	DimensionC               *int      `json:"dimensionC"`
	DimensionD               *int      `json:"dimensionD"`
	PositionFixingDeviceType int       `json:"positionFixingDeviceType"`
	ReportClass              string    `json:"reportClass"`
}

// IsZero is true iff the receiver is a default-valued CombinedFullJson struct.
func (a CombinedFullJson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// CombinedSimpleJson is a response to Combined when requesting ModelType "Simple" and ModelFormat "Json"
type CombinedSimpleJson struct {
	CourseOverGround *float64  `json:"courseOverGround"`
	Latitude         *float64  `json:"latitude"`
	Longitude        *float64  `json:"longitude"`
	Name             string    `json:"name"`
	RateOfTurn       *float64  `json:"rateOfTurn"`
	ShipType         *int      `json:"shipType"`
	SpeedOverGround  *float64  `json:"speedOverGround"`
	TrueHeading      *int      `json:"trueHeading"`
	Mmsi             int       `json:"mmsi"`
	Msgtime          time.Time `json:"msgtime"`
}

// IsZero is true iff the receiver is a default-valued CombinedSimpleJson struct.
func (a CombinedSimpleJson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// CombinedFullGeojson is a response to Combined when requesting ModelType "Full" and ModelFormat "Geojson"
type CombinedFullGeojson struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Mmsi                     int       `json:"mmsi"`
		Name                     string    `json:"name"`
		Msgtime                  time.Time `json:"msgtime"`
		SpeedOverGround          *float64  `json:"speedOverGround"`
		CourseOverGround         *float64  `json:"courseOverGround"`
		NavigationalStatus       int       `json:"navigationalStatus"`
		RateOfTurn               *float64  `json:"rateOfTurn"`
		ShipType                 *int      `json:"shipType"`
		TrueHeading              *int      `json:"trueHeading"`
		CallSign                 string    `json:"callSign"`
		Destination              string    `json:"destination"`
		Eta                      string    `json:"eta"`
		ImoNumber                *int      `json:"imoNumber"`
		DimensionA               *int      `json:"dimensionA"`
		DimensionB               *int      `json:"dimensionB"`
		DimensionC               *int      `json:"dimensionC"`
		DimensionD               *int      `json:"dimensionD"`
		Draught                  *int      `json:"draught"`
		ShipLength               *int      `json:"shipLength"`
		ShipWidth                *int      `json:"shipWidth"`
		PositionFixingDeviceType int       `json:"positionFixingDeviceType"`
		ReportClass              string    `json:"reportClass"`
	} `json:"properties"`
}

// IsZero is true iff the receiver is a default-valued CombinedFullJson struct.
func (a CombinedFullGeojson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// CombinedSimpleGeojson is a response to Combined when requesting ModelType "Simple" and ModelFormat "Geojson"
type CombinedSimpleGeojson struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Mmsi             int       `json:"mmsi"`
		Name             string    `json:"name"`
		Msgtime          time.Time `json:"msgtime"`
		SpeedOverGround  *float64  `json:"speedOverGround"`
		CourseOverGround *float64  `json:"courseOverGround"`
		RateOfTurn       *float64  `json:"rateOfTurn"`
		ShipType         *int      `json:"shipType"`
		TrueHeading      *int      `json:"trueHeading"`
	} `json:"properties"`
}

// IsZero is true iff the receiver is a default-valued CombinedSimpleGeojson struct.
func (a CombinedSimpleGeojson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}
