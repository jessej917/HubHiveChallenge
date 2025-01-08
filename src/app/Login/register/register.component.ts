import { Component } from '@angular/core';
import { Login } from '../login/login.component';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import {environment} from '../../../environments/environment';
import { AuthService } from '../../Services/auth.service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {

  loginObj: Login;
  constructor(private http: HttpClient, private router: Router, private AS: AuthService) {
    this.loginObj = new Login();
  }

  onLogin() {
    let json = JSON.stringify(this.loginObj)
    this.http.post(environment.serverUrl + '/createUser', json).subscribe((res:any)=>{
      if(res.result) {
        alert("Registered Successfully!")
        this.AS.setUser(this.loginObj.Username);
        this.router.navigateByUrl('/dashboard')
      } else {
        alert(res.message)
      }
    })
  }
}
