package split

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func fixTilesets(tm *Tilemap) {
	for tsindex := range tm.Tilesets {
		ts := &tm.Tilesets[tsindex]
		ts.Image = ts.ImageXML.Source
		ts.ImageHeight = ts.ImageXML.Height
		ts.ImageWidth = ts.ImageXML.Width
		ts.ImageXML = XMLImage{}

		// ts.Terrains = ts.XMLTerrains.Terrain
		// ts.XMLTerrains = TerrainTypes{}

		// for tindex, tile := range ts.Tiles {
		// ts.Tiles[tindex].Properties = tile.XMLProperties.Properties
		// ts.Tiles[tindex].XMLProperties = Properties{}

		// ts.Tiles[tindex].Animation = tile.XMLAnimation.Frames
		// ts.Tiles[tindex].XMLAnimation = Animations{}
		// }

	}
}

func fixLayers(tm *Tilemap) {
	for layerindex := range tm.Layers {
		layer := &tm.Layers[layerindex]
		layer.Compression = Compression(layer.XMLData.Compression)
		layer.Encoding = Encoding(layer.XMLData.Encoding)
		layer.Data = strings.TrimSpace(layer.XMLData.Data)
		layer.XMLData = XMLData{}

		layer.Type = TileLayer

		// if layer.DrawOrder == "" {
		// 	layer.DrawOrder = DrawOrderTopDown
		// }

		if layer.Opacity == 0 {
			// Really should differentiate between 0 and missing from XML
			layer.Opacity = 1
		}

		if !layer.Visible {
			// Really should differentiate between false and missing from XML
			// as it stands, a layer can never be invisible :(
			layer.Visible = true
		}
	}

	for _, layer := range tm.XMLObjectGroup {
		layer.Type = ObjectGroup

		if layer.DrawOrder == "" {
			layer.DrawOrder = DrawOrderTopDown
		}

		if layer.Opacity == 0 {
			// Really should differentiate between 0 and missing from XML
			layer.Opacity = 1
		}

		if !layer.Visible {
			// Really should differentiate between false and missing from XML
			// as it stands, a layer can never be invisible :(
			layer.Visible = true
		}

		tm.Layers = append(tm.Layers, layer)
	}
}

func ParseXML(r io.Reader) (Tilemap, error) {
	tilemap := Tilemap{}
	if err := xml.NewDecoder(r).Decode(&tilemap); err != nil {
		return Tilemap{}, fmt.Errorf("failed to xml decode: %w", err)
	}

	fixTilesets(&tilemap)
	fixLayers(&tilemap)

	return tilemap, nil
}

func parsePointList(l string) []Point {
	var points []Point
	for _, p := range strings.Split(l, " ") {
		xy := strings.Split(p, ",")
		x, err := strconv.ParseFloat(xy[0], 64)
		if err != nil {
			logrus.Warnf("failed to parse point '%s' as float, will skip: %v", p, err)
		}
		y, err := strconv.ParseFloat(xy[1], 64)
		if err != nil {
			logrus.Warnf("failed to parse point '%s' as float, will skip: %v", p, err)
		}
		points = append(points, Point{X: x, Y: y})

	}
	return points
}

func (o *Object) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	buf := struct {
		ID         int        `xml:"id,attr"`
		Name       string     `xml:"name,attr"`
		Type       string     `xml:"type,attr"`
		X          float64    `xml:"x,attr"`
		Y          float64    `xml:"y,attr"`
		Width      float64    `xml:"width,attr"`
		Height     float64    `xml:"height,attr"`
		Rotation   float64    `xml:"rotation,attr"`
		GID        int        `xml:"gid,attr"`
		Visible    bool       `xml:"visible,attr"`
		Template   string     `xml:"template,attr"`
		Properties []Property `xml:"properties>property"`
		Ellipse    *struct{}  `xml:"ellipse"`
		Point      *struct{}  `xml:"point"`
		Polygon    *struct {
			Points string `xml:"points,attr"`
		} `xml:"polygon"`
		Polyline *struct {
			Points string `xml:"points,attr"`
		} `xml:"polyline"`

		Text *struct{} `xml:"-"` // not yet supported
	}{}

	if err := d.DecodeElement(&buf, &start); err != nil {
		return err
	}

	o.ID = buf.ID
	o.Name = buf.Name
	o.Type = buf.Type
	o.X = buf.X
	o.Y = buf.Y
	o.Width = buf.Width
	o.Height = buf.Height
	o.Rotation = buf.Rotation
	o.GID = buf.GID
	o.Visible = buf.Visible
	o.Template = buf.Template
	o.Properties = buf.Properties

	if !o.Visible {
		// Unfortunate that we don't differentiate between missing value (where true is default) and actual `false`
		o.Visible = true
	}

	if buf.Ellipse != nil {
		o.Ellipse = true
	} else if buf.Point != nil {
		o.Point = true
	} else if buf.Polygon != nil {
		o.Polygon = parsePointList(buf.Polygon.Points)
	} else if buf.Polyline != nil {
		o.Polyline = parsePointList(buf.Polyline.Points)
	}

	return nil
}

func (t *TileTerrain) UnmarshalXMLAttr(attr xml.Attr) error {
	corners := strings.Split(attr.Value, ",")

	var n []int
	for _, c := range corners {
		v := -1
		if c != "" {
			i64, err := strconv.ParseInt(c, 10, 32)
			if err != nil {
				return err
			}
			v = int(i64)
		}
		n = append(n, v)
	}

	*t = n

	return nil
}

func (p *Property) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	buf := struct {
		Name  string       `xml:"name,attr"`
		Type  PropertyType `xml:"type,attr"`
		Value string       `xml:"value,attr"`
	}{}

	if err := d.DecodeElement(&buf, &start); err != nil {
		return err
	}

	p.Name = buf.Name
	p.Type = buf.Type

	var err error
	switch buf.Type {
	default:
		fallthrough
	case PropertyTypeString:
		p.Type = "string"
		p.Value = buf.Value

	case PropertyTypeColor:
		p.Value = buf.Value

	case PropertyTypeFile:
		p.Value = buf.Value

	case PropertyTypeBool:
		p.Value, err = strconv.ParseBool(buf.Value)

	case PropertyTypeFloat:
		p.Value, err = strconv.ParseFloat(buf.Value, 64)

	case PropertyTypeInt:
		p.Value, err = strconv.ParseInt(buf.Value, 10, 64)
	}

	if err != nil {
		return err
	}

	return nil
}
