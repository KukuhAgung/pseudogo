import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { Observable, map } from 'rxjs';

interface ConvertResponse {
  go_code?: string;
  error?: string;
}

@Service()
export class Converter {
  private http = inject(HttpClient);

  convert(pseudocode: string): Observable<string> {
    return this.http.post<ConvertResponse>('http://localhost:8080/convert', { pseudocode }).pipe(
      map((res) => {
        if (res.error) throw new Error(res.error);
        return res.go_code ?? '';
      }),
    );
  }
}
