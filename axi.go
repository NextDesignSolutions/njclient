package njclient

import (
	"errors"
	"fmt"
	"net/http"
)

type AxiHandle struct {
	client *AxiService
	index  int
}

func NewAxiHandle(client *AxiService, index int) *AxiHandle {
	a := &AxiHandle{
		client: client,
		index:  index,
	}
	return a
}

type AxiCacheAttributes struct {
	Bufferable   bool
	Modifiable   bool
	Read_alloc   bool
	Write_alloc  bool
	Query_string string
}

type AxiWriteBody struct {
	Value []uint32 `json:"value"`
}

func (ca *AxiCacheAttributes) generateQueryString() {
	s := "&cache="
	if ca.Bufferable {
		s += "b"
	}
	if ca.Modifiable {
		s += "m"
	}
	if ca.Read_alloc {
		s += "r"
	}
	if ca.Write_alloc {
		s += "w"
	}
	ca.Query_string = s
}

func NewAxiCacheAttributes(bufferable bool, modifiable bool, read_alloc bool, write_alloc bool) *AxiCacheAttributes {
	b := AxiCacheAttributes{
		Bufferable:  bufferable,
		Modifiable:  modifiable,
		Read_alloc:  read_alloc,
		Write_alloc: write_alloc,
	}
	b.generateQueryString()
	return &b
}

type AxiTransactionOptions struct {
	Incr_mode    bool
	Count        int
	Query_string string
}

func (ato *AxiTransactionOptions) generateQueryString() {
	s := ""
	if ato.Incr_mode == true {
		s += "+"
	}
	fmt.Printf("ato.Count = %d", ato.Count)
	cnt := fmt.Sprintf("?count=%d", ato.Count)
	ato.Query_string = fmt.Sprintf("%s%s", s, cnt)
}

func NewAxiTransactionOptions(incr_mode bool, count int) *AxiTransactionOptions {
	b := AxiTransactionOptions{
		Incr_mode: incr_mode,
		Count:     count,
	}
	b.generateQueryString()
	return &b
}

type AxiTransaction struct {
	Address          uint64
	Rnw              bool
	Options          *AxiTransactionOptions
	Cache_attributes *AxiCacheAttributes
	Query_string     string
	Data             *[]uint32
}

func (t *AxiTransaction) buildQueryStr() {
	urlstr := fmt.Sprintf("0x%08x%s%s", t.Address, t.Options.Query_string, t.Cache_attributes.Query_string)
	t.Query_string = urlstr
}

func NewAxiTransaction(address uint64, rnw bool, options *AxiTransactionOptions, cache_attributes *AxiCacheAttributes, data *[]uint32) *AxiTransaction {
	a := AxiTransaction{
		Address:          address,
		Rnw:              rnw,
		Options:          options,
		Cache_attributes: cache_attributes,
		Data:             data,
	}
	if a.Options == nil {
		a.Options = default_transaction_options
	}
	if a.Cache_attributes == nil {
		a.Cache_attributes = default_cache_attributes
	}

	a.buildQueryStr()
	return &a
}

type AxiResult struct {
	Value    *[]uint32
	Response string
}

var default_cache_attributes = NewAxiCacheAttributes(true, true, false, false)
var default_transaction_options = NewAxiTransactionOptions(false, 1)

func (a *AxiHandle) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	msg := fmt.Sprintf("%d/%s", a.index, urlstr)
	return a.client.NewRequest(msg, method, body)
}

func (a *AxiHandle) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return a.client.Do(req, into)
}

func (a *AxiHandle) IssueTransaction(t *AxiTransaction) (*AxiResult, error) {
	if t.Rnw == true {
		req, err := a.NewRequest(t.Query_string, "GET", nil)
		if err != nil {
			return nil, err
		}
		type ReadAxiResult struct {
			Value    []uint32  `json:"value"`
			Response string    `json:"response"`
			Error    *NjtError `json:"error"`
		}

		var result ReadAxiResult
		_, err = a.Do(req, &result)
		if err != nil {
			return nil, err
		}
		if result.Error != nil {
			return nil, errors.New(result.Error.Msg)
		}

		final_result := AxiResult{
			Value:    &result.Value,
			Response: result.Response,
		}
		return &final_result, nil

	} else {

		var b AxiWriteBody
		b.Value = *t.Data
		req, err := a.NewRequest(t.Query_string, "PUT", b)

		if err != nil {
			return nil, err
		}

		type WriteAxiResult struct {
			Response string    `json:"response"`
			Error    *NjtError `json:"error"`
		}

		var result WriteAxiResult
		_, err = a.Do(req, &result)
		if err != nil {
			return nil, err
		}
		if result.Error != nil {
			return nil, errors.New(result.Error.Msg)
		}

		var final_result AxiResult
		final_result.Value = nil
		final_result.Response = result.Response
		return &final_result, nil
	}
	return nil, errors.New("only read_axi or write_axi supported")
}
