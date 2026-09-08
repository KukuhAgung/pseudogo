import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Looping } from './looping';

describe('Looping', () => {
  let component: Looping;
  let fixture: ComponentFixture<Looping>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Looping],
    }).compileComponents();

    fixture = TestBed.createComponent(Looping);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
