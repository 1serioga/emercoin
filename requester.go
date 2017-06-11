package requester

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type requester struct {
	dsn         string
	parentValet string
	dpo         string
}

type response struct {
	error  string `json:"error"`
	result string `json:"result"`
}

// create new request
func NewRequester(dsn, parentValet, dpo string) *requester {
	return &requester{
		dsn: dsn,
		parentValet: parentValet,
		dpo: dpo,
	}
}

// add new entry to emercoin dpo
func (r *requester) Add(name, val string, days int) (*response, error) {
	params := []string{
		r.prepareDpo(name),
		val,
	}
	return r.request("name_new", params)
}

// get entry from emercoin dpo
func (r *requester) Get(name string) (*response, error) {
	params := []string{
		r.prepareDpo(name),
	}
	return r.request("name_show", params)
}

// delete entry from emercoin dpo
func (r *requester) Delete(name string) (*response, error) {
	params := []string{
		r.prepareDpo(name),
	}
	return r.request("name_delete", params)
}

func (r *requester) request(method string, params []string) (*response, error) {
	type reqBody struct {
		method string
		params []string
	}
	rb := &reqBody{
		method: method,
		params: params,
	}
	jsonStr, err := json.Marshal(rb)
	if err != nil {
		return nil, fmt.Errorf("could not encode params to json %v", jsonStr)
	}
	jsonStr = []byte(jsonStr)
	httpReq, err := http.NewRequest("POST", r.dsn, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("could not make a request %v", r.dsn)
	}
	httpReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Add("Accept-Charset", "utf-8;q=0.7,*;q=0.7")
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("could not make a request %v", r.dsn)
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body from response %v", httpResp)
	}

	var resp response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("could decode body %v", body)
	}

	return &resp, nil
}

func (r *requester) prepareDpo(name string) string {
	return fmt.Sprintf(r.dpo, name)
}