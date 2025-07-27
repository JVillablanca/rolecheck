package depend

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

func (a *Accessos) downloadFile(cuenta dominio.Cuenta, jobId string, userAdmin string, passAdmin string) (string, error) {
	// Construye el XML para descargar el archivo
	// 1) Construye el XML a mano
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:typ="http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/types/">
  <soapenv:Header/>
  <soapenv:Body>
    <typ:downloadESSJobExecutionDetails>
      <typ:requestId>%s</typ:requestId>
      <typ:fileType>All</typ:fileType>
    </typ:downloadESSJobExecutionDetails>
  </soapenv:Body>
</soapenv:Envelope>`, jobId)

	// 2) Prepara la petición HTTP
	action := "http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/downloadESSJobExecutionDetails"
	url := "https://" + cuenta.Host + ":443/fscmService/ErpIntegrationService"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(envelope))
	if err != nil {
		log.Fatalf("Error creando request: %v", err)
		return "", err
	}
	req.SetBasicAuth(userAdmin, passAdmin)
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", `"`+action+`"`)

	// 3) Manda la petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error llamando al servicio: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("Status %s: %s", resp.Status, b)
		return "", fmt.Errorf("status %s: %s", resp.Status, b)
	}

	// 4) Lee el multipart response (XOP) y extrae la parte ZIP
	//    Aquí simplificamos: volcamos TODO el cuerpo a un archivo .zip
	//    La mayoría de clientes SOAP pactan que la segunda parte es el attachment.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error leyendo respuesta: %v", err)
		return "", err
	}

	// 5) Guarda el ZIP
	// El zip se guarda en una carpeta con el nombre del ambiente+nombre de usuario+jobId y el nombre
	// del archivo se llamara output.zip

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error obteniendo ruta del ejecutable: %v", err)
		return "", err
	}
	exeDir := exePath
	if info, err := os.Stat(exePath); err == nil && !info.IsDir() {
		exeDir = filepath.Dir(exePath)
	}
	// Usar filepath.Join para construir la ruta correctamente en Windows

	outDir := fmt.Sprintf("downloads/%s_%s_%s", cuenta.NombreAmbiente, cuenta.Username, jobId)
	fullOutDir := filepath.Join(exeDir, outDir)

	if err := os.MkdirAll(fullOutDir, 0755); err != nil {
		log.Fatalf("Error creando directorio %s: %v", fullOutDir, err)
		return "", err
	}
	outFile := filepath.Join(fullOutDir, "output.zip")
	if err := os.WriteFile(outFile, bodyBytes, 0644); err != nil {
		log.Fatalf("Error escribiendo %s: %v", outFile, err)
		return "", err
	}
	log.Printf("✅ Archivo guardado como %s (%d bytes)", outFile, len(bodyBytes))

	return fullOutDir, nil
}
