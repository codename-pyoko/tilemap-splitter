// Package format formats a tilemap to an output format
package format

import (
	"fmt"
	"io"
	"text/template"

	"github.com/codename-pyoko/tilemap-splitter/split"
)

var TypescriptTemplate = `
{{- range $i, $e := .Tilemaps -}}
import t{{$i}} from '../../static/{{ $e.URL }}';
{{ end -}}

{{ range $i, $e := .Tilesets -}}
import s{{$i}} from '../../static/{{ $e.SpritesheetURL }}';
{{- end }}

const map = {
    tilemaps: [
        {{- range $i, $e := .Tilemaps }}
        {
            key: '{{ $e.Key }}',
            url: (t{{ $i }} as unknown) as string,
            tileX: {{ $e.TileX }},
            tileY: {{ $e.TileY }},
            widthInTiles: {{ $e.WidthInTiles }},
            heightInTiles: {{ $e.HeightInTiles }},
        },
        {{- end }}
    ],
    tilesets: [
        {{- range $i, $e := .Tilesets }}
        {
            spritesheetKey: '{{ $e.SpritesheetKey }}',
            spritesheetUrl: s{{ $i }},
            frameWidth: {{ $e.FrameWidth }},
            frameHeight: {{ $e.FrameHeight }},
            tilesetKey: '{{ $e.TilesetKey }}',
        },
        {{- end }}
    ],
};

export { map };
`

func FormatTypescript(w io.Writer, master split.MasterFile) error {
	templ, err := template.New("").Parse(TypescriptTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := templ.Execute(w, master); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}
