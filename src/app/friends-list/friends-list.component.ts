import { Component } from '@angular/core';
import { AuthService } from '../Services/auth.service';
import { HelloWorldService } from '../hello-world.service';
import { Friend } from '../Posts/posts/posts.component';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

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

  constructor(private http: HttpClient, private hw: HelloWorldService, private AS: AuthService) { }

  ngOnInit() {
    this.hw.getFriends(this.AS.getUser()).subscribe(data => {
      this.friends = data;
    });

  }

  RemoveFriend(friend: string) {
      let person: Friend = {
        Username: this.AS.getUser(),
        Friend: friend,
        Remove: true
      };
      if (person.Friend != person.Username) {
        let json = JSON.stringify(person)
        this.hw.AddFriend(json).subscribe((res: any) => {
          if (res.result) {
            alert("Removed Friend: " + friend)
            window.location.reload();
          } else {
            alert(res.message)
          }
        })
      }
    }
}
