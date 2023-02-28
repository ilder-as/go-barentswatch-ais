package responsetype

type Ais string

const (
	Position   Ais = "Position"
	Aton           = "Aton"
	Staticdata     = "Staticdata"
)

type Combined string

const (
	FullJson      Combined = "FullJson"
	SimpleJson             = "SimpleJson"
	FullGeojson            = "FullGeojson"
	SimpleGeojson          = "SimpleGeojson"
)
