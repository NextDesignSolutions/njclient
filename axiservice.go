package njclient

import (
	"fmt"
	"net/http"
)

type AxiService struct {
	client     *Fpga
	prefix     string
	AxiHandles []*AxiHandle
}

func NewAxiService(client *Fpga) *AxiService {
	a := &AxiService{
		client: client,
		prefix: "axi32",
	}

	b := NewAxiHandle(a, 0)
	if b == nil {
		return nil
	}

	a.AxiHandles = append(a.AxiHandles, b)
	return a
}

func (a *AxiService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	msg := fmt.Sprintf("%s/%s", a.prefix, urlstr)
	return a.client.NewRequest(msg, method, body)
}

func (a *AxiService) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return a.client.Do(req, into)
}

func (a *AxiService) GetAvailableAxiHandle() *AxiHandle {
	return a.AxiHandles[0]
}
