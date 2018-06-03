import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  constructor() { }

  ngOnInit() {
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
