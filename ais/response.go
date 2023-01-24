package ais

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sync"
	"time"
)

type CancelFunc func()

type Response[T any] http.Response

func (r *Response[T]) Unmarshal() (T, error) {
	var obj T
	defer r.Body.Close()
	return obj, json.NewDecoder(r.Body).Decode(&obj)
}

type StreamResponse[T any] http.Response

func (r *StreamResponse[T]) UnmarshalStream(ctx context.Context) (<-chan T, <-chan error, CancelFunc, error) {
	scan := bufio.NewScanner(r.Body)
	if scan == nil {
		return nil, nil, nil, errors.New("received nil scanner")
	}

	out := make(chan T)
	done, cancel := makeSignalCh()

	errCh := make(chan error, 1)

	go func() {
		defer r.Body.Close()

		for {
			select {
			case <-done:
				close(out)
				return
			case <-ctx.Done():
				cancel()
			default:
				success := scan.Scan()
				if !success {
					errCh <- errors.New("failure on scan")
				}

				if err := scan.Err(); err != nil {
					errCh <- err
				}

				var res T
				if err := json.Unmarshal(scan.Bytes(), &res); err != nil {
					errCh <- err
					cancel()
				} else {
					out <- res
				}
			}
		}
	}()

	return out, errCh, cancel, nil
}

type SSEStreamResponse[T any] http.Response

func (r *SSEStreamResponse[T]) UnmarshalStream(ctx context.Context) (<-chan T, <-chan error, CancelFunc, error) {
	scan := bufio.NewScanner(r.Body)
	if scan == nil {
		return nil, nil, nil, errors.New("received nil scanner")
	}

	out := make(chan T)
	done, cancel := makeSignalCh()

	errCh := make(chan error, 1)

	go func() {
		defer r.Body.Close()

		for {
			select {
			case <-done:
				close(out)
				return
			case <-ctx.Done():
				cancel()
			default:
				success := scan.Scan()
				if !success {
					errCh <- errors.New("failure on scan")
				}

				if err := scan.Err(); err != nil {
					errCh <- err
				}

				var res T
				var err error
				res, err = unmarshalSSEData[T](scan.Bytes())
				if errors.Is(err, empty) || errors.Is(err, noMatch) {
					continue
				} else if err != nil {
					errCh <- err
					cancel()
				} else {
					out <- res
				}
			}
		}
	}()

	return out, errCh, cancel, nil
}

var (
	noMatch = errors.New("no match")
	empty   = errors.New("empty")
)

// unmarshalSSEData
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

type AisResponseType string

const (
	PositionType   AisResponseType = "Position"
	AtonType                       = "Aton"
	StaticdataType                 = "Staticdata"
)

// AisMultiple holds a union of the multiple response types that an AIS data request can return.
// Use the Type property to inspect which type the message is.
type AisMultiple struct {
	Type AisResponseType
	Position
	Aton
	Staticdata
}

func (a *AisMultiple) UnmarshalJSON(data []byte) error {
	typ := struct {
		Type AisResponseType `json:"type"`
	}{}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	a.Type = typ.Type

	switch a.Type {
	case PositionType:
		return json.Unmarshal(data, &a.Position)
	case AtonType:
		return json.Unmarshal(data, &a.Aton)
	case StaticdataType:
		return json.Unmarshal(data, &a.Staticdata)
	default:
		return fmt.Errorf("unknown type: %s", a.Type)
	}
}

func (a AisMultiple) IsPosition() bool {
	return a.Type == PositionType
}

func (a AisMultiple) IsAton() bool {
	return a.Type == AtonType
}

func (a AisMultiple) IsStaticdata() bool {
	return a.Type == StaticdataType
}

func (a AisMultiple) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

func (a AisMultiple) AsPosition() Position {
	return a.Position
}

func (a AisMultiple) AsAton() Aton {
	return a.Aton
}

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

func (a Staticdata) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

type CombinedMultiple struct {
	Type string `json:"type"`
	CombinedSimpleJson
	CombinedFullJson
	CombinedSimpleGeojson
	CombinedFullGeojson
}

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
			c.Type = "FullGeojson"
			return json.Unmarshal(data, &c.CombinedFullGeojson)
		} else {
			c.Type = "SimpleGeojson"
			return json.Unmarshal(data, &c.CombinedSimpleGeojson)
		}
	} else {
		if isFull {
			c.Type = "FullJson"
			return json.Unmarshal(data, &c.CombinedFullJson)
		} else {
			c.Type = "SimpleJson"
			return json.Unmarshal(data, &c.CombinedSimpleJson)
		}
	}
}

func (c CombinedMultiple) IsSimpleJson() bool {
	return c.Type == "SimpleJson"
}

func (c CombinedMultiple) IsFullJson() bool {
	return c.Type == "FullJson"
}

func (c CombinedMultiple) IsFullGeojson() bool {
	return c.Type == "FullGeojson"
}

func (c CombinedMultiple) IsSimpleGeojson() bool {
	return c.Type == "SimpleGeojson"
}

func (c CombinedMultiple) IsZero() bool {
	return reflect.ValueOf(c).IsZero()
}

func (c CombinedMultiple) AsSimpleJson() CombinedSimpleJson {
	return c.CombinedSimpleJson
}

func (c CombinedMultiple) AsFullJson() CombinedFullJson {
	return c.CombinedFullJson
}

func (c CombinedMultiple) AsSimpleGeojson() CombinedSimpleGeojson {
	return c.CombinedSimpleGeojson
}

func (c CombinedMultiple) AsFullGeojson() CombinedFullGeojson {
	return c.CombinedFullGeojson
}

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

func (a CombinedFullJson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

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

func (a CombinedSimpleJson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

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

func (a CombinedFullGeojson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

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

func (a CombinedSimpleGeojson) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

type Combined struct {
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

func (a Combined) IsZero() bool {
	return reflect.ValueOf(a).IsZero()
}

// makeSignalCh creates and returns a new buffered channel of type struct{} and a cancel function.
//
// The channel is used as a signaling mechanism, where sending a value to the channel signals an event has occurred.
// Receiving from the channel is used to wait for the event.
//
// The cancel function, cancelFunc, is used to close the channel, indicating that no more events will be sent.
// The function uses a mutex to make sure it is only closed once, even if called concurrently.
// If the channel has not been closed yet, it closes the channel and sets closed to true.
func makeSignalCh() (<-chan struct{}, CancelFunc) {
	mu := sync.Mutex{}
	closed := false
	ch := make(chan struct{})
	cancelFunc := func() {
		mu.Lock()
		defer mu.Unlock()
		if !closed {
			closed = true
			close(ch)
		}
	}

	return ch, cancelFunc
}
