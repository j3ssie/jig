package core

// Options global options
type Options struct {
	ConfigFile   string
	LogFile      string
	Timeout      int
	Concurrency  int
	Retry        int
	Redirect     bool
	Proxy        string
	Input        string
	Output       string
	OutputFolder string
	OutputType   string
	NoBanner     bool
	NoOutput     bool
	Verbose      bool
	Quite        bool
	Force        bool
	Helper       bool
	UseChrome    bool
	Debug        bool
	Params       []string
}

// Request all information about request
type Request struct {
	Timeout  int
	Repeat   int
	Scheme   string
	Host     string
	Port     string
	Path     string
	URL      string
	Proxy    string
	Method   string
	Redirect bool
	Headers  []map[string]string
	Body     string
	Beautify string
}

// Response all information about response
type Response struct {
	HasPopUp       bool
	StatusCode     int
	Status         string
	Headers        []map[string]string
	Body           string
	ResponseTime   float64
	Length         int
	Beautify       string
	Location       string
	Cookies        string
	BeautifyHeader string
}
