package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	"github.com/codename-pyoko/tilemap-splitter/format"
	"github.com/codename-pyoko/tilemap-splitter/split"
	"github.com/sirupsen/logrus"
)

var outputFmt string
var masterFile string
var pretty bool

func saveTilemaps(tilemaps []split.Tilemap) error {
	for index, tm := range tilemaps {
		f, err := os.Create(fmt.Sprintf(outputFmt, index))
		if err != nil {
			logrus.Fatalf("failed to create output file: %w", err)
		}
		defer f.Close()

		encoder := json.NewEncoder(f)
		if pretty {
			encoder.SetIndent("", "\t")
		}

		if err := encoder.Encode(&tm); err != nil {
			logrus.Fatalf("failed to encode tilemap: %v", err)
		}
		logrus.Infof("saved tilemap to %s", f.Name())
	}

	logrus.Infof("saved %d tilemaps", len(tilemaps))
	return nil
}

func saveMasterFile(master split.MasterFile) error {
	f, err := os.Create(masterFile)
	if err != nil {
		return fmt.Errorf("failed to open master file: %w", err)
	}

	defer f.Close()

	if err := format.FormatTypescript(f, master); err != nil {
		return fmt.Errorf("failed to format master file: %w", err)
	}

	logrus.Infof("saved masterfile to %s", f.Name())
	return nil
}

func main() {
	tiledJSON := flag.String("json", "", "Tiled JSON tilemap")
	tiledXML := flag.String("tmx", "", "Tiled TMX (xml) tilemap")

	flag.StringVar(&outputFmt, "out", "", "Output fmt string. %d for index")
	flag.BoolVar(&pretty, "pretty", false, "If output should be pretty printed")
	flag.StringVar(&masterFile, "master", "", "Master output file")

	var chunkHeight, chunkWidth int
	flag.IntVar(&chunkWidth, "chunkwidth", 100, "Width of each chunk")
	flag.IntVar(&chunkHeight, "chunkheight", 100, "Height of each chunk")

	flag.Parse()

	if *tiledJSON == "" && *tiledXML == "" {
		logrus.Errorf("Must specify tilemap source file")
		flag.Usage()
		return
	}

	if *tiledJSON != "" && *tiledXML != "" {
		logrus.Errorf("Cannot specify both json and tmx")
		flag.Usage()
		return
	}

	var sourceFile string
	var parser split.TilemapDecoder
	if *tiledJSON != "" {
		sourceFile = *tiledJSON
		parser = split.ParseJSON
	} else if *tiledXML != "" {
		sourceFile = *tiledXML
		parser = split.ParseXML
	}

	sourceNoExt := strings.TrimSuffix(sourceFile, path.Ext(sourceFile))

	if masterFile == "" {
		masterFile = fmt.Sprintf("%s-master.ts", sourceNoExt)
	}

	f, err := os.Open(sourceFile)
	if err != nil {
		logrus.Fatalf("failed to open source file: %v", err)
	}

	if outputFmt == "" {
		outputFmt = fmt.Sprintf("%s-%%d.json", sourceNoExt)
	}

	tilemap, err := parser(f)
	if err != nil {
		logrus.Fatalf("failed to parse tilemap: %v", err)
	}

	tilemaps, err := split.Run(tilemap, chunkWidth, chunkHeight)
	if err != nil {
		logrus.Fatalf("failed to split map: %v", err)
	}

	widthInTilemaps := int(math.Ceil(math.Max(float64(tilemap.WidthInTiles)/float64(chunkHeight), 1.0)))
	master, err := split.CreateMasterFile(tilemaps, path.Base(sourceFile), widthInTilemaps)
	if err != nil {
		logrus.Fatalf("failed to create master file: %v", err)
	}

	if err := saveMasterFile(master); err != nil {
		logrus.Fatalf("failed to save master file: %v", err)
	}

	if err := saveTilemaps(tilemaps); err != nil {
		logrus.Fatalf("failed to save tilemaps: %v", err)
	}

	logrus.Infof("tilemap successfully split into %d chunks", len(tilemaps))
}
