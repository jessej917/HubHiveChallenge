import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Router } from '@angular/router';
import { HelloWorldService } from '../../hello-world.service';
import { AuthService } from '../../Services/auth.service';
import { UploadPostComponent } from '../upload-post/upload-post.component';
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [NgIf, FormsModule, UploadPostComponent],
  providers: [HelloWorldService],
  templateUrl: './posts.component.html',
  styleUrl: './posts.component.css'
})
export class PostsComponent implements OnInit {
  posts: any[] = [];
  postObj: Post;
  user: any;

  constructor(private http: HttpClient, private router: Router, private hw: HelloWorldService, private AS: AuthService) {
    this.postObj = new Post();
  }

  ngOnInit() {

    console.log(this.AS.getUser());
    this.user = this.AS.getUser();

    this.hw.getPosts().subscribe(data => {
      this.posts = data;
    });

    //console.log(this.title);
  }

  CreatePost() {
    let json = JSON.stringify(this.postObj)
    this.http.post(environment.serverUrl + '/createPost', json).subscribe((res: any) => {
      if (res.result) {
        alert("Created Post Successfully")
        this.AS.setUser(this.postObj.Title);
        this.router.navigateByUrl('/dashboard')
      } else {
        alert(res.message)
      }
    })
  }

  AddFriend(friend: string) {
    let person: Friend = {
      Username: this.AS.getUser(),
      Friend: friend
    };
    if (person.Friend != person.Username) {
      let json = JSON.stringify(person)
      this.http.post(environment.serverUrl + '/addFriend', json).subscribe((res: any) => {
        if (res.result) {
          alert("Friended " + friend)
          //this.router.navigateByUrl('/dashboard')
        } else {
          alert(res.message)
        }
      })
    }
  }

}

export class Post {
  Username: string;
  Title: string;
  Body: string;
  Image: string;
  Date: string
  constructor() {
    this.Username = '';
    this.Title = '';
    this.Body = '';
    this.Image = '';
    this.Date = '';
  }
}

export class Friend {
  Username: string;
  Friend: string;
  constructor() {
    this.Username = '';
    this.Friend = '';

  }
}
