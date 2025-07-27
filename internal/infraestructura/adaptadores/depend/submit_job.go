package depend

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jvillablanca/rolecheck/internal/aplicacion/dominio"
)

type Envelope struct {
	Body Body `xml:"Body"`
}

type Body struct {
	SubmitESSJobRequestResponse SubmitESSJobRequestResponse `xml:"submitESSJobRequestResponse"`
}

type SubmitESSJobRequestResponse struct {
	Result string `xml:"result"`
}

func (a *Accessos) getJobId(cuenta dominio.Cuenta, userAdmin string, passAdmin string) string {

	// 1) Construye el XML a mano
	envelope := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
xmlns:typ="http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/types/">
   <soapenv:Header/>
   <soapenv:Body>
      <typ:submitESSJobRequest>
         <typ:jobPackageName>/oracle/apps/ess/hcm/users</typ:jobPackageName>
         <typ:jobDefinitionName>RetrieveEntitlementsJob</typ:jobDefinitionName>
         <!--Zero or more repetitions:-->
         <typ:paramList>ORA_USER_NAME</typ:paramList>
         <typ:paramList>%s</typ:paramList>
         <typ:paramList></typ:paramList>
         <typ:paramList></typ:paramList>
         <typ:paramList></typ:paramList>
         <typ:paramList></typ:paramList>
         <typ:paramList></typ:paramList>
         <typ:paramList>true</typ:paramList>
         <typ:paramList>false</typ:paramList>
      </typ:submitESSJobRequest>
   </soapenv:Body>
</soapenv:Envelope>`, cuenta.Username)

	// 2) Prepara la petición HTTP
	action := "http://xmlns.oracle.com/apps/financials/commonModules/shared/model/erpIntegrationService/submitESSJobRequest"
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

	//fmt.Printf("Respuesta XML: %s\n", xmlBody)

	var envelopeResponse Envelope
	if err := xml.Unmarshal(xmlBody, &envelopeResponse); err != nil {
		log.Fatalf("Error parseando XML: %v", err)
	}

	//debug
	fmt.Printf("Respuesta del job: %s\n", envelopeResponse.Body.SubmitESSJobRequestResponse.Result)

	return envelopeResponse.Body.SubmitESSJobRequestResponse.Result
}
