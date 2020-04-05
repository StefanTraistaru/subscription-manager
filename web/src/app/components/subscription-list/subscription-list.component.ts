import { Component, OnInit, ViewChild } from '@angular/core';

import { Subscription } from 'src/app/model/subscription'
// import { subscriptions } from '../../subscriptions'

import { WebService } from '../../web.service';
import { MatTableDataSource } from '@angular/material/table';
import {MatPaginator} from '@angular/material/paginator';
import {MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';

const ELEMENT_DATA: Subscription[] = [
  {name: 'Hydrogen', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Helium', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Lithium', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Beryllium', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Boron', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Carbon', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Nitrogen', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Oxygen', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Fluorine', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
  {name: 'Neon', price: '12', details: 'asd', date_d: 'a', date_m: 's', date_y: 'd'},
];


@Component({
  selector: 'app-subscription-list',
  templateUrl: './subscription-list.component.html',
  styleUrls: ['./subscription-list.component.scss']
})
export class SubscriptionListComponent implements OnInit {

  // public my_subs;
  my_subs: Subscription[];
  // displayedColumns: string[] = ['name', 'price', 'day', 'month', 'year'];
  displayedColumns: string[] = ['name', 'price', 'details', 'date'];
  expandedElement: Subscription | null;

  dataSource = new MatTableDataSource<Subscription>(this.my_subs)
  // dataSource = ELEMENT_DATA;

  @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;

  constructor(
    public webService: WebService
  ) { }

  ngOnInit(): void {
    // this.my_subs = subscriptions;
    this.getSubscriptions();
  }

  getSubscriptions(): void {
    this.webService.getSubscriptions().subscribe(
      subs => {
        this.my_subs = subs;
        this.dataSource = new MatTableDataSource<Subscription>(this.my_subs)
        this.dataSource.paginator = this.paginator;
        console.log("Got subscriptions.");
      }
    )
  }

}
