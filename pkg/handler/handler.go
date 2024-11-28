package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GoPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	State   struct {
		Zones []string `json:"zones"`
	} `json:"state"`
}

type Request struct {
	Go     GoPackage `json:"go"`
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
	Strict bool      `json:"strict"`
}

type File struct {
	Name string `json:"name"`
}

type ArtifactGo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ArtifactState struct {
	Status string `json:"status"`
}

type Artifact struct {
	Go    ArtifactGo    `json:"go"`
	State ArtifactState `json:"state"`
}

type Response struct {
	Artifacts []Artifact `json:"artifacts"`
}

func SendPackageRequest(packageName string) *http.Response {
	reqData := Request{
		Go: GoPackage{
			Name:    packageName,
			Version: "",
			State: struct {
				Zones []string `json:"zones"`
			}{Zones: []string{}},
		},
		Offset: 0,
		Limit:  50,
		Strict: false,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("https://repository.rt.ru/gateway/artifacts/findArtifacts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	return resp
}

func ParseResponse(responseData *http.Response) (*Response, error) {
	var resp Response
	err := json.NewDecoder(responseData.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
