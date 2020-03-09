import { Component, OnInit } from '@angular/core';

import { subscriptions } from '../../subscriptions'

@Component({
  selector: 'app-subscription-list',
  templateUrl: './subscription-list.component.html',
  styleUrls: ['./subscription-list.component.scss']
})
export class SubscriptionListComponent implements OnInit {

  public my_subs;

  constructor() { }

  ngOnInit(): void {
    this.my_subs = subscriptions;
  }

}
