import {Component, OnInit} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import {HelloWorldService} from './hello-world.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
  ],
  providers: [HelloWorldService],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  titles: any[] = [];;
  todos: any[] = [];

  constructor(private hw: HelloWorldService) {}

  ngOnInit() {
    this.hw.getTodos().subscribe(response => {
      this.todos = response;
    });
    //console.log(this.todos);

    // this.hw.getTitle().subscribe(data => {
    //   this.titles = data;
    // });

    //console.log(this.title);
  }

}
