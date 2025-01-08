import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {environment} from '../environments/environment';

@Injectable()
export class HelloWorldService {

  constructor(private http: HttpClient) { }

  getTodos() {
    return this.http.get<any[]>('https://jsonplaceholder.typicode.com/todos');
  }

  getTitle() {
    return this.http.get<any[]>(environment.serverUrl + '/getUsers');
  }

  getPosts() {
    return this.http.get<any[]>(environment.serverUrl + '/getPosts');
  }

}