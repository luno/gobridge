
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class ExampleService {
	const url = 'http://localhost:8080'

	constructor(private http: HttpClient) {}
	public async Name(payload: RequestName): Promise<ResponseName> {
		return await this.http.post(this.url + '/Name', JSON.stringify(payload)).toPromise() as ResponseName;
	}
}

export interface RequestName {
}

export interface ResponseName {
	string: string;
}


