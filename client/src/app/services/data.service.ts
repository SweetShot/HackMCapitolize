import { Injectable } from '@angular/core';
import { Http, HttpModule, Response } from '@angular/http';
// Add the RxJS Observable operators we need in this app.
import { map } from 'rxjs/operators';
import 'rxjs/add/operator/map';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  constructor(public http: Http) { }

  getIdeas() {
    return this.http.get('100.96.247.217:8081/Ideas').map(res => res.json());
  }
}
