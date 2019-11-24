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

func (f *Fpga) getTemperature() (*MinValueMaxError, error) {
	var result MinValueMaxError
	req, err := f.NewRequest("temperature", "GET", nil)
	if err != nil {
		return nil, err
	}
	_, err = f.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *Fpga) getVccaux() (*MinValueMaxError, error) {
	var result MinValueMaxError
	req, err := f.NewRequest("vccaux", "GET", nil)
	if err != nil {
		return nil, err
	}
	_, err = f.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *Fpga) getVccbram() (*MinValueMaxError, error) {
	var result MinValueMaxError
	req, err := f.NewRequest("vccbram", "GET", nil)
	if err != nil {
		return nil, err
	}
	_, err = f.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *Fpga) getVccint() (*MinValueMaxError, error) {
	var result MinValueMaxError
	req, err := f.NewRequest("vccint", "GET", nil)
	if err != nil {
		return nil, err
	}
	_, err = f.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *Fpga) Temperature() (min float32, current float32, max float32, err error) {
	result, err := f.getTemperature()
	if err != nil {
		return 0, 0, 0, err
	}
	if result.Error != nil {
		return 0, 0, 0, fmt.Errorf("%s(%d)", result.Error.Msg, result.Error.Code)
	}
	return result.Min, result.Value, result.Max, nil
}

func (f *Fpga) Vccaux() (min float32, current float32, max float32, err error) {
	result, err := f.getVccaux()
	if err != nil {
		return 0, 0, 0, err
	}
	if result.Error != nil {
		return 0, 0, 0, fmt.Errorf("%s(%d)", result.Error.Msg, result.Error.Code)
	}
	return result.Min, result.Value, result.Max, nil
}

func (f *Fpga) Vccbram() (min float32, current float32, max float32, err error) {
	result, err := f.getVccbram()
	if err != nil {
		return 0, 0, 0, err
	}
	if result.Error != nil {
		return 0, 0, 0, fmt.Errorf("%s(%d)", result.Error.Msg, result.Error.Code)
	}
	return result.Min, result.Value, result.Max, nil
}

func (f *Fpga) Vccint() (min float32, current float32, max float32, err error) {
	result, err := f.getVccint()
	if err != nil {
		return 0, 0, 0, err
	}
	if result.Error != nil {
		return 0, 0, 0, fmt.Errorf("%s(%d)", result.Error.Msg, result.Error.Code)
	}
	return result.Min, result.Value, result.Max, nil
}
