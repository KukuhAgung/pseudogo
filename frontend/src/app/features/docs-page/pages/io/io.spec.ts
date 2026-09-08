import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Io } from './io';

describe('Io', () => {
  let component: Io;
  let fixture: ComponentFixture<Io>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Io],
    }).compileComponents();

    fixture = TestBed.createComponent(Io);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
