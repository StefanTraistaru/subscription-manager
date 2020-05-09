import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { User } from 'src/app/model/user';
import { WebService } from 'src/app/web.service'
import { throwIfEmpty } from 'rxjs/operators';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  showCreateForm = false;
  username: string;

  constructor(
    public webService: WebService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.getUsername();
  }

  logoutUser(): void {
    this.webService.logoutUser();
    this.router.navigateByUrl('/');
  }

  openCreateForm() {
    this.showCreateForm = true;
  }

  closeCreateForm() {
    this.showCreateForm = false;
  }

  getUsername() {
    var aux: User;
    aux = this.getFromLocal();
    this.username = aux.username;
  }

  getFromLocal() {
    return JSON.parse(localStorage.getItem("user"));
  }
}
