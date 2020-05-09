import { Component, OnInit, ViewChild } from '@angular/core';

import { Subscription } from 'src/app/model/subscription'
import { ResponseResult } from 'src/app/model/response'

import { WebService } from '../../web.service';
import { MatTableDataSource } from '@angular/material/table';
import {MatPaginator} from '@angular/material/paginator';
import {MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';

@Component({
  selector: 'app-subscription-list',
  templateUrl: './subscription-list.component.html',
  styleUrls: ['./subscription-list.component.scss']
})
export class SubscriptionListComponent implements OnInit {

  my_subs: Subscription[];
  displayedColumns: string[] = ['name', 'details', 'date', 'price'];
  expandedElement: Subscription | null;
  dataSource = new MatTableDataSource<Subscription>(this.my_subs)
  totalCosts: number;

  @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;

  constructor(
    public webService: WebService
  ) { }

  ngOnInit(): void {
    this.getSubscriptions();
    this.webService.subCreated.subscribe(
      _ => {
        this.refreshWindow();
      }
    );
  }

  getSubscriptions(): void {
    this.webService.getSubscriptions().subscribe(
      res => {
        console.log(res);
        this.my_subs = res.data;
        this.dataSource = new MatTableDataSource<Subscription>(this.my_subs)
        this.dataSource.paginator = this.paginator;
        console.log("Got subscriptions.");
        this.calculateCosts();
      }
    )
  }

  calculateCosts(): void {
    this.totalCosts = 0;
    this.my_subs.forEach(sub => {
      if (!isNaN(+sub.price)) {
        this.totalCosts += +sub.price;
      }
    });
  }

  refreshWindow() {
    this.getSubscriptions();
  }

}
