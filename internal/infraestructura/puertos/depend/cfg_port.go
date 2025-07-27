package depend

import (
	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

type Cfg interface {
	GetCuenta1() dominio.Cuenta
	GetCuenta2() dominio.Cuenta
	GetUserAdmin() string
	GetPassAdmin() string
}
