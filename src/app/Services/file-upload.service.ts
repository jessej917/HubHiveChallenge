import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class FileUploadService {

  constructor(private httpClient: HttpClient) {}

  uploadFile(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('image', file, file.name);
    debugger;
    var thing = this.httpClient.post(environment.serverUrl + '/uploadImage', formData, {
      responseType: 'text',
    });
    return thing;
  }
}
