package alert

import "github.com/milosrs/go-hls-server/view"

type Alert struct {
	Color       view.Color
	Heading     string
	Description string
}
