package creator

import (
	"sync"

	ad "github.com/jvillablanca/rolecheck/internal/infraestructura/adaptadores/depend/cfg"
	id "github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/depend"
)

var (
	Crea *Creator
	once sync.Once
	cfg  *ad.Cfg
)

type Creator struct {
}

func (c *Creator) IniCfg(username1, host1, username2, host2, userAdmin, passAdmin, nombreAmbiente1, nombreAmbiente2 string) id.Cfg {
	cfg = ad.NewCfg(username1, host1, username2, host2, userAdmin, passAdmin, nombreAmbiente1, nombreAmbiente2)
	return cfg
}

func (c *Creator) GetCfg() id.Cfg {
	if cfg == nil {
		panic("Configuration not initialized. Call IniCfg first.")
	}
	return cfg
}

func init() {
	once.Do(func() {
		Crea = &Creator{}
	})
}
