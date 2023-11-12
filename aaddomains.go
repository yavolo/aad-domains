package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"bufio"
)

// Structs to represent the XML structure
type Envelope struct {
	Body Body `xml:"Body"`
}

type Body struct {
	Response GetFederationInformationResponse `xml:"GetFederationInformationResponseMessage"`
}

type GetFederationInformationResponse struct {
	Response Response `xml:"Response"`
}

type Response struct {
	Domains Domains `xml:"Domains"`
}

type Domains struct {
	DomainList []string `xml:"Domain"`
}

func main() {

	var domains []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}

	if len(domains) == 0 {
		fmt.Println("No domains provided.")
		return
	}

	var xmlPayloadBuffer bytes.Buffer

	for _, domain := range domains {
		url := "https://autodiscover-s.outlook.com/autodiscover/autodiscover.svc"
		xmlPayload := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
			<soap:Envelope xmlns:exm="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:ext="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
				<soap:Header>
					<a:Action soap:mustUnderstand="1">http://schemas.microsoft.com/exchange/2010/Autodiscover/Autodiscover/GetFederationInformation</a:Action>
					<a:To soap:mustUnderstand="1">https://autodiscover-s.outlook.com/autodiscover/autodiscover.svc</a:To>
					<a:ReplyTo>
						<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
					</a:ReplyTo>
				</soap:Header>
				<soap:Body>
					<GetFederationInformationRequestMessage xmlns="http://schemas.microsoft.com/exchange/2010/Autodiscover">
						<Request>
							<Domain>%s</Domain>
						</Request>
					</GetFederationInformationRequestMessage>
				</soap:Body>
			</soap:Envelope>`, domain)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(xmlPayload)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Type", "text/xml; charset=utf-8")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("SOAPAction", "http://schemas.microsoft.com/exchange/2010/Autodiscover/Autodiscover/GetFederationInformation")
		req.Header.Set("User-Agent", "AutodiscoverClient")
		req.Header.Set("Content-Length", fmt.Sprint(xmlPayloadBuffer.Len()))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Parse the XML response
		var envelope Envelope
		decoder := xml.NewDecoder(resp.Body)
		err = decoder.Decode(&envelope)
		if err != nil {
			fmt.Println("Error decoding XML response:", err)
			return
		}

		for _, domain := range envelope.Body.Response.Response.Domains.DomainList {
			fmt.Println(domain)
		}
	}
}
