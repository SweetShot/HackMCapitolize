import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
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
  constructor(private dataService: DataService, private router: Router ) { }

  ngOnInit() {
  }

  login(username, password) {
    console.log(username);
    this.dataService.postLogin(username, password).subscribe((status) => {
      this.loginStatus = status;
      this.dataService.auth = '1 ' + status.token;
      this.router.navigate(['/dashboard']);
      console.log(this.dataService.auth);
    });
  }
}
