package api

import (
	"github.com/OneOfOne/xxhash"
	"github.com/golang/protobuf/proto"
	protos "github.com/pogodevorg/POGOProtos-go"
)

const hashSeed = uint64(0x1B845238) // Static xxhash seed

func protoToXXHash64(seed uint64, pb proto.Message) (uint64, error) {
	h := xxhash.NewS64(seed)
	b, err := proto.Marshal(pb)
	if err != nil {
		return uint64(0), ErrFormatting
	}
	_, err = h.Write(b)
	if err != nil {
		return uint64(0), ErrFormatting
	}
	return h.Sum64(), nil
}

func protoToXXHash32(seed uint32, pb proto.Message) (uint32, error) {
	h := xxhash.NewS32(seed)
	b, err := proto.Marshal(pb)
	if err != nil {
		return uint32(0), ErrFormatting
	}
	_, err = h.Write(b)
	if err != nil {
		return uint32(0), ErrFormatting
	}
	return h.Sum32(), nil
}

func locationToXXHash32(seed uint32, location *Location) (uint32, error) {
	h := xxhash.NewS32(seed)
	b := location.GetBytes()
	_, err := h.Write(b)
	if err != nil {
		return uint32(0), ErrFormatting
	}
	return h.Sum32(), nil
}

func generateRequestHash(authTicket *protos.AuthTicket, request *protos.Request) (uint64, error) {
	h, err := protoToXXHash64(hashSeed, authTicket)
	if err != nil {
		return h, ErrFormatting
	}
	h, err = protoToXXHash64(h, request)
	if err != nil {
		return h, ErrFormatting
	}

	return h, nil
}

func generateLocation1(authTicket *protos.AuthTicket, location *Location) (uint32, error) {
	h, err := protoToXXHash32(uint32(hashSeed), authTicket)
	if err != nil {
		return h, ErrFormatting
	}
	h, err = locationToXXHash32(h, location)
	if err != nil {
		return h, ErrFormatting
	}
	return h, nil
}

func generateLocation2(location *Location) (uint32, error) {
	h, err := locationToXXHash32(uint32(hashSeed), location)
	if err != nil {
		return h, ErrFormatting
	}
	return h, nil
}
