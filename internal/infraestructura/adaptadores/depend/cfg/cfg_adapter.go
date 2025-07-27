package cfg

import (
	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

type Cfg struct {
	cuenta1   dominio.Cuenta
	cuenta2   dominio.Cuenta
	userAdmin string
	passAdmin string
}

func NewCfg(username1, host1, username2, host2, userAdmin, passAdmin, nombreAmbiente1, nombreAmbiente2 string) *Cfg {
	return &Cfg{
		cuenta1: dominio.Cuenta{
			Username:       username1,
			Host:           host1,
			NombreAmbiente: nombreAmbiente1,
		},
		cuenta2: dominio.Cuenta{
			Username:       username2,
			Host:           host2,
			NombreAmbiente: nombreAmbiente2,
		},
		userAdmin: userAdmin,
		passAdmin: passAdmin,
	}
}

func (c *Cfg) GetCuenta1() dominio.Cuenta {
	return c.cuenta1
}

func (c *Cfg) GetCuenta2() dominio.Cuenta {
	return c.cuenta2
}

func (c *Cfg) GetUserAdmin() string {
	return c.userAdmin
}

func (c *Cfg) GetPassAdmin() string {
	return c.passAdmin
}
