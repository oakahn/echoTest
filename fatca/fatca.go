package fatca

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func Call(c echo.Context) error {
	req := new(Request)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error(), "statuscode": "400"})
		// return err
	}

	dataXML, err := convertToXML(req)

	out := c.Get("output").(io.Writer)
	c.Logger().SetOutput(out)
	c.Logger().Print("Hello", c.Response().Header().Get(echo.HeaderXRequestID))

	if err != nil {
		return err
	}

	resp, err := getData(dataXML)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp.Body.GetPartyFATCAInfoResponse)
}

func convertToXML(input *Request) (string, error) {
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
		return "", err
	}

	return string(resp), nil
}

func getData(text string) (EnvelopeResponse, error) {

	// resp, _ := API.Post(Url.Facta, Data.MockData())
	url := "http://10.2.15.105:19080/FATCAHttpRouter/services/FATCA"
	resp, err := post(url, text)

	if err != nil {
		return EnvelopeResponse{}, err
	}

	envelopeResponse := EnvelopeResponse{}

	err = xml.Unmarshal(resp, &envelopeResponse)

	if err != nil {
		return EnvelopeResponse{}, err
	}

	return envelopeResponse, nil
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
