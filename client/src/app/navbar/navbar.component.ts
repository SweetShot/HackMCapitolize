import { DataService } from './../services/data.service';
import { LoginComponent } from './../components/login/login.component';
import { Component, OnInit } from '@angular/core';
import { Http } from '@angular/http';
import { HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
  public isLogin = false;
  public dataService: DataService;

  constructor(private dataService_: DataService) {
    this.dataService = dataService_;
  }

  ngOnInit() {
  }

  onLogin() {
    console.log('clicked');
  }

}
