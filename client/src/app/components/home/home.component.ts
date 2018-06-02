import { DataService } from './../../services/data.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  Ideas: Idea[];

  constructor(private dataService: DataService ) { }

  ngOnInit() {
    this.dataService.getIdeas().subscribe((Ideas) => {
      console.log(Ideas);
      this.Ideas = Ideas;
    });
  }

}

interface SupportOptions {
  title: string;
  description: string;
  delivery: boolean;
  expected_date: number;
  price: number;
}

interface Idea {
  username: string;
  title: string;
  image: ImageData;
  description: string;
  total_funds_required: number;
  total_funds_raised: number;
  date_posted: number;
  date_end: number;
  num_contributors: number;
  beneficiary: string;
  category: string;
  summary: string;
  support_options: SupportOptions[];
}
