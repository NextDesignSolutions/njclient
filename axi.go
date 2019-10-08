package main

type AxiService struct {
	client *Fpga
	prefix string
}

type AxiCacheAttributes struct {
	Bufferable   bool
	Modifiable   bool
	Read_alloc   bool
	Write_alloc  bool
	Query_string string
}

func (ca *AxiCacheAttributes) generateQueryString() {
	s := "cache="
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

func (ato *AxiTransationOptions) generateQueryString() {
	s := ""
	if ato.Incr_mode == true {
		s += "+"
	}
	s += "/?count=" + str(ato.Count)

	ato.Query_string = s
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
}

func (t *AxiTransaction) buildQueryStr() {
	urlstr := "/" + str(t.Address) + t.Options.Query_string + t.Cache_attributes.Query_string
	t.Query_string = urlstr
}

func NewAxiTransaction(address uint64, rnw bool, options *AxiTransactionOptions, cache_attributes *AxiCacheAttributes, []data uint32) *AxiTransaction {
	a := AxiTransaction{
		Address:          address,
		Rnw:              rnw,
		Options:          options,
		Cache_attributes: cache_attributes,
	}
	if a.Options == nil {
		a.Options = default_transaction_options
	}
	if a.Cache_attributes == nil {
		a.Cache_attributes = default_cache_attributes
	}

	a.buildQueryStr()
}

type AxiResult struct {
	Value    []uint32
	Response int
}

var default_cache_attributes = NewAxiCacheAttributes(true, true, false, false)
var default_transaction_options = NewAxiTransactionOptions(false, 1)

func NewAxiService(client *Fpga) *AxiService {
	a := &AxiService{
		client: client,
		prefix: "axi32",
	}
	return a
}

func (a *AxiService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	return a.client.NewRequest(a.prefix+urlstr, method, body)
}

func (a *AxiService) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return a.client.Do(req, into)
}

func (a *AxiService) IssueTransaction(t *AxiTransaction) AxiResult, error {
	if t.rnw == true {
		req, err := a.NewRequest(t.Query_string, "GET", nil)
		if err != nil {
			return err
		}
		type ReadAxiResult struct {
			Value    []uint32  `json:"value"`
			Response string    `json:"response"`
			Error    *NjtError `json:"error"`
		}

		var result ReadAxiResult
		resp, err := a.Do(req, &result)
		if err != nil {
			return err
		}
		t.Data = &result.Value

	}

}
