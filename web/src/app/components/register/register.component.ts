import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { WebService } from 'src/app/web.service'
import { User } from '../../model/user';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss']
})
export class RegisterComponent implements OnInit {

  public data: User;
  showWarning = false;

  constructor(
    public webService: WebService,
    private router: Router
  ) { }


  ngOnInit(): void {
    this.data = new User("","","","","");
  }

  registerUser(): void {
    if (this.data.username === "" || this.data.firstname === "" ||
        this.data.lastname === "" || this.data.password === "") {
          this.showWarning = true;
          return;
        }
    this.showWarning = false;
    console.log(this.data);
    this.webService.registerUser(this.data).subscribe(
      _ => {
        console.log('done');
        this.data = new User("","","","","");
        this.router.navigateByUrl('/login');
      }
    )
  }

}
