package api

import (
	"encoding/binary"
	"math"

	"github.com/golang/geo/s2"
	"github.com/pogodevorg/pogo-protos"
)

const cellIDLevel = 15
const earthRadiusInMeters = 6378100

// CellIDs is a slice of uint64s
type CellIDs []uint64

func (a CellIDs) Len() int           { return len(a) }
func (a CellIDs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CellIDs) Less(i, j int) bool { return a[i] < a[j] }

// Location consists of coordinates in longitude, latitude and altitude
type Location struct {
	Lon float64
	Lat float64
	Alt float64
}

// GetCellIDs will return a slice of the closed neighbourhood cell ids for the current coordinates
func (l *Location) GetCellIDs() CellIDs {
	origin := s2.CellIDFromLatLng(s2.LatLngFromDegrees(l.Lat, l.Lon)).Parent(cellIDLevel)

	var cellIDs = make([]uint64, 0)
	cellIDs = append(cellIDs, uint64(origin))
	for _, cellID := range origin.EdgeNeighbors() {
		cellIDs = append(cellIDs, uint64(cellID))
	}

	return cellIDs
}

// DistanceToFort returns distance between the location and a fort using the Haversine formula
// Reference: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
func (l *Location) DistanceToFort(fort *protos.FortData) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2 float64
	la1 = l.Lat * math.Pi / 180
	lo1 = l.Lon * math.Pi / 180
	la2 = fort.Latitude * math.Pi / 180
	lo2 = fort.Longitude * math.Pi / 180

	// calculate
	dla := math.Sin(0.5 * (la2 - la1))
	dlo := math.Sin(0.5 * (lo2 - lo1))
	h := dla*dla + math.Cos(la1)*math.Cos(la2)*dlo*dlo

	return 2 * earthRadiusInMeters * math.Asin(math.Sqrt(h))
}

// GetBytes returns a byte slice of the location coordinates
func (l *Location) GetBytes() []byte {
	b := make([]byte, 24)
	binary.BigEndian.PutUint64(b[0:8], math.Float64bits(l.Lat))
	binary.BigEndian.PutUint64(b[8:16], math.Float64bits(l.Lon))
	binary.BigEndian.PutUint64(b[16:24], math.Float64bits(l.Alt))
	return b[:24]
}
