package dominio

type Permiso struct {
	// Estructura para manejar roles y permisos
	Rol         string
	RolHeredado string
	Aplicacion  string
	Permiso     string
}

type DataAccess struct {
	Rol         string
	RolHeredado string
	Aplicacion  string
	Objeto      string
}

type Accessos struct {
	Permisos []Permiso
	Data     []DataAccess
}
