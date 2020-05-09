import { Component, OnInit, Inject } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDatepicker } from '@angular/material/datepicker';

import { WebService } from 'src/app/web.service'
import { Subscription } from 'src/app/model/subscription'

@Component({
  selector: 'app-create-subscription',
  templateUrl: './create-subscription.component.html',
  styleUrls: ['./create-subscription.component.scss']
})
export class CreateSubscriptionComponent implements OnInit {

  public data: Subscription;
  showWarning = false;

  constructor(
    public webService: WebService,
  ) { }

  ngOnInit(): void {
    this.data = new Subscription("","","","","","");
  }

  createAsset(): void {
    if (this.data.name === "" || this.data.price === "" || this.data.details === "" ||
        this.data.date_d === "" || this.data.date_m === "" || this.data.date_y === "") {
          this.showWarning = true;
          return;
        }
    this.showWarning = false;
    console.log(this.data);
    this.webService.createSubscription(this.data).subscribe(
      _ => {
        console.log('done');
        this.data = new Subscription("","","","","","");
        this.webService.subCreated.next(true);
      }
    )
  }

}
