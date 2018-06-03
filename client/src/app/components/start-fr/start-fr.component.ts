import { FormsModule } from '@angular/forms';
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
  public Idea: Idea;
  constructor(private dataService: DataService, private router: Router ) {
  }

  ngOnInit() {
  }

  newIdea(idea: Idea) {
    console.log(idea);
    this.dataService.postIdeas(idea).subscribe((status) => {
      this.router.navigate(['/dashboard']);
      console.log(status);
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

