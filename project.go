package gproject

import (
	"encoding/json"
	"fmt"
	"gopkg.in/dougEfresh/toggl-http-client.v8"
)

// Toggl Project Definition
type Project struct {
	Id   uint64 `json:"id"`
	WId  uint64 `json:"wid"`
	CId  uint64 `json:"cid"`
	Name string `json:"name"`
}

type Projects []Project

const Endpoint = "/projects"

//Return a ProjectClient. An error is also returned when some configuration option is invalid
//    thc,err := gtoggl.NewClient("token")
//    pc,err := gproject.NewClient(tc)
func NewClient(thc *ghttp.TogglHttpClient) *ProjectClient {
	tc := &ProjectClient{
		thc: thc,
	}
	tc.endpoint = thc.Url + Endpoint
	return tc
}

type ProjectClient struct {
	thc      *ghttp.TogglHttpClient
	endpoint string
}

func (tc *ProjectClient) List() (Projects, error) {
	body, err := tc.thc.GetRequest(tc.endpoint)
	var projects []Project
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*body, &projects)
	return projects, err
}

func (tc *ProjectClient) Get(id uint64) (*Project, error) {
	return projectResponse(tc.thc.GetRequest(fmt.Sprintf("%s/%d", tc.endpoint, id)))
}

func (tc *ProjectClient) Create(p *Project) (*Project, error) {
	put := projectUpdateRequest{Project: p}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}
	return projectResponse(tc.thc.PostRequest(tc.endpoint, body))
}

func (tc *ProjectClient) Update(p *Project) (*Project, error) {
	put := projectUpdateRequest{Project: p}
	body, err := json.Marshal(put)
	if err != nil {
		return nil, err
	}

	return projectResponse(tc.thc.PutRequest(fmt.Sprintf("%s/%d", tc.endpoint, p.Id), body))
}

func (tc *ProjectClient) Delete(id uint64) error {
	_, err := tc.thc.DeleteRequest(fmt.Sprintf("%s/%d", tc.endpoint, id), nil)
	return err
}

func projectResponse(response *json.RawMessage, error error) (*Project, error) {
	if error != nil {
		return nil, error
	}
	var tResp ghttp.TogglResponse
	err := json.Unmarshal(*response, &tResp)
	if err != nil {
		return nil, err
	}
	var p Project
	err = json.Unmarshal(*tResp.Data, &p)
	if err != nil {
		return nil, err
	}
	return &p, err
}

//Configures a Client.
/*
    func SetURL(url string) ToggleClientOptionFunc {
	return func(c *TogglClient) error {
	    c.Url = url
	}
    }
*/
type ProjectClientOptionFunc func(*ProjectClient) error

type projectResponseTest struct {
	Data []byte `json:"data"`
}

type projectUpdateRequest struct {
	Project *Project `json:"project"`
}
