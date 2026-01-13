package domain

import "github.com/google/uuid"

// Coordinates groups the spatial data
type Coordinates struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

type Address struct {
    ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
    RawAddress        string
    NormalizedAddress string      `gorm:"uniqueIndex"`
    Coordinates       Coordinates `gorm:"embedded"` // This groups them in Go
    Accuracy          string
    Source            string
    Geom              string      `gorm:"type:geography(Point,4326)"`
}