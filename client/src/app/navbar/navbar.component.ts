import { LoginComponent } from './../components/login/login.component';
import { Component, OnInit } from '@angular/core';
import { Http } from '@angular/http';
import { HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { DataService } from '../services/data.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
  public isLogin = false;

  constructor(private dataService: DataService) {
    if (this.dataService.auth !== '') {
      this.isLogin = true;
      console.log('Logged In');
    } else {
      this.isLogin = false;
      console.log('Logged Out');
    }
  }

  ngOnInit() {
  }

  onLogin() {
    console.log('clicked');
  }

}
