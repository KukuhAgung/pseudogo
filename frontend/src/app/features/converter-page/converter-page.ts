import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Converter } from '../../core/services/converter';
import { CodeEditor } from '../../shared/components/code-editor/code-editor';
import { UiStateService } from '../../core/services/ui-state';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { remixTerminalFill, remixCloseFill, remixErrorWarningLine, remixFileCodeFill } from '@ng-icons/remixicon';
import { boxGoLangLogo } from '@ng-icons/boxicons/logos';

@Component({
  selector: 'app-converter-page',
  standalone: true,
  imports: [FormsModule, CodeEditor, NgIcon],
  providers: [
    provideIcons({ remixTerminalFill, remixCloseFill, remixErrorWarningLine, boxGoLangLogo, remixFileCodeFill }),
  ],
  templateUrl: './converter-page.html',
  styleUrl: './converter-page.css',
})
export class ConverterPage {
  ui = inject(UiStateService);
  private converter = inject(Converter);

  pseudocode = signal(
    ['Program NamaProgram', 'Kamus:', '    ', 'Algoritma:', '    ', 'EndProgram'].join('\n'),
  );
  goCode = signal('');
  errorMessage = signal<string | null>(null);
  isLoading = signal(false);

  stdinInput = signal('');
  runOutput = signal('');
  runStderr = signal('');
  runTimedOut = signal(false);
  runError = signal('');
  isRunning = signal(false);

  clearError() {
    this.errorMessage.set(null);
    this.runError.set('');
    this.runStderr.set('');
  }

  onConvert() {
    this.isLoading.set(true);
    this.errorMessage.set('');
    this.stdinInput.set('');
    this.converter.convert(this.pseudocode()).subscribe({
      next: (result) => {
        this.goCode.set(result);
        this.isLoading.set(false);
        const needsInput = result.includes('fmt.Scan') || result.includes('fmt.Scanf');

        if (needsInput) {
          this.ui.flagUnreadIfClosed();
        }
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
        this.ui.flagUnreadIfClosed();
      },
      error: (err) => {
        this.runError.set(err.message ?? 'Terjadi kesalahan');
        this.isRunning.set(false);
        this.ui.sidebarOpen.set(false);
      },
    });
  }
}
