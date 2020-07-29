package main

// Struct used to marshal data for "/routes" response
type Node struct {
	Source string   `json:"source"`
	Edge []Edge `json:"routes"`
}
type Edge struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}


// Struct used to unmarshal json from OSMR
type ExtendedNode struct {
	Routes    []Routes    `json:"routes"`
	Waypoints []Waypoints `json:"waypoints"`
	Code      string      `json:"code"`
}
type Legs struct {
	Summary  string        `json:"summary"`
	Weight   int           `json:"weight"`
	Duration float64       `json:"duration"`
	Steps    []interface{} `json:"steps"`
	Distance float64       `json:"distance"`
}
type Routes struct {
	Legs       []Legs  `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     int     `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}
type Waypoints struct {
	Hint     string    `json:"hint"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}
