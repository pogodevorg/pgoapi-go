package api

import (
	"sort"

	"github.com/golang/geo/s2"
)

type CellIDs []uint64

func (a CellIDs) Len() int           { return len(a) }
func (a CellIDs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CellIDs) Less(i, j int) bool { return a[i] < a[j] }

func CalculateCellIDs(loc *Location) CellIDs {
	latlng := s2.LatLngFromDegrees(loc.Lat, loc.Lon)
	origin := s2.CellIDFromLatLng(latlng).Parent(15)
	cellIDs := CellIDs{uint64(origin)}

	prev := origin.Prev()
	next := origin.Next()
	for i := 0; i < 10; i++ {
		cellIDs = append(cellIDs, uint64(prev))
		cellIDs = append(cellIDs, uint64(next))

		prev = prev.Prev()
		next = next.Next()
	}

	sort.Sort(cellIDs)
	return cellIDs
}
