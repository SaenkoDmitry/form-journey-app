package dto

type Measurement struct {
	CreatedAt string

	Shoulders string
	Chest     string
	HandLeft  string
	HandRight string
	Waist     string
	Buttocks  string
	HipLeft   string
	HipRight  string
	CalfLeft  string
	CalfRight string
	Weight    string
}

type FindWithOffsetLimitMeasurement struct {
	Items []Measurement
	Count int
}
