import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor(private http: HttpClient) {}

  // @ts-ignore
  public async HasPermission(payload: HasPermissionRequest): Promise<ResponseHasPermission> {
    // tslint:disable-next-line:max-line-length
    return await this.http.post(environment.BackendURL + '/usersservice/haspermission', JSON.stringify(payload)).toPromise() as HasPermissionResponse;
  }
}

export interface HasPermissionRequest {
  R: Role[];
}

export interface HasPermissionResponse {
  Bool: boolean;
}

export interface Toy {
  Design: string;
}

export interface User {
  ID: number;
  Name: string;
  Role: Role;
  T: Toy;
}

export enum Role {
  RoleAdmin = 2,
  RoleUnknown = 0,
  RoleUser = 1,
}
