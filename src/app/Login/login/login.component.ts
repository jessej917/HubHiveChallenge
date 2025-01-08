import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {environment} from '../../../environments/environment';
import { Router } from '@angular/router';
import { AuthService } from '../../Services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {

  loginObj: Login;
  constructor(private http: HttpClient, private router: Router, private AS: AuthService) {
    this.loginObj = new Login();
  }

  ngOnInit() {
    console.log(this.AS.getUser());
    if(this.AS.getUser() != null) {
      this.router.navigateByUrl('/dashboard')
    }
  }

  onLogin() {
    let json = JSON.stringify(this.loginObj)
    this.http.post(environment.serverUrl + '/login', json).subscribe((res:any)=>{
      if(res.result) {
        alert("Login Success")
        this.AS.setUser(this.loginObj.Username);
        this.router.navigateByUrl('/dashboard')
      } else {
        alert(res.message)
      }
    })
  }
}

export class Login {
  Username: string;
  Password: string;
  constructor() {
    this.Username = '';
    this.Password = '';
  }
}