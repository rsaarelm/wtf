package zendesk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Resource struct {
	Response interface{}
	Raw      string
}

func (widget *Widget) api(meth string, path string, params string) (*Resource, error) {
	trn := &http.Transport{}

	client := &http.Client{
		Transport: trn,
	}

	baseURL := fmt.Sprintf("https://%v.zendesk.com/api/v2", widget.settings.subdomain)
	URL := baseURL + "/tickets.json?sort_by=status"

	req, err := http.NewRequest(meth, URL, bytes.NewBufferString(params))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	apiUser := fmt.Sprintf("%v/token", widget.settings.username)
	req.SetBasicAuth(apiUser, widget.settings.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Resource{Response: &resp, Raw: string(data)}, nil
}
