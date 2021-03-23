package templates

import (
	"os"
	"text/template"
)

type TSService struct {
	Name string
	Types   []TSType
	MethodNames []string
}

type TSType struct {
	Name   string
	Fields map[string]string
}

func (tss *TSService) AddTo(file *os.File) error {
	return template.Must(template.New("").Parse(tsServiceTemplate)).Execute(file, tss)
}

var tsServiceTemplate = `
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class {{.Name}}Service {
	const url = 'http://localhost:8080'

	constructor(private http: HttpClient) {}

{{- range $key, $value := .MethodNames }}
	public async {{$value}}(payload: Request{{$value}}): Promise<Response{{$value}}> {
		return await this.http.post(this.url + '/{{$value}}', JSON.stringify(payload)).toPromise() as Response{{$value}};
	}
{{- end }}
}
{{ range $key, $value := .Types }}
export interface {{$value.Name}} {
{{- range $key2, $value2 := $value.Fields }}
	{{ $key2 }}: {{ $value2 }};
{{- end }}
}
{{ end }}

`
