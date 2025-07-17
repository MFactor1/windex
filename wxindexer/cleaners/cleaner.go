package cleaners

import (
	"wxindexer/containers"
)

type Cleaner interface {
	Clean(text string) containers.Doc
}
