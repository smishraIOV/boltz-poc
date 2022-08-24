package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RestClientError struct {
	Error string `json:"error"`
}

type RestClient struct {
	endpoint string
}

func NewRestClient(endpoint string) (*RestClient, error) {
	return &RestClient{
		endpoint: endpoint,
	}, nil
}

func (client *RestClient) Get(path string, reponseType any) error {

	resp, err := http.Get(client.endpoint + path)
	if err != nil {
		return fmt.Errorf("getpairs get %v: %w", client.endpoint+path, err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(reponseType)
	if err != nil {
		return fmt.Errorf("json decode (status: %v): %w", resp.Status, err)
	}

	return nil

}

func (client *RestClient) Post(path string, request any, response any) error {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(request)
	if err != nil {
		err = fmt.Errorf("json encode %v: %w", request, err)
		return err
	}

	resp, err := http.Post(client.endpoint+path, "application/json", buffer)
	if err != nil {
		err = fmt.Errorf("swapstatus post %v: %w", client.endpoint+path, err)
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		e := RestClientError{}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			err = fmt.Errorf("json decode (status: %v): %w", resp.Status, err)
			return err
		}
		badRequestError := BadRequestError(e.Error)
		err = fmt.Errorf("createswap result (status: %v) %w", resp.Status, &badRequestError)
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		err = fmt.Errorf("json decode (status ok): %w", err)
		return err
	}

	return nil
}
