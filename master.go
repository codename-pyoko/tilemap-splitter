package tmsplit

import (
	"fmt"
	"path"
	"strings"
)

type MasterTilemapEntry struct {
	Key           string `json:"key"`
	URL           string `json:"url"`
	TileX         int    `json:"tileX"`
	TileY         int    `json:"tileY"`
	WidthInTiles  int    `json:"widthInTiles"`
	HeightInTiles int    `json:"heightInTiles"`
}

type MasterTileset struct {
	SpritesheetKey string `json:"spritesheetKey"`
	SpritesheetURL string `json:"spritesheetUrl"`
	FrameWidth     int    `json:"frameWidth"`
	FrameHeight    int    `json:"frameHeight"`
	TilesetKey     string `json:"tilesetKey"`
}

type Spawn struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MasterFile struct {
	Spawn    Spawn                `json:"spawn"`
	Tilesets []MasterTileset      `json:"tilesets"`
	Tilemaps []MasterTilemapEntry `json:"tilemaps"`
}

func containsTileset(tilesets []MasterTileset, spritesheetKey string) bool {
	for _, ts := range tilesets {
		if ts.SpritesheetKey == spritesheetKey {
			return true
		}
	}
	return false
}

func CreateMasterFile(tilemaps []Tilemap, sourceFileBase string, nChunksWidth int) (MasterFile, error) {
	var mtilesets []MasterTileset
	var mtilemaps []MasterTilemapEntry
	var spawn Spawn
	for tmindex, tm := range tilemaps {
		for _, ts := range tm.Tilesets {
			spritesheetKey := fmt.Sprintf("spritesheet-%s", ts.Name)
			if containsTileset(mtilesets, spritesheetKey) {
				continue
			}

			mts := MasterTileset{
				SpritesheetKey: spritesheetKey,
				TilesetKey:     ts.Name,
				FrameWidth:     ts.TileWidth,
				FrameHeight:    ts.TileHeight,
				SpritesheetURL: path.Base(ts.Image),
			}

			mtilesets = append(mtilesets, mts)
		}

		noExt := strings.TrimSuffix(sourceFileBase, path.Ext(sourceFileBase))
		mtm := MasterTilemapEntry{
			Key:           fmt.Sprintf("%s-%d", noExt, tmindex),
			URL:           fmt.Sprintf("%s-%d.json", noExt, tmindex),
			HeightInTiles: tm.HeightInTiles,
			WidthInTiles:  tm.WidthInTiles,
			TileX:         tmindex % nChunksWidth * tm.WidthInTiles,
			TileY:         tmindex / nChunksWidth * tm.HeightInTiles,
		}

		for _, l := range tm.Layers {
			if l.Type != ObjectGroup {
				continue
			}

			for _, o := range l.Objects {
				if o.Type != "spawn" {
					continue
				}
				if spawn.X == 0 && spawn.Y == 0 || o.Properties.HasProperty("type", "primary") {
					spawn.X = mtm.TileX*tm.TileWidth + int(o.X)
					spawn.Y = mtm.TileY*tm.TileHeight + int(o.Y)
				}
			}
		}

		mtilemaps = append(mtilemaps, mtm)
	}

	return MasterFile{
		Spawn:    spawn,
		Tilesets: mtilesets,
		Tilemaps: mtilemaps,
	}, nil
}
