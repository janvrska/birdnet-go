package observation

import (
	"fmt"
	"strings"
	"time"

	"github.com/tphakala/go-birdnet/pkg/config"
)

// Observation represents a single observation data point
type Note struct {
	ID             uint `gorm:"primaryKey"`
	Date           string
	Time           string
	ScientificName string
	CommonName     string
	Confidence     float64
	Latitude       float64
	Longitude      float64
	Threshold      float64
	Sensitivity    float64
	ClipName       string
}

// ParseSpeciesString extracts the scientific and common name from the species string.
func ParseSpeciesString(species string) (string, string) {
	parts := strings.SplitN(species, "_", 2) // Split into 2 parts at most: ScientificName and CommonName
	if len(parts) == 2 {
		return parts[0], parts[1] // Return ScientificName (parts[0]) and CommonName (parts[1])
	}
	// Log this to see what is being returned
	fmt.Println("Species string has an unexpected format:", species)
	return species, species
}

// New creates a new Observation.
func New(species string, confidence float64, predStart, predEnd float64) Note {
	scientificName, commonName := ParseSpeciesString(species)
	return Note{
		Date:           time.Now().Format("2006-01-02"), // Using the current date
		Time:           time.Duration(predEnd - predStart).String(),
		ScientificName: scientificName,
		CommonName:     commonName,
		Confidence:     confidence,
		//... fill other fields if needed
	}
}

// LogNote is the central function for logging observations.
func LogNote(cfg *config.Settings, note Note) error {
	// 1. Save to Database
	if err := SaveToDatabase(note); err != nil {
		return fmt.Errorf("failed to save note to database: %v", err)
	}

	use24HourFormat := true

	// 2. Save to Log File
	if err := LogNoteToFile(cfg.LogPath, note, use24HourFormat); err != nil {
		return fmt.Errorf("failed to log note to file: %v", err)
	}

	return nil
}