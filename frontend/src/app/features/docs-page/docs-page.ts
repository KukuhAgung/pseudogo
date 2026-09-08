import { Component, ElementRef, ViewChild, OnInit, OnDestroy } from '@angular/core';
import { filter, Subscription } from 'rxjs';
import { RouterOutlet, Router, NavigationEnd } from '@angular/router';
import { DocsSidebar } from './components/docs-sidebar/docs-sidebar';

@Component({
  imports: [DocsSidebar, RouterOutlet],
  selector: 'app-docs-page',
  styleUrl: './docs-page.css',
  templateUrl: './docs-page.html',
})
export class DocsPage implements OnInit, OnDestroy {
  @ViewChild('scrollContainer') scrollContainer!: ElementRef;

  private routerSub!: Subscription;

  constructor(private router: Router) {}

  ngOnInit() {
    this.routerSub = this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe(() => {
        if (this.scrollContainer) {
          this.scrollContainer.nativeElement.scrollTo({ top: 0, behavior: 'instant' });
        }
      });
  }

  ngOnDestroy() {
    if (this.routerSub) {
      this.routerSub.unsubscribe();
    }
  }
}
