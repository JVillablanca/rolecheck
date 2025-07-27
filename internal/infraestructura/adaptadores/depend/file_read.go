package depend

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

func (a *Accessos) extrae_accessos(downloadDir string) (dominio.Accessos, error) {

	permisos, err := a.leePermisos(downloadDir)
	if err != nil {
		return dominio.Accessos{}, err
	}
	dataAccess, err := a.leeDataAccess(downloadDir)
	if err != nil {
		return dominio.Accessos{}, err
	}

	return dominio.Accessos{
		Permisos: permisos,
		Data:     dataAccess,
	}, nil
}

func (a *Accessos) leePermisos(downloadDir string) ([]dominio.Permiso, error) {

	// 1. Buscar el archivo *Hierarchical.zip
	var zipFile string
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "Hierarchical.zip") {
			zipFile = path
			return io.EOF // para detener el walk
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return nil, err
	}
	if zipFile == "" {
		return nil, errors.New("no se encontró archivo Hierarchical.zip")
	}

	// 2. Abrir el ZIP
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 3. Buscar el archivo *_Hierarchical.csv dentro del ZIP
	var csvFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "Hierarchical.csv") {
			csvFile = f
			break
		}
	}
	if csvFile == nil {
		return nil, errors.New("no se encontró archivo Hierarchical.csv en el ZIP")
	}

	// 4. Abrir el archivo CSV
	rc, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)
	reader.FieldsPerRecord = -1 // permite variable cantidad de campos

	// 5. Leer cabecera y encontrar índices
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	idxRol := -1
	idxRolHeredado := -1
	idxAplicacion := -1
	idxPermiso := -1
	for i, h := range header {
		switch strings.TrimSpace(h) {
		case "ROLE NAME":
			idxRol = i
		case "INHERITED ROLE NAME":
			idxRolHeredado = i
		case "APPLICATION NAME":
			idxAplicacion = i
		case "ENTITLEMENT":
			idxPermiso = i
		}
	}
	if idxRol == -1 || idxRolHeredado == -1 || idxAplicacion == -1 || idxPermiso == -1 {
		return nil, errors.New("no se encontraron todas las columnas requeridas en el CSV")
	}

	// 6. Leer línea por línea y acumular permisos
	var permisos []dominio.Permiso
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		perm := dominio.Permiso{
			Rol:         record[idxRol],
			RolHeredado: record[idxRolHeredado],
			Aplicacion:  record[idxAplicacion],
			Permiso:     record[idxPermiso],
		}
		permisos = append(permisos, perm)
	}

	return permisos, nil

}

func (a *Accessos) leeDataAccess(downloadDir string) ([]dominio.DataAccess, error) {
	// 1. Buscar el archivo *DataSec.zip
	var zipFile string
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "DataSec.zip") {
			zipFile = path
			return io.EOF // para detener el walk
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return nil, err
	}
	if zipFile == "" {
		return nil, errors.New("no se encontró archivo DataSec.zip")
	}

	// 2. Abrir el ZIP
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 3. Buscar el archivo *_DataSec.csv dentro del ZIP
	var csvFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "DataSec.csv") {
			csvFile = f
			break
		}
	}
	if csvFile == nil {
		return nil, errors.New("no se encontró archivo DataSec.csv en el ZIP")
	}

	// 4. Abrir el archivo CSV
	rc, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)
	reader.FieldsPerRecord = -1 // permite variable cantidad de campos

	// 5. Leer cabecera y encontrar índices
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	idxRol := -1
	idxRolHeredado := -1
	idxAplicacion := -1
	idxObjeto := -1
	for i, h := range header {
		switch strings.TrimSpace(h) {
		case "ROLE NAME":
			idxRol = i
		case "INHERITED ROLE NAME":
			idxRolHeredado = i
		case "APPLICATION NAME":
			idxAplicacion = i
		case "OBJECT NAME":
			idxObjeto = i
		}
	}
	if idxRol == -1 || idxRolHeredado == -1 || idxAplicacion == -1 || idxObjeto == -1 {
		return nil, errors.New("no se encontraron todas las columnas requeridas en el CSV")
	}
	// 6. Leer línea por línea y acumular data access
	var dataAccess []dominio.DataAccess
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data := dominio.DataAccess{
			Rol:         record[idxRol],
			RolHeredado: record[idxRolHeredado],
			Aplicacion:  record[idxAplicacion],
			Objeto:      record[idxObjeto],
		}
		dataAccess = append(dataAccess, data)
	}
	return dataAccess, nil
}
