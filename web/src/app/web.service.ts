import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { Subscription } from './model/subscription'

@Injectable({
  providedIn: 'root'
})
export class WebService {

  private apiUrl = 'http://localhost:5000'

  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  }

  constructor(private http: HttpClient) { }

  getSubscriptions(): Observable<Subscription[]> {
    const url = `${this.apiUrl}/subscriptions`;
    return this.http.get<Subscription[]>(url).pipe();
  }

  createSubscription(sub: Subscription) {
    const url = `${this.apiUrl}/subscriptions`;
    return this.http.post<Subscription>(url, sub).pipe(
      map(_ => {
        console.log("Asset created");
      })
    )
  }


}
