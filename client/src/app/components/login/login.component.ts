import { Http } from '@angular/http';
import { DataService } from './../../services/data.service';
import { Component, OnInit } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  loginStatus: string;
  token: string;
  constructor(private dataService: DataService ) { }

  ngOnInit() {
  }

  login(username, password) {
    this.dataService.postLogin(username, password).subscribe((status) => {
      this.loginStatus = status;
      console.log(status);
    });
  }
}
