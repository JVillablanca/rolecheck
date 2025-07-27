package depend

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

var inProgress = map[string]struct{}{
	"WAIT": {}, "PENDING_VALIDATION": {}, "READY": {}, "RUNNING": {},
	"COMPLETED": {}, "ERROR_AUTO_RETRY": {}, "RETRYING": {}, "PAUSED": {},
	"BLOCKED": {}, "CANCELING": {}, "SCHEDULE_ENDED": {},
}

func isInProgress(status string) bool {
	s := strings.ToUpper(strings.TrimSpace(status))
	_, ok := inProgress[s]
	return ok
}

func (a *Accessos) waitForJobCompletion(cuenta dominio.Cuenta, jobId string, userAdmin string, passAdmin string) (string, error) {
	// Si el estado es "Running", esperar 5 segundos y volver a consultar
	for {
		status, err := a.getESSJobStatus(cuenta, jobId, userAdmin, passAdmin)
		if err != nil {
			return "", fmt.Errorf("error al obtener el estado del job: %w", err)
		}

		fmt.Println("El job ", jobId, " en ambiente ", cuenta.NombreAmbiente, " está en estado:", status)

		if !isInProgress(status) {
			return status, nil // ya dejó de trabajar
		}

		time.Sleep(10 * time.Second)
	}
}

func (a *Accessos) getESSJobStatus(cuenta dominio.Cuenta, jobId string, userAdmin string, passAdmin string) (string, error) {
	// Construye el XML para consultar el estado del job
	envelope := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:typ="http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/types/">
   <soapenv:Header/>
   <soapenv:Body>
      <typ:getESSJobStatus>
         <typ:requestId>%s</typ:requestId>
      </typ:getESSJobStatus>
   </soapenv:Body>
</soapenv:Envelope>`, jobId)

	// Prepara la petición HTTP
	action := "http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/getESSJobStatus"
	url := "https://" + cuenta.Host + ":443/fscmService/ErpIntegrationService"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(envelope))
	if err != nil {
		log.Fatalf("Error creando request: %v", err)
	}
	req.SetBasicAuth(userAdmin, passAdmin)
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", `"`+action+`"`)

	// 3) Manda la petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error llamando al servicio: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("Status %s: %s", resp.Status, b)
	}

	// Leer y parsear el XML de respuesta
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error leyendo respuesta: %v", err)
	}
	// Extraer solo el XML del cuerpo MIME
	start := bytes.Index(b, []byte("<env:Envelope"))
	end := bytes.LastIndex(b, []byte("</env:Envelope>"))
	if start == -1 || end == -1 {
		log.Fatalf("No se encontró el XML Envelope en la respuesta")
	}
	xmlBody := b[start : end+len("</env:Envelope>")]

	var envelopeResponse StatusEnvelope
	if err := xml.Unmarshal(xmlBody, &envelopeResponse); err != nil {
		log.Fatalf("Error parseando XML: %v", err)
	}

	return envelopeResponse.Body.GetESSJobStatusResponse.Result, nil
}

type StatusEnvelope struct {
	Body StatusBody `xml:"Body"`
}

type StatusBody struct {
	GetESSJobStatusResponse GetESSJobStatusResponse `xml:"getESSJobStatusResponse"`
}

type GetESSJobStatusResponse struct {
	Result string `xml:"result"`
}
