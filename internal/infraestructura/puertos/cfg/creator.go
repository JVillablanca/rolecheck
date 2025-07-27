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
	println("Inicializando configuración con los siguientes parámetros:")
	println("username1:", username1)
	println("host1:", host1)
	println("username2:", username2)
	println("host2:", host2)
	println("userAdmin:", userAdmin)
	println("passAdmin:", passAdmin)
	println("nombreAmbiente1:", nombreAmbiente1)
	println("nombreAmbiente2:", nombreAmbiente2)
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
