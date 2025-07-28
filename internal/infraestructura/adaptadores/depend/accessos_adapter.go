package depend

import (
	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
	cfg "github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/cfg"
)

type Accessos struct {
}

func NewAccessos() *Accessos {
	return &Accessos{}
}

func (a *Accessos) GetAccessos(cuenta dominio.Cuenta) dominio.Accessos {

	// Se submite un job para obtener los accessos
	cfg := cfg.Crea.GetCfg()
	jobId := a.getJobId(cuenta, cfg.GetUserAdmin(), cfg.GetPassAdmin())
	if jobId == "" {
		panic("No se pudo obtener el jobId en el ambiente " + cuenta.NombreAmbiente)
	}
	// Se espera a que el job se complete
	status, err := a.waitForJobCompletion(cuenta, jobId, cfg.GetUserAdmin(), cfg.GetPassAdmin())
	if err != nil {
		panic("Error esperando la finalización del job " + jobId + ": " + err.Error())
	}
	if status != "SUCCEEDED" {
		panic("El job " + jobId + " no se completó exitosamente en el ambiente " + cuenta.NombreAmbiente + ". Estado: " + status)
	}

	// Se descargan los archivos generados por el job
	downloadDir, err := a.downloadFile(cuenta, jobId, cfg.GetUserAdmin(), cfg.GetPassAdmin())
	if err != nil {
		panic("Error descargando archivos del job " + jobId + ": " + err.Error())
	}

	// Se descomprime el ZIP descargado
	if err := a.descompactaZip(downloadDir); err != nil {
		panic("Error descomprimiendo el ZIP descargado: " + err.Error())
	}

	// Se extraen los accessos del ZIP descomprimido
	accessos, err := a.extrae_accessos(downloadDir)
	if err != nil {
		panic("Error extrayendo accessos del ZIP: " + err.Error())
	}

	return accessos
}
