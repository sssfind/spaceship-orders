package model

type Category int32

const (
	CategoryUnknown  Category = 0
	CategoryEngine   Category = 1
	CategoryFuel     Category = 2
	CategoryPorthole Category = 3
	CategoryWing     Category = 4
)

type Part struct {
	UUID          string
	Name          string
	Price         float64
	StockQuantity int64
	Description   string
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
}

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}
