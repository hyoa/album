package mock

import "github.com/hyoa/album/api/internal/cdn"

type CDN struct{}

func (*CDN) SignGetUri(key string, size cdn.MediaSize, kind cdn.MediaKind) string {
	return "asigneduriformedia"
}
