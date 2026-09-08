import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DocsPage } from './docs-page';

describe('DocsPage', () => {
  let component: DocsPage;
  let fixture: ComponentFixture<DocsPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DocsPage],
    }).compileComponents();

    fixture = TestBed.createComponent(DocsPage);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
