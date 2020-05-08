import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { WebService } from 'src/app/web.service'
import { throwIfEmpty } from 'rxjs/operators';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(
    public webService: WebService,
    private router: Router
  ) { }

  ngOnInit(): void {
  }

  logoutUser(): void {
    this.webService.logoutUser();
    this.router.navigateByUrl('/');
  }
}
