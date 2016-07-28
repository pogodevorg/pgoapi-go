package api

import (
	"github.com/pkmngo-odi/pogo-protos"
)

const feedBufferSize = 256

// Feed is an intermediary that holds a buffer of encountered game items and a reporter to notify about them
type Feed struct {
	repoter Reporter
	buffer  chan interface{}
}

// NewFeed construct a new Feed
func NewFeed(r Reporter) *Feed {
	return &Feed{
		buffer:  make(chan interface{}, feedBufferSize),
		repoter: r,
	}
}

// Push is used to get an entry on to the feed buffer
func (f *Feed) Push(entry interface{}) {
	select {
	default:
		return
	case f.buffer <- entry:
		return
	}
}

// Report is a blocking operation that starts retrieving data from the feed buffer
func (f *Feed) Report() {
	for {
		select {
		default:
			// NOOP: No entries in the buffer
		case entry := <-f.buffer:
			switch e := entry.(type) {
			default:
				// NOOP: Cannot report type
			case *protos.GetMapObjectsResponse:
				cells := e.GetMapCells()
				for _, cell := range cells {
					pokemons := cell.GetWildPokemons()
					if len(pokemons) > 0 {
						f.repoter.WildPokemons(pokemons)
					}
					forts := cell.GetForts()
					if len(forts) > 0 {
						f.repoter.Forts(forts)
					}
				}
				// Report Forts & Pokemons
			}
		}
	}
}
