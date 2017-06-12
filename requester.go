package requester

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type Requester struct {
	dsn         string
	parentValet string
	dpo         string
}

type Response struct {
	error  string `json:"error"`
	result string `json:"result"`
}

// create new request
func New(dsn, parentValet, dpo string) *Requester {
	return &Requester{
		dsn: dsn,
		parentValet: parentValet,
		dpo: dpo,
	}
}

// add new entry to emercoin dpo
func (r *Requester) Add(name, val string, days int) (*Response, error) {
	params := []string{
		r.prepareDpo(name),
		val,
	}
	return r.request("name_new", params)
}

// get entry from emercoin dpo
func (r *Requester) Get(name string) (*Response, error) {
	params := []string{
		r.prepareDpo(name),
	}
	return r.request("name_show", params)
}

// delete entry from emercoin dpo
func (r *Requester) Delete(name string) (*Response, error) {
	params := []string{
		r.prepareDpo(name),
	}
	return r.request("name_delete", params)
}

func (r *Requester) request(method string, params []string) (*Response, error) {
	type reqBody struct {
		Method string
		Params []string
	}
	rb := &reqBody{
		Method: method,
		Params: params,
	}
	jsonStr, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	jsonStr = []byte(jsonStr)
	httpReq, err := http.NewRequest("POST", r.dsn, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Add("Accept-Charset", "utf-8;q=0.7,*;q=0.7")
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *Requester) prepareDpo(name string) string {
	return fmt.Sprintf(r.dpo, name)
}