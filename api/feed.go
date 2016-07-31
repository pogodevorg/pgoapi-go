package api

// Feed is a common interface to act on encountered
type Feed interface {
	// Push is used to put response messages on to the feed
	Push(entry interface{})
}

// VoidFeed is a feed that does nothing with the data
type VoidFeed struct {
}

// Push pushes the entry in to nothing
func (f *VoidFeed) Push(entry interface{}) {
	// NOOP
}
