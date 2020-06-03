package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"frontendForPam/helper"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/create", create)
	http.ListenAndServe(":8080", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// json header default
	w.Header().Set("Content-Type", "application/json")

	// disable ssl validation
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// capture param
	fullName := r.FormValue("full_name")
	age := r.FormValue("age")

	// create new instances
	reader := strings.NewReader(fmt.Sprintf("{\"application\": "+
		"{\"com.edw.project01.User\": "+
		"{\"age\": %s,\"name\": \"%s\",\"failed\": null,\"success\": null}}}",
		age,
		fullName))
	request, _ := http.NewRequest("POST", "https://pam01:password@myapp-kieserver-pam.apps.edwin-cluster."+
		"sandbox1243.opentlc.com/services/rest/server/containers/Project01/processes/"+
		"Project01.Business01/instances", reader)

	// set header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "{\"message\":\"fail\"}")
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// get instance result
	request, _ = http.NewRequest("GET", fmt.Sprintf("https://pam01:password@myapp-kieserver-pam.apps."+
		"edwin-cluster.sandbox1243.opentlc.com/services/rest/server/queries/processes/instances/"+
		"%s/variables/instances/",
		string(body)), nil)
	client = &http.Client{}
	resp, err = client.Do(request)

	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "{\"message\":\"fail\"}")
		return
	}

	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Fprint(w, ParseXmlToGenerateResponse(body))
	return
}

// converting xml to response true / false status
func ParseXmlToGenerateResponse(xmlContent []byte) string {
	variableInstanceList := helper.VariableInstanceList{}
	xml.Unmarshal(xmlContent, &variableInstanceList)

	var success string
	var name string
	var imageUrl string

	for _, variableInstance := range variableInstanceList.VariableInstance {
		if variableInstance.Name == "status" {
			if variableInstance.Value == "true" {
				success = "success"
			} else {
				success = "failed"
			}
		}
		if variableInstance.Name == "nameResponse" {
			name = variableInstance.Value
		}
		if variableInstance.Name == "imageUrl" {
			imageUrl = variableInstance.Value
		}
	}

	return fmt.Sprintf("{\"message\":\"%s\", \"name\":\"%s\", \"imageUrl\":\"%s\"}", success, name, imageUrl)
}
