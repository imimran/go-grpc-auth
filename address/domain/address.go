package domain

import "github.com/google/uuid"

type Address struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	RawAddress        string
	NormalizedAddress string `gorm:"uniqueIndex"`
	Latitude          float64
	Longitude         float64
	Accuracy          string
	Source            string
	Geom              string `gorm:"type:geography(Point,4326)"`
}
