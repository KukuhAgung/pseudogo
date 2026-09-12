import { Service, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map, catchError, throwError } from 'rxjs';

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

function extractErrorMessage(err: any): string {
  return err.error?.error ?? 'Terjadi kesalahan tak terduga';
}

@Service()
export class Converter {
  private http = inject(HttpClient);

  convert(pseudocode: string): Observable<string> {
    return this.http.post<ConvertResponse>('http://localhost:8080/convert', { pseudocode }).pipe(
      map((res) => res.go_code ?? ''),
      catchError((err) => throwError(() => new Error(extractErrorMessage(err)))),
    );
  }

  run(pseudocode: string, input: string): Observable<RunResult> {
    return this.http.post<RunResponse>('http://localhost:8080/run', { pseudocode, input }).pipe(
      map((res) => ({
        output: res.output ?? '',
        stderr: res.stderr ?? '',
        timedOut: res.timed_out ?? false,
      })),
      catchError((err) => throwError(() => new Error(extractErrorMessage(err)))),
    );
  }
}
