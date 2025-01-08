import { Routes } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './Login/login/login.component';
import { RegisterComponent } from './Login/register/register.component';
import { PostsComponent } from './Posts/posts/posts.component';
import { FriendsListComponent } from './friends-list/friends-list.component';

export const routes: Routes = [
    { path: '', redirectTo: 'login', pathMatch:'full' },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: '', component: HomeComponent, 
        children: [     // Contains HomeComponent Template
            { path: 'dashboard', component: FriendsListComponent },
            { path: 'posts', component: PostsComponent },
        ]
     },
    { path: 'about', component: FriendsListComponent }
];
