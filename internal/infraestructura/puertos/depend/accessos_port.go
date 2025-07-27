package depend

import (
	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

type Accessos interface {
	GetAccessos(cuenta dominio.Cuenta) dominio.Accessos
}
