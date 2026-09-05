import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ConverterPage } from './converter-page';

describe('ConverterPage', () => {
  let component: ConverterPage;
  let fixture: ComponentFixture<ConverterPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConverterPage],
    }).compileComponents();

    fixture = TestBed.createComponent(ConverterPage);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
