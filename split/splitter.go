// Package split can split tilemaps
package split

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

func decodeLayerData(b64enc string) ([]uint32, error) {
	b, err := base64.StdEncoding.DecodeString(b64enc)
	if err != nil {
		return nil, err
	}

	var ui32 uint32
	var decodedSize = binary.Size(&ui32)

	layerData := make([]uint32, len(b)/decodedSize)
	if err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, layerData); err != nil {
		return nil, err
	}

	return layerData, nil
}

func encodeLayerData(data []uint32) (string, error) {
	buf := bytes.Buffer{}
	if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func flattenMap(data [][]uint32, width int) []uint32 {
	ret := make([]uint32, len(data)*len(data[0]))
	for row, cols := range data {
		for col, gid := range cols {
			ret[row*width+col] = gid
		}
	}

	return ret
}

func countLayerType(tm Tilemap, layerType LayerType) int {
	c := 0
	for _, l := range tm.Layers {
		if l.Type == layerType {
			c++
		}
	}
	return c
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Run(tilemap Tilemap, chunkHeight, chunkWidth int) ([]Tilemap, error) {
	logrus.Debugf("tilemap widthInTiles: %d, heightInTiles: %d", tilemap.WidthInTiles, tilemap.HeightInTiles)
	widthInTilemaps := int(math.Ceil(math.Max(float64(tilemap.WidthInTiles)/float64(chunkHeight), 1.0)))
	heightInTilemaps := int(math.Ceil(math.Max(float64(tilemap.HeightInTiles)/float64(chunkWidth), 1.0)))
	logrus.Debugf("widthInTilemaps: %d, heightInTilemaps: %d", widthInTilemaps, heightInTilemaps)

	ntilemaps := int(math.Ceil(float64(widthInTilemaps) * float64(heightInTilemaps)))
	nlayers := countLayerType(tilemap, TileLayer)
	logrus.Debugf("creating %d tilemap(s) with %d layer(s) each", ntilemaps, nlayers)

	var decodedLayerData [][]uint32
	for _, layer := range tilemap.Layers {
		if layer.Type != TileLayer {
			continue
		}

		logrus.Debugf("adding layer: %v, %v", layer.Name, layer.Type)
		data, err := decodeLayerData(layer.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode layer data: %w", err)
		}

		logrus.Debugf("decoded %d gids", len(data))
		decodedLayerData = append(decodedLayerData, data)
	}

	var chunkedTilemaps []Tilemap

	for chunkIndex := 0; chunkIndex < ntilemaps; chunkIndex++ {

		buf := bytes.Buffer{}
		if err := json.NewEncoder(&buf).Encode(&tilemap); err != nil {
			return nil, fmt.Errorf("failed to encoder original tilemap to json: %w", err)
		}

		tm := Tilemap{}
		if err := json.NewDecoder(&buf).Decode(&tm); err != nil {
			return nil, fmt.Errorf("failed to decode original tilemap buffer from json: %w", err)
		}

		left := (chunkIndex % widthInTilemaps) * chunkWidth
		top := (chunkIndex / widthInTilemaps) * chunkHeight
		tm.WidthInTiles = min(chunkWidth, tilemap.WidthInTiles-left)
		tm.HeightInTiles = min(chunkHeight, tilemap.HeightInTiles-top)
		logrus.Debugf("tilemap %d: %d,%d (%dx%d)", chunkIndex, left, top, tm.WidthInTiles, tm.HeightInTiles)

		chunkoffset := tilemap.WidthInTiles*top + left

		for layerIndex, layerData := range decodedLayerData {

			var ll []uint32

			for itop := 0; itop < tm.HeightInTiles; itop++ {
				begin := chunkoffset + itop*tm.WidthInTiles
				end := begin + tm.WidthInTiles
				ll = append(ll, layerData[begin:end]...)
			}
			encoded, err := encodeLayerData(ll)
			if err != nil {
				return nil, fmt.Errorf("failed to encode layer data: %w", err)
			}
			tm.Layers[layerIndex].WidthInTiles = tm.WidthInTiles
			tm.Layers[layerIndex].HeightInTiles = tm.HeightInTiles
			tm.Layers[layerIndex].Data = encoded
		}

		for layerIndex, layer := range tm.Layers {
			if layer.Type != ObjectGroup {
				continue
			}

			objects := []Object{}
			for _, object := range layer.Objects {

				tileX := int(object.X / float64(tm.TileWidth))
				tileY := int(object.Y / float64(tm.TileHeight))
				if tileX >= left && tileX < left+chunkWidth && tileY >= top && tileY < top+chunkHeight {
					chunkX := math.Floor(float64(tileX) / float64(chunkWidth))
					chunkY := math.Floor(float64(tileY) / float64(chunkHeight))
					object.X -= chunkX * float64(chunkWidth*tm.TileWidth)
					object.Y -= chunkY * float64(chunkHeight*tm.TileHeight)
					objects = append(objects, object)
				}
			}
			tm.Layers[layerIndex].Objects = objects
		}
		chunkedTilemaps = append(chunkedTilemaps, tm)
	}

	return chunkedTilemaps, nil
}
