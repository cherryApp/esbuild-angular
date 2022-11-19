import { Component, Inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'The Fastest Builder ever!!!';

  constructor(
    private router: Router,
    private ar: ActivatedRoute,
  ) {}
}
