package creator

import (
	"sync"

	ad "github.com/jvillablanca/rolecheck/internal/infraestructura/adaptadores/depend"
	pd "github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/depend"
)

var (
	Crea *Creator
	once sync.Once
)

type Creator struct {
}

func NewCreator() *Creator {
	return &Creator{}
}

func (c *Creator) GetRecuperaAccessos() pd.Accessos {
	return ad.NewAccessos()
}

func init() {
	once.Do(func() {
		Crea = &Creator{}
	})
}
