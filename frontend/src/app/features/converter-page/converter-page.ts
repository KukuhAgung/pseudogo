import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Converter } from '../../core/services/converter';

@Component({
  imports: [FormsModule],
  standalone: true,
  selector: 'app-converter-page',
  styleUrl: './converter-page.css',
  templateUrl: './converter-page.html',
})
export class ConverterPage {
  private converter = inject(Converter);

  pseudocode = signal('');
  goCode = signal('');
  errorMessage = signal('');
  isLoading = signal(false);

  onConvert() {
    this.isLoading.set(true);
    this.errorMessage.set('');
    this.converter.convert(this.pseudocode()).subscribe({
      next: (result) => {
        this.goCode.set(result);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.errorMessage.set(err.message ?? 'Terjadi Kesalahan');
        this.isLoading.set(false);
      },
    });
  }
}
