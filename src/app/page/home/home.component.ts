import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  router: Router = inject(Router);

  ar: ActivatedRoute = inject(ActivatedRoute);

  constructor() { }

  ngOnInit(): void {
    console.log("Home inited Bond, James Bond!");
  }

}
