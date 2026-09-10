import {
  Component,
  ElementRef,
  ViewChild,
  OnInit,
  OnDestroy,
  inject,
  AfterViewInit,
  signal,
} from '@angular/core';
import { filter, Subscription } from 'rxjs';
import { RouterOutlet, Router, NavigationEnd } from '@angular/router';
import { DocsSidebar } from './components/docs-sidebar/docs-sidebar';
import { UiStateService } from '../../core/services/ui-state';

@Component({
  imports: [DocsSidebar, RouterOutlet],
  standalone: true,
  selector: 'app-docs-page',
  styleUrl: './docs-page.css',
  templateUrl: './docs-page.html',
})
export class DocsPage implements AfterViewInit, OnInit, OnDestroy {
  ui = inject(UiStateService);
  private el = inject(ElementRef);
  private observer: IntersectionObserver | null = null;

  sections = signal<{ id: string; title: string; level: string }[]>([]);
  activeSection = signal<string>('');

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
        setTimeout(() => {
          this.setupScrollSpy();
        }, 50);
      });
  }

  ngAfterViewInit() {
    setTimeout(() => {
      this.setupScrollSpy();
    }, 50);
  }

  private setupScrollSpy() {
    if (this.observer) {
      this.observer.disconnect();
    }

    const headings = Array.from(this.el.nativeElement.querySelectorAll('h2, h3')) as HTMLElement[];

    const parsedSections = headings
      .filter((h) => h.id)
      .map((h) => ({
        id: h.id,
        title: h.innerText,
        level: h.tagName.toLowerCase(),
      }));

    this.sections.set(parsedSections);

    const options = {
      root: this.scrollContainer ? this.scrollContainer.nativeElement : null,
      rootMargin: '-10% 0px -40% 0px',
      threshold: 0,
    };

    this.observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          this.activeSection.set(entry.target.id);
        }
      });
    }, options);

    headings.forEach((h) => this.observer?.observe(h));
  }

  scrollTo(id: string, event: Event) {
    event.preventDefault();
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
      this.activeSection.set(id);
    }
  }

  ngOnDestroy() {
    if (this.routerSub) {
      this.routerSub.unsubscribe();
    }
    this.observer?.disconnect();
  }
}
