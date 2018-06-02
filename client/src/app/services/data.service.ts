import { Injectable } from '@angular/core';
import { Http, HttpModule, Response } from '@angular/http';
// Add the RxJS Observable operators we need in this app.
import { Observable, Subject, asapScheduler, pipe, of, from, interval, merge, fromEvent } from 'rxjs';
import { map, filter, scan } from 'rxjs/operators';
import 'rxjs/add/operator/map';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  auth: string;
  constructor(public http: Http) { }

  getIdeas() {
    return this.http.get('http://100.96.247.217:8081/Ideas').pipe(map(res => res.json()));
  }

  postLogin(username: string, password: string) {
    return this.http.post('http://100.96.247.217:8081/Login', username, password).pipe(map(res => res.json()));
  }
}
