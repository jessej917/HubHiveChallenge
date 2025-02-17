import { NgIf } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../Services/auth.service';
import { environment } from '../../../environments/environment';
import { FileUploadService } from '../../Services/file-upload.service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-upload-post',
  standalone: true,
  imports: [NgIf, FormsModule],
  templateUrl: './upload-post.component.html',
  styleUrl: './upload-post.component.css'
})
export class UploadPostComponent {
  isOpen = false;
  file: File | null = null; // Variable to store file
  url: string | ArrayBuffer | null = "";

  postObj: Post;
  constructor(private http: HttpClient, private router: Router, private AS: AuthService, private FUS: FileUploadService) {
    this.postObj = new Post();
  }

  // On file Select
  onChange(event: any) {
    this.file = event.target.files[0];

    if (this.file) {
      const mimeType = this.file.type;
      if (mimeType.match(/image\/*/) == null) {
        //"Only images are supported."
        return;
      }

      const reader = new FileReader();
      reader.readAsDataURL(this.file);
      reader.onload = (_event) => {
        this.url = reader.result;
      }
    }
  }

  onUpload() {
    debugger;
    this.postObj.Username = this.AS.getUser();
    if (this.file) {
      this.FUS.uploadFile(this.file).subscribe((res: any) => {
          // Short link via api response
          this.postObj.Image = res;

          this.WriteToDatabase();
        }
      );
    } else {
      this.WriteToDatabase();
    }
  }

  WriteToDatabase() {
    let json = JSON.stringify(this.postObj)
    this.http.post(environment.serverUrl + '/createPost', json).subscribe((res: any) => {
      if (res.result) {
        alert("Posted Successfully!")
        window.location.reload();
      } else {
        alert(res.message)
      }
    })

  }

}

export class Post {
  Username: string;
  Title: string;
  Body: string;
  Image: string;
  constructor() {
    this.Username = '';
    this.Title = '';
    this.Body = '';
    this.Image = '';
  }
}