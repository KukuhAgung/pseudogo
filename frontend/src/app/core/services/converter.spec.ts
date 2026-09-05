import { TestBed } from '@angular/core/testing';
import { Converter } from './converter';

describe('Converter', () => {
  let service: Converter;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Converter);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
