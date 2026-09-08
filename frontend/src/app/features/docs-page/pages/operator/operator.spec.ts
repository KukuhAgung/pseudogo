import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Operator } from './operator';

describe('Operator', () => {
  let component: Operator;
  let fixture: ComponentFixture<Operator>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Operator],
    }).compileComponents();

    fixture = TestBed.createComponent(Operator);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
