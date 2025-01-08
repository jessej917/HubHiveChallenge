import {Component, OnInit} from '@angular/core';
import { CourseComponent } from '../course/course.component';
import {HelloWorldService} from '../hello-world.service';
import { AuthService } from '../Services/auth.service';

@Component({
  selector: 'app-courses',
  standalone: true,
  imports: [CourseComponent],
  providers: [HelloWorldService],
  templateUrl: './courses.component.html',
  styleUrl: './courses.component.css'
})
export class CoursesComponent implements OnInit {
  titles: any[] = [];;
  todos: any[] = [];

  constructor(private hw: HelloWorldService, private AS: AuthService) {}

  ngOnInit() {
    this.hw.getTodos().subscribe(response => {
      this.todos = response;
      console.log("Todos");
    });
    //console.log(this.todos);

    this.hw.getTitle().subscribe(data => {
      this.titles = data;
    });
    //console.log(this.title);

  }

}
