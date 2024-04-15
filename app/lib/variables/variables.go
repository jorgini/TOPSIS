package variables

import "errors"

type Method string
type Type string
type Variants int

const (
	TOPSIS      = "topsis"
	SMART       = "smart"
	Individuals = "individual"
	Group       = "group"
)

const (
	Benefit bool = true
	Cost    bool = false
)

const (
	None      = 0
	Triangle  = 1
	Trapezoid = 2
)

const (
	NormalizeWithSum           = 0b00000
	NormalizeWeightsByMidPoint = 0b00001
	NormalizeValueWithMax      = 0b00010
	Sengupta                   = 0b00011
	AlphaSlices                = 0b00100
	SqrtDistance               = 0b00101
	CbrtDistance               = 0b00110
	AggregateMatrix            = 0b00111
	AggregateFinals            = 0b01000
	Default                    = 0b01010
)

var (
	InvalidCaseOfOperation = errors.New("invalid case of operation")
	OutOfBounds            = errors.New("invalid index")
	InvalidSize            = errors.New("invalid size")
	EmptyValues            = errors.New("no values specified")
	IncompatibleTypes      = errors.New("incompatible types for operation")
	NoUsageMethod          = errors.New("this method shouldn't have usage")
)
