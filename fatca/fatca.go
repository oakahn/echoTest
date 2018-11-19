package fatca

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func Call(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func convertToXML(input Request) string {
	control := input.Control
	data := Control{
		Branch:        control.Branch,
		Channel:       control.Channel,
		RequesterName: control.RequesterName,
		RequestId:     control.RequestId,
		User:          control.User,
	}

	request := Request{
		Control:        data,
		CustomerID:     input.CustomerID,
		CustomerSource: input.CustomerSource,
	}

	getPartyFATCAInfo := GetPartyFATCAInfo{
		Request: request,
	}

	body := Body{
		GetPartyFATCAInfo: getPartyFATCAInfo,
	}

	soapenvEnvelope := EnvelopeRequest{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Ejbs:    "http://ejbs",
		Body:    body,
	}

	resp, err := xml.MarshalIndent(soapenvEnvelope, " ", "  ")

	if err != nil {
		log.Fatal(err)
	}

	return string(resp)
}

func getData(text string) EnvelopeResponse {

	// resp, _ := API.Post(Url.Facta, Data.MockData())
	url := "http://10.2.15.105:19080/FATCAHttpRouter/services/FATCA"
	resp, err := post(url, text)

	if err != nil {
		log.Fatal(err)
	}

	envelopeResponse := EnvelopeResponse{}

	jsonErr := xml.Unmarshal(resp, &envelopeResponse)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return envelopeResponse
}

func post(url string, payload string) ([]byte, error) {

	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("SOAPAction", "getPartyFATCAInfo")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}
