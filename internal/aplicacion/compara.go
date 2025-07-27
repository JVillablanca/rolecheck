package aplicacion

import (
	"fmt"
	"os"
	"time"

	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
	cfg "github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/cfg"
	"github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/creator"
)

func ComparaCuentas() {
	c := cfg.Crea.GetCfg()
	p := creator.Crea

	accesos1 := p.GetRecuperaAccessos().GetAccessos(c.GetCuenta1())
	accesos2 := p.GetRecuperaAccessos().GetAccessos(c.GetCuenta2())

	now := time.Now()
	filename := fmt.Sprintf("comparacion_%s.txt", now.Format("20060102T1504"))
	f, err := os.Create(filename)
	if err != nil {
		println("Error creando archivo:", err.Error())
		return
	}
	defer f.Close()

	// Sección 1: Descripción de cuentas y fecha/hora
	fmt.Fprintf(f, "Comparación de cuentas\n")
	fmt.Fprintf(f, "Fecha y hora: %s\n", now.Format("2006-01-02 15:04"))
	fmt.Fprintf(f, "Cuenta 1: %s\n", c.GetCuenta1())
	fmt.Fprintf(f, "Cuenta 2: %s\n\n", c.GetCuenta2())

	// Sección 2: Roles únicos (Permiso.Rol)
	roles1 := make(map[string]bool)
	roles2 := make(map[string]bool)
	for _, p1 := range accesos1.Permisos {
		roles1[p1.Rol] = true
	}
	for _, p2 := range accesos2.Permisos {
		roles2[p2.Rol] = true
	}
	var rolesFaltantes1 []string
	var rolesFaltantes2 []string
	fmt.Fprintf(f, "Roles en Cuenta 1 (%s) y no en Cuenta 2 (%s):\n", c.GetCuenta1().NombreAmbiente, c.GetCuenta2().NombreAmbiente)
	for r := range roles1 {
		if !roles2[r] {
			fmt.Fprintf(f, "- %s\n", r)
			rolesFaltantes2 = append(rolesFaltantes2, r) // faltan en cuenta2
		}
	}
	fmt.Fprintf(f, "\nRoles en Cuenta 2 (%s) y no en Cuenta 1 (%s):\n", c.GetCuenta2().NombreAmbiente, c.GetCuenta1().NombreAmbiente)
	for r := range roles2 {
		if !roles1[r] {
			fmt.Fprintf(f, "- %s\n", r)
			rolesFaltantes1 = append(rolesFaltantes1, r) // faltan en cuenta1
		}
	}
	fmt.Fprintf(f, "\n")

	rolesFaltantes1Map := make(map[string]bool)
	rolesFaltantes2Map := make(map[string]bool)
	for _, r := range rolesFaltantes1 {
		rolesFaltantes1Map[r] = true
	}
	for _, r := range rolesFaltantes2 {
		rolesFaltantes2Map[r] = true
	}

	// Sección 3: Permisos únicos (todos los campos de Permiso)
	type PermisoKey struct {
		Rol         string
		RolHeredado string
		Aplicacion  string
		Permiso     string
	}
	perm1 := make(map[PermisoKey]dominio.Permiso)
	perm2 := make(map[PermisoKey]dominio.Permiso)
	for _, p := range accesos1.Permisos {
		if rolesFaltantes2Map[p.Rol] {
			continue // ignorar permisos cuyo rol falta en cuenta2
		}
		k := PermisoKey{p.Rol, p.RolHeredado, p.Aplicacion, p.Permiso}
		perm1[k] = p
	}
	for _, p := range accesos2.Permisos {
		if rolesFaltantes1Map[p.Rol] {
			continue // ignorar permisos cuyo rol falta en cuenta1
		}
		k := PermisoKey{p.Rol, p.RolHeredado, p.Aplicacion, p.Permiso}
		perm2[k] = p
	}
	fmt.Fprintf(f, "Permisos en Cuenta 1 (%s) y no en Cuenta 2 (%s):\n", c.GetCuenta1().NombreAmbiente, c.GetCuenta2().NombreAmbiente)
	for k, p := range perm1 {
		if _, ok := perm2[k]; !ok {
			fmt.Fprintf(f, "- %+v\n", p)
		}
	}
	fmt.Fprintf(f, "\nPermisos en Cuenta 2 (%s) y no en Cuenta 1 (%s):\n", c.GetCuenta2().NombreAmbiente, c.GetCuenta1().NombreAmbiente)
	for k, p := range perm2 {
		if _, ok := perm1[k]; !ok {
			fmt.Fprintf(f, "- %+v\n", p)
		}
	}
	fmt.Fprintf(f, "\n")

	// Sección 4: Data Access únicos (todos los campos de DataAccess)
	type DataKey struct {
		Rol         string
		RolHeredado string
		Aplicacion  string
		Objeto      string
	}
	data1 := make(map[DataKey]dominio.DataAccess)
	data2 := make(map[DataKey]dominio.DataAccess)
	for _, d := range accesos1.Data {
		if rolesFaltantes2Map[d.Rol] {
			continue // ignorar data access cuyo rol falta en cuenta2
		}
		k := DataKey{d.Rol, d.RolHeredado, d.Aplicacion, d.Objeto}
		data1[k] = d
	}
	for _, d := range accesos2.Data {
		if rolesFaltantes1Map[d.Rol] {
			continue // ignorar data access cuyo rol falta en cuenta1
		}
		k := DataKey{d.Rol, d.RolHeredado, d.Aplicacion, d.Objeto}
		data2[k] = d
	}
	fmt.Fprintf(f, "Data Access en Cuenta 1 (%s) y no en Cuenta 2 (%s):\n", c.GetCuenta1().NombreAmbiente, c.GetCuenta2().NombreAmbiente)
	for k, d := range data1 {
		if _, ok := data2[k]; !ok {
			fmt.Fprintf(f, "- %+v\n", d)
		}
	}
	fmt.Fprintf(f, "\nData Access en Cuenta 2 (%s) y no en Cuenta 1 (%s):\n", c.GetCuenta2().NombreAmbiente, c.GetCuenta1().NombreAmbiente)
	for k, d := range data2 {
		if _, ok := data1[k]; !ok {
			fmt.Fprintf(f, "- %+v\n", d)
		}
	}
	fmt.Fprintf(f, "\n")
}
