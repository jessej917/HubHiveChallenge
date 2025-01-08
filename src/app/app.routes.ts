import { Routes } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { AppComponent } from './app.component';
import { CoursesComponent } from './courses/courses.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './Login/login/login.component';
import { RegisterComponent } from './Login/register/register.component';
import { PostsComponent } from './Posts/posts/posts.component';

export const routes: Routes = [
    { path: '', redirectTo: 'login', pathMatch:'full' },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: '', component: HomeComponent, 
        children: [     // Contains HomeComponent Template
            { path: 'dashboard', component: CoursesComponent },
            { path: 'posts', component: PostsComponent },
        ]
     },
    { path: 'about', component: CoursesComponent }
];
