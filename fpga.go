package njclient

import (
	"fmt"
	"net/http"
)

type FpgaInfo struct {
	Name string `json:"name"`
	Dna  string `json:"dna"`
}

type Fpga struct {
	client     *FpgaService
	Index      int
	Name       string
	Dna        string
	AxiService *AxiService
}

func (f *Fpga) Board() *Board {
	return f.client.Board()
}

func (f *Fpga) UpdatefromInfo(fi *FpgaInfo) error {
	f.Name = fi.Name
	f.Dna = fi.Dna
	return nil
}

func NewFpga(service *FpgaService) (*Fpga, error) {
	f := Fpga{
		client: service,
	}
	f.AxiService = NewAxiService(&f)
	return &f, nil
}

func (f *Fpga) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	msg := fmt.Sprintf("/%d/%s", f.Index, urlstr)
	return f.client.NewRequest(msg, method, body)
}

func (f *Fpga) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return f.client.Do(req, into)
}
