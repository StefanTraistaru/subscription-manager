import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { SubscriptionListComponent } from './components/subscription-list/subscription-list.component';


const routes: Routes = [
  { path: '', redirectTo: '/init', pathMatch: 'full' },
  { path: 'init', component: SubscriptionListComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
