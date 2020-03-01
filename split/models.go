package split

type Compression string

const (
	NoCompression Compression = ""
	Zlib          Compression = "zlib"
	Gzip          Compression = "gzip"
)

type DrawOrder string

const (
	DrawOrderTopDown DrawOrder = "topdown"
	DrawOrderIndex   DrawOrder = "index"
)

type Encoding string

const (
	EncodingCSV    Encoding = "csv"
	EncodingBase64 Encoding = "base64"
)

type LayerType string

const (
	TileLayer   LayerType = "tilelayer"
	ObjectGroup LayerType = "objectgroup"
	ImageLayer  LayerType = "imagelayer"
	Group       LayerType = "group"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Object struct {
	Ellipse    bool       `json:"ellipse,omitempty"`
	GID        int        `json:"gid"`
	Height     float64    `json:"height"`
	ID         int        `json:"id" xml:"id,attr"`
	Name       string     `json:"name,omitempty" xml:"name,attr"`
	Point      bool       `json:"point,omitempty" xml:"point,attr"`
	Polygon    []Point    `json:"polygon,omitempty" xml:"polygon"`
	Polyline   []Point    `json:"polyline,omitempty" xml:"polyline"`
	Properties Properties `json:"properties,omitempty" xml:"properties>property"`
	Rotation   float64    `json:"rotation,omitempty"`
	Template   string     `json:"template,omitempty"`
	Text       string     `json:"text,omitempty"`
	Type       string     `json:"type,omitempty" xml:"type,attr"`
	Visible    bool       `json:"visible"`
	Width      float64    `json:"width"`
	X          float64    `json:"x"`
	Y          float64    `json:"y"`
}

type Orientation string

const (
	Orthogonal Orientation = "orthogonal"
	Isometrict Orientation = "isometric"
	Staggered  Orientation = "staggered"
	Hexagonal  Orientation = "hexagonal"
)

type RenderOrder string

const (
	RightDown RenderOrder = "right-down"
	RightUp   RenderOrder = "right-up"
	LeftDown  RenderOrder = "left-down"
	LeftUp    RenderOrder = "left-up"
)

type PropertyType string

const (
	PropertyTypeString PropertyType = "string"
	PropertyTypeInt    PropertyType = "int"
	PropertyTypeFloat  PropertyType = "float"
	PropertyTypeBool   PropertyType = "bool"
	PropertyTypeColor  PropertyType = "color"
	PropertyTypeFile   PropertyType = "file"
)

type Property struct {
	Name  string       `json:"name,omitempty" xml:"name,attr"`
	Type  PropertyType `json:"type,omitempty" xml:"type,attr"`
	Value interface{}  `json:"value,omitempty" xml:"value,attr"`
}

type Properties []Property

type PropertyValue struct{}

type Grid struct {
	Height      int         `json:"height,omitempty" xml:"height,attr"`
	Orientation Orientation `json:"orientation,omitempty" xml:"orientation,attr"`
	Width       int         `json:"width,omitempty" xml:"width,attr"`
}

type Terrain struct {
	Name       string     `json:"name" xml:"name,attr"`
	Properties Properties `json:"properties,omitempty" xml:"properties,attr"`
	Tile       int        `json:"tile" xml:"tile,attr"`
}

type TileOffset struct {
	X int `json:"x,omitempty" xml:"x,attr"`
	Y int `json:"y,omitempty" xml:"y,attr"`
}

type Frame struct {
	Duration int `json:"duration" xml:"duration,attr"`
	TileID   int `json:"tileid" xml:"tileid,attr"`
}

type TilesetType string

const (
	TilesetTypeTileset TilesetType = "tileset"
)

type WangSet struct {
	// TODO
}

type XMLImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type XMLData struct {
	Data        string `xml:",innerxml"`
	Encoding    string `xml:"encoding,attr"`
	Compression string `xml:"compression,attr"`
}

type XMLObjectGroup struct {
	ID        int       `xml:"id,attr"`
	Name      string    `xml:"name,attr"`
	Color     string    `xml:"color,attr"`
	Opacity   float64   `xml:"opacity,attr"`
	Visible   bool      `xml:"visible,attr"`
	OffsetX   float64   `xml:"offsetx,attr"`
	OffsetY   float64   `xml:"offsety,attr"`
	DrawOrder DrawOrder `xml:"draworder,attr"`

	Properties Properties `xml:"property"`
	Objects    []Object   `xml:"object"`
}

type Layer struct {
	Compression      Compression `json:"compression"`
	XMLData          XMLData     `json:"-" xml:"data"`
	Data             string      `json:"data,omitempty"`
	DrawOrder        DrawOrder   `json:"draworder,omitempty" xml:"draworder,attr"`
	Encoding         Encoding    `json:"encoding,omitempty"`
	HeightInTiles    int         `json:"height,omitempty" xml:"height,attr"`
	ID               int         `json:"id,omitempty" xml:"id,attr"`
	Image            string      `json:"image,omitempty"`
	Layers           []Layer     `json:"layers,omitempty"`
	Name             string      `json:"name,omitempty" xml:"name,attr"`
	Objects          []Object    `json:"objects" xml:"object"`
	OffsetX          float64     `json:"offsetx,omitempty" xml:"offsetx,attr"`
	OffsetY          float64     `json:"offsety,omitempty" xml:"offsety,attr"`
	Opacity          float64     `json:"opacity,omitempty" xml:"opacity,attr"`
	Properties       Properties  `json:"properties,omitempty" xml:"properties>property"`
	StartX           int         `json:"startx,omitempty"`
	StartY           int         `json:"starty,omitempty"`
	TransparentColor string      `json:"transparentcolor,omitempty"`
	Type             LayerType   `json:"type,omitempty"`
	Visible          bool        `json:"visible,omitempty" xml:"visible,attr"`
	WidthInTiles     int         `json:"width,omitempty" xml:"width,attr"`
	X                int         `json:"x" xml:"x,attr"`
	Y                int         `json:"y" xml:"y,attr"`
}

type TileTerrain []int

type Tile struct {
	Animation   []Frame     `json:"animation,omitempty" xml:"animation>frame"`
	ID          int         `json:"id" xml:"id,attr"`
	Image       string      `json:"image,omitempty"`
	ImageHeight int         `json:"imageheight,omitempty"`
	ImageWidth  int         `json:"imagewidth,omitempty"`
	ObjectGroup Layer       `json:"objectgroup,omitempty"`
	Probability float64     `json:"probability,omitempty"`
	Properties  Properties  `json:"properties,omitempty" xml:"properties>property"`
	Terrain     TileTerrain `json:"terrain,omitempty" xml:"terrain,attr"`
	Tile        int         `json:"tile,omitempty"`
}

type Tileset struct {
	BackgroundColor  string      `json:"backgroundcolor,omitempty"`
	Columns          int         `json:"columns,omitempty" xml:"columns,attr"`
	FirstGID         int         `json:"firstgid,omitempty" xml:"firstgid,attr"`
	Grid             Grid        `json:"grid,omitempty" xml:"grid"`
	ImageXML         XMLImage    `json:"-" xml:"image"`
	Image            string      `json:"image,omitempty"`
	ImageHeight      int         `json:"imageheight,omitempty"`
	ImageWidth       int         `json:"imagewidth,omitempty"`
	Margin           int         `json:"margin,omitempty" xml:"margin"`
	Name             string      `json:"name,omitempty" xml:"name,attr"`
	Properties       Properties  `json:"properties,omitempty" xml:"properties>property"`
	Source           string      `json:"source,omitempty" xml:"source,attr"`
	Spacing          int         `json:"spacing,omitempty" xml:"spacing,attr"`
	Terrains         []Terrain   `json:"terrains,omitempty" xml:"terraintypes>terrain"`
	TileCount        int         `json:"tilecount,omitempty" xml:"tilecount,attr"`
	TiledVersion     string      `json:"tiledversion,omitempty"`
	TileHeight       int         `json:"tileheight,omitempty" xml:"tileheight,attr"`
	TileOffset       TileOffset  `json:"tileoffset,omitempty" xml:"tileoffset"`
	Tiles            []Tile      `json:"tiles,omitempty" xml:"tile"`
	TileWidth        int         `json:"tilewidth,omitempty" xml:"tilewidth,attr"`
	TransparentColor string      `json:"transparentcolor,omitempty"`
	Type             TilesetType `json:"type,omitempty"`
	Version          float64     `json:"version,omitempty"`
	WangSets         []WangSet   `json:"wangsets,omitempty" xml:"wangsets>wangset"`
}

type Tilemap struct {
	HeightInTiles  int         `json:"height,omitempty" xml:"height,attr"`
	WidthInTiles   int         `json:"width,omitempty" xml:"width,attr"`
	TileHeight     int         `json:"tileheight,omitempty" xml:"tileheight,attr"`
	TileWidth      int         `json:"tilewidth,omitempty" xml:"tilewidth,attr"`
	Layers         []Layer     `json:"layers,omitempty" xml:"layer"`
	XMLObjectGroup []Layer     `json:"-" xml:"objectgroup"`
	Infinite       bool        `json:"infinite" xml:"infinite,attr"`
	NextLayerID    int         `json:"nextlayerid,omitempty" xml:"nextlayerid,attr"`
	NextObjectID   int         `json:"nextobjectid,omitempty" xml:"nextobjectid,attr"`
	Orientation    Orientation `json:"orientation,omitempty" xml:"orientation,attr"`
	Properties     Properties  `json:"properties,omitempty" xml:"properties>property"`
	RenderOrder    RenderOrder `json:"renderorder,omitempty" xml:"renderorder,attr"`
	StaggerAxis    string      `json:"staggeraxis,omitempty" xml:"staggeraxis,attr"`
	StaggerIndex   string      `json:"staggerindex,omitempty" xml:"staggerindex,attr"`
	TiledVersion   string      `json:"tiledversion,omitempty" xml:"tiledversion,attr"`
	Tilesets       []Tileset   `json:"tilesets,omitempty" xml:"tileset"`
	Version        float64     `json:"version,omitempty" xml:"version,attr"`
}

func (props Properties) HasProperty(name, value string) bool {
	for _, p := range props {
		if p.Name == name && p.Value == value {
			return true
		}
	}

	return false
}
