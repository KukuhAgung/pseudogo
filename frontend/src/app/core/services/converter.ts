import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { Observable, map } from 'rxjs';

interface ConvertResponse {
  go_code?: string;
  error?: string;
}

interface RunResponse {
  output?: string;
  stderr?: string;
  timed_out?: boolean;
  error?: string;
}

export interface RunResult {
  output: string;
  stderr: string;
  timedOut: boolean;
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

  run(pseudocode: string, input: string): Observable<RunResult> {
    return this.http.post<RunResponse>('http://localhost:8080/run', { pseudocode, input }).pipe(
      map((res) => {
        if (res.error) throw new Error(res.error);
        return {
          output: res.output ?? '',
          stderr: res.stderr ?? '',
          timedOut: res.timed_out ?? false,
        };
      }),
    );
  }
}
