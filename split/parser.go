package split

import (
	"io"
)

type TilemapDecoder = func(io.Reader) (Tilemap, error)
