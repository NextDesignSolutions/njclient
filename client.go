

type VersionInfo struct {
	Minor   int    `json:"minor"`
	Major   int    `json:"major"`
	Sha1    string `json:"sha1"`
	Version string `json:"version"`
}

type Config struct {
	APIVersion string
}

type Client struct {
	client  *http.Client
	BaseURL *url.URL
	UserAgent string
	Config *Config
	BoardService *BoardService
}

func NewClient(config *Config, string api_uri, client * http.Client) *Client {
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}
	baseURL, _ := url.Parse(api_uri)
	baseURL.Path = "api/" + config.APIVersion

	cl = client
	if cl == nil {
		cl = http.DefaultClient
	}

	c := &Client{
		client: cl,
		BaseURL: baseURL,
		userAgent: userAgent,
		Config: config,
	}

	c.initBoardService()
	return c;
}

func (c *Client) initBoardService() {
	bs := &BoardService {
			client: c
	}
	c.BoardService = bs;
}


func (c *Client) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.newEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", mediatype)
	req.header.Set("User-Agent", userAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, into interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(into)
	return resp, err
}

func (c *Client) GetServerVersion() (*VersionInfo, error) {
	v := new(VersionInfo)
	req, err := c.NewRequest("/version", "GET", nil)
	if err != nil {
		return nil, err
	}

	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}


