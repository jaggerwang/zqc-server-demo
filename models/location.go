package models

type Position []float32

type Point struct {
	Type        string   `bson:"type"`
	Coordinates Position `bson:"coordinates"`
}

func NewPoint(lon float32, lat float32) *Point {
	return &Point{"Point", Position{lon, lat}}
}
