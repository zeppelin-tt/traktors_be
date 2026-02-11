package main

// EngineType тип двигателя.
type EngineType string

const (
	Diesel   EngineType = "diesel"
	Gasoline EngineType = "gasoline"
	Gas      EngineType = "gas"
)

// Tractor модель трактора.
type Tractor struct {
	ID          string      `bson:"_id"                   json:"id"`
	Name        string      `bson:"name"                  json:"name"`
	Images      []string    `bson:"images"                json:"images"`
	Brand       *string     `bson:"brand,omitempty"       json:"brand"`
	Color       *string     `bson:"color,omitempty"       json:"color"`
	EngineType  *EngineType `bson:"engineType,omitempty"  json:"engineType"`
	Horsepower  *float64    `bson:"horsepower,omitempty"  json:"horsepower"`
	Year        *int        `bson:"year,omitempty"        json:"year"`
	Mileage     *float64    `bson:"mileage,omitempty"     json:"mileage"`
	VIN         *string     `bson:"vin,omitempty"         json:"vin"`
	PTS         *string     `bson:"pts,omitempty"         json:"pts"`
	PTSOwners   *int        `bson:"ptsOwners,omitempty"   json:"ptsOwners"`
	Location    *string     `bson:"location,omitempty"    json:"location"`
	Phone       *string     `bson:"phone,omitempty"       json:"phone"`
	Description *string     `bson:"description,omitempty" json:"description"`
	Price       *float64    `bson:"price,omitempty"       json:"price"`
}
