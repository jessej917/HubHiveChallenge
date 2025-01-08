import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import {environment} from '../environments/environment';

@Injectable()
export class HelloWorldService {

  constructor(private http: HttpClient) { }

  getTodos() {
    return this.http.get<any[]>('https://jsonplaceholder.typicode.com/todos');
  }

  getFriends(user:any) {
    const params = new HttpParams()
      .set('username', user)
    return this.http.get<any[]>(environment.serverUrl + '/getFriends', { params });
  }

  getPosts() {
    return this.http.get<any[]>(environment.serverUrl + '/getPosts');
  }

}