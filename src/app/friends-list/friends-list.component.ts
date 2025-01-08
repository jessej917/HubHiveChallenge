import { Component } from '@angular/core';
import { HelloWorldService } from '../hello-world.service';
import { AuthService } from '../Services/auth.service';

@Component({
  selector: 'app-friends-list',
  standalone: true,
  imports: [],
  providers: [HelloWorldService, AuthService],
  templateUrl: './friends-list.component.html',
  styleUrl: './friends-list.component.css'
})
export class FriendsListComponent {
  friends: any[] = [];

  constructor(private hw: HelloWorldService, private AS: AuthService) { }

  ngOnInit() {
    this.hw.getFriends(this.AS.getUser()).subscribe(data => {
      this.friends = data;
    });
    //console.log(this.title);

  }
}
