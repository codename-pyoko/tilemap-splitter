package split

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseJSON(r io.Reader) (Tilemap, error) {
	tilemap := Tilemap{}
	if err := json.NewDecoder(r).Decode(&tilemap); err != nil {
		return Tilemap{}, fmt.Errorf("failed to json decode: %w", err)
	}

	return tilemap, nil
}
