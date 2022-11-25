import { Component, inject, Inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'The Fastest Builder ever!!!';

  private ar: ActivatedRoute = inject(ActivatedRoute);

  private router: Router = inject(Router);

  constructor(
  ) {}
}
