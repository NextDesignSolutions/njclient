package njclient

import (
	"fmt"
	"net/http"
)

type SysmonReadResp struct {
	Value uint32    `json:"value"`
	Error *NjtError `json:"error"`
}

type SysmonWriteResp struct {
	Error *NjtError `json:"error"`
}

type SysmonWriteBody struct {
	Value uint32 `json:"value"`
}

type Sysmon struct {
	client *Fpga
}

func NewSysmon(fpga *Fpga) (*Sysmon, error) {
	s := &Sysmon{
		client: fpga,
	}
	return s, nil
}

func (s *Sysmon) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	msg := fmt.Sprintf("/sysmon/%s", urlstr)
	return s.client.NewRequest(msg, method, body)
}

func (s *Sysmon) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return s.client.Do(req, into)
}

func (s *Sysmon) ReadRegister(address uint32) (uint32, error) {
	req, err := s.NewRequest(fmt.Sprintf("0x%x", address), "GET", nil)
	if err != nil {
		return 0, nil
	}
	var ret SysmonReadResp
	_, err = s.Do(req, &ret)
	if err != nil {
		return 0, nil
	}
	if ret.Error != nil {
		return 0, fmt.Errorf("Failed to read SYSMON register at addr 0x%08x: %s(%d)", address, ret.Error.Msg, ret.Error.Code)
	}
	return ret.Value, nil
}

func (s *Sysmon) WriteRegister(address uint32, value uint32) error {
	body := SysmonWriteBody{Value: value}
	req, err := s.NewRequest(fmt.Sprintf("0x%x", address), "PUT", body)
	if err != nil {
		return nil
	}
	var ret SysmonWriteResp
	_, err = s.Do(req, &ret)
	if err != nil {
		return nil
	}
	if ret.Error != nil {
		return fmt.Errorf("Failed to write SYSMON register at addr 0x%08x to value 0x%08x: %s(%d)", address, value, ret.Error.Msg, ret.Error.Code)
	}
	return nil

}
