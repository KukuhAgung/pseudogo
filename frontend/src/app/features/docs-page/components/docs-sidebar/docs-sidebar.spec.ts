import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DocsSidebar } from './docs-sidebar';

describe('DocsSidebar', () => {
  let component: DocsSidebar;
  let fixture: ComponentFixture<DocsSidebar>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DocsSidebar],
    }).compileComponents();

    fixture = TestBed.createComponent(DocsSidebar);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
