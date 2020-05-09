import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { Subscription } from './model/subscription';
import { ResponseResult } from './model/response';
import { User } from './model/user';

@Injectable({
  providedIn: 'root'
})
export class WebService {

  private apiUrl = 'http://127.0.0.1:5003';
  public subCreated = new Subject<boolean>();

  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  }

  constructor(private http: HttpClient) { }


  // LOGIN API

  registerUser(user: User) {
    const url = `${this.apiUrl}/register`;
    return this.http.post<ResponseResult>(url, user).pipe(
      map(response => {
        console.log(response);
        if (response.error === "") {
          console.log("User registered successfuly");
        } else {
          console.log("There was an error in processing your request");
          console.log(response.error);
        }
      })
    )
  }

  loginUser(user: User) {
    console.log("S-a facut update 3");
    const url = `${this.apiUrl}/login`;
    var reqHeader = new HttpHeaders({
      'Content-Type': 'application/json',
      'Access-Control-Allow-Headers': 'accept, content-type',
      'Access-Control-Allow-Methods': 'GET, POST',
      'Access-Control-Allow-Origin': '*'
    });
    return this.http.post<ResponseResult>(url, user, { headers: reqHeader }).pipe(
      map(response => {
        console.log(response);
        if (response.error === "") {
          console.log("User logged in successfuly");
          // this.token = response.result;
          // this.username = user.username;
          user.token = response.result;
          this.saveInLocal(user);
          // localStorage.setItem('dataSource', "test local storage");
          var aux = this.getFromLocal();
          console.log("From local storage: " + aux);
        } else {
          console.log("There was an error in processing your request");
          console.log(response.error);
        }
      })
    )
  }

  logoutUser() {
    this.deleteLocal();
  }


  // OPERATIONS API
  getSubscriptions(): Observable<ResponseResult> {
    var aux = this.getFromLocal();
    console.log(aux.username)
    console.log(aux.token)
    const url = `${this.apiUrl}/subscriptions/${aux.username}`;
    var reqHeader = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + aux.token
    });
    return this.http.get<ResponseResult>(url, { headers: reqHeader }).pipe();
  }

  createSubscription(sub: Subscription) {
    var aux = this.getFromLocal();
    console.log(aux.username)
    console.log(aux.token)
    const url = `${this.apiUrl}/subscriptions/${aux.username}`;
    var reqHeader = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + aux.token
    });
    return this.http.post<ResponseResult>(url, sub, { headers: reqHeader }).pipe(
      map(_ => {
        console.log("Asset created");
      })
    )
  }


  // Util functions

  saveInLocal(user: User) {
    console.log('recieved user: ' + user);
    localStorage.setItem("user", JSON.stringify(user));
    // this.storage.set(key, val);
    // this.data[key]= this.storage.get(key);
  }

  getFromLocal() {
      return JSON.parse(localStorage.getItem("user"));
  }

  deleteLocal() {
    localStorage.removeItem("user");
  }
}
