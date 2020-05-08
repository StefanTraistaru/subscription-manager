import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { WebService } from 'src/app/web.service'
import { User } from '../../model/user';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  public data: User;
  showWarning = false;

  constructor(
    public webService: WebService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.data = new User("","","","","");
  }

  loginUser(): void {
    if (this.data.username === "" || this.data.password === "") {
          this.showWarning = true;
          return;
        }
    this.showWarning = false;
    console.log(this.data);
    this.webService.loginUser(this.data).subscribe(
      _ => {
        console.log('done');
        this.data = new User("","","","","");
        this.router.navigateByUrl('/main');
      }
    )
  }

}
