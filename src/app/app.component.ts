import { Component, Inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'esbuild-angular';

  constructor(
    private router: Router,
    private ar: ActivatedRoute,
  ) {}
}
