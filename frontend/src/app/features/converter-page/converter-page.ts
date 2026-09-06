import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Converter } from '../../core/services/converter';
import { CodeEditor } from '../../shared/components/code-editor/code-editor';

@Component({
  selector: 'app-converter-page',
  standalone: true,
  imports: [FormsModule, CodeEditor],
  templateUrl: './converter-page.html',
  styleUrl: './converter-page.css',
})
export class ConverterPage {
  private converter = inject(Converter);

  pseudocode = signal(['Program NamaProgram', 'Kamus:', '    ', 'Algoritma:', '    ', 'EndProgram'].join('\n'));
  goCode = signal('');
  errorMessage = signal('');
  isLoading = signal(false);

  stdinInput = signal('');
  runOutput = signal('');
  runStderr = signal('');
  runTimedOut = signal(false);
  runError = signal('');
  isRunning = signal(false);

  onConvert() {
    this.isLoading.set(true);
    this.errorMessage.set('');
    this.converter.convert(this.pseudocode()).subscribe({
      next: (result) => {
        this.goCode.set(result);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err.message ?? 'Terjadi kesalahan');
        this.isLoading.set(false);
      },
    });
  }

  onRun() {
    this.isRunning.set(true);
    this.runError.set('');
    this.runOutput.set('');
    this.runStderr.set('');
    this.runTimedOut.set(false);
    this.converter.run(this.pseudocode(), this.stdinInput()).subscribe({
      next: (result) => {
        this.runOutput.set(result.output);
        this.runStderr.set(result.stderr);
        this.runTimedOut.set(result.timedOut);
        this.isRunning.set(false);
      },
      error: (err) => {
        this.runError.set(err.message ?? 'Terjadi kesalahan');
        this.isRunning.set(false);
      },
    });
  }
}
