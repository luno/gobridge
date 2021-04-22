package templates

import (
	"os"
	"text/template"
)

type HttpClient struct {
	Service string

	Type    string
	Method  string
	Params  []string
	Results []string

	RequestType string
	Request     map[string]string

	ResponseType   string
	ResponseParams []string

	Return []string
}

func (cl *HttpClient) AddTo(file *os.File) error {
	return template.Must(template.New("").Parse(httpClientTemplate)).Execute(file, cl)
}

var httpClientTemplate = `
func (c * Client) {{.Method}}(ctx context.Context{{ range $key, $value := .Params }}, {{ $value }}{{ end }}) ({{- range $key, $value := .Results }}{{if $key}}, {{end}}{{ $value }}{{- end }}) {
	req := {{.RequestType}} {
	{{- range $key, $value := .Request }}
		{{$key}}: {{$value}},
	{{- end}}
	}

	b, err := json.Marshal(req)
	if err != nil {
		return {{ range $key, $value := .Return }}{{if $key}}, {{end}}{{ $value }}{{ end }}
	}

	uniquePath := "/{{.Service}}/{{.Method}}"
	buf := bytes.NewBuffer(b)
	httpResp, err := ctxhttp.Post(ctx, c.HttpClient, c.Address + uniquePath, "application/json", buf)
	if err != nil {
		return {{ range $key, $value := .Return }}{{if $key}}, {{end}}{{ $value }}{{ end }}
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return {{ range $key, $value := .Return }}{{if $key}}, {{end}}{{ $value }}{{ end }}
	}

	var resp {{.ResponseType}}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return {{ range $key, $value := .Return }}{{if $key}}, {{end}}{{ $value }}{{ end }}
	}

	{{ if not (eq (len .ResponseParams) 0)}}return {{range $key, $value := .ResponseParams }}{{if $key}}, {{end}}{{if eq $value "_"}}{{else if $value}}resp.{{end}}{{$value}}{{end}}, nil{{end}}{{ if eq (len .ResponseParams) 0}}return nil{{end}}
}
`