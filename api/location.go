package api

import "github.com/golang/geo/s2"

const cellIDLevel = 15

// Location consists of coordinates in longitude, latitude and altitude
type Location struct {
	Lon float64
	Lat float64
	Alt float64
}

// GetCellIDs will return a slice of the closed neighbourhood cell ids for the current coordinates
func (l *Location) GetCellIDs() []uint64 {
	origin := s2.CellIDFromLatLng(s2.LatLngFromDegrees(l.Lat, l.Lon)).Parent(cellIDLevel)

	var cellIDs = make([]uint64, 0)
	cellIDs = append(cellIDs, uint64(origin))
	for _, cellID := range origin.EdgeNeighbors() {
		cellIDs = append(cellIDs, uint64(cellID))
	}

	return cellIDs
}
