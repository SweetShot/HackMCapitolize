import { DataService } from './../../services/data.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-start-fr',
  templateUrl: './start-fr.component.html',
  styleUrls: ['./start-fr.component.css']
})
export class StartFRComponent implements OnInit {
  constructor(private dataService: DataService ) { }

  ngOnInit() {
    if ((this.dataService.auth) !== '') {
      console.log('Login Failed');
    } else {
      console.log('Logged in');
    }
  }
}
