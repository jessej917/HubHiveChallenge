import {Component, OnInit} from '@angular/core';
import {HelloWorldService} from '../hello-world.service';
import { Router, RouterOutlet } from '@angular/router';
import { AuthService } from '../Services/auth.service';
import { PostsComponent } from '../Posts/posts/posts.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [PostsComponent, RouterOutlet],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})

export class HomeComponent implements OnInit {
  posts: any[] = [];
  user: any;

  constructor(private router: Router, private hw: HelloWorldService, private AS: AuthService) {}

  ngOnInit() {
    
    console.log(this.AS.getUser());
    this.user = this.AS.getUser();
    if(this.user == null) {
      this.router.navigateByUrl('/')
    }

    this.hw.getPosts().subscribe(data => {
      this.posts = data;
    });

  }

  LogOut() {
    this.AS.setUser(null);
    this.router.navigateByUrl('/')
  }

}
