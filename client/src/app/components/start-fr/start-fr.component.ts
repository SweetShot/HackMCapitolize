import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { DataService } from './../../services/data.service';
import { Component, OnInit } from '@angular/core';
import { LoginComponent } from '../login/login.component';

@Component({
  selector: 'app-start-fr',
  templateUrl: './start-fr.component.html',
  styleUrls: ['./start-fr.component.css']
})
export class StartFRComponent implements OnInit {
  constructor(private dataService: DataService, private router: Router ) {
  }

  ngOnInit() {
  }
}
