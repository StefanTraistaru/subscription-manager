import { Component, OnInit } from '@angular/core';

import { Subscription } from 'src/app/model/subscription'
// import { subscriptions } from '../../subscriptions'

import { WebService } from '../../web.service';

@Component({
  selector: 'app-subscription-list',
  templateUrl: './subscription-list.component.html',
  styleUrls: ['./subscription-list.component.scss']
})
export class SubscriptionListComponent implements OnInit {

  // public my_subs;
  my_subs: Subscription[];


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
        console.log("Got subscriptions.");
      }
    )
  }

}
