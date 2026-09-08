import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { remixArrowRightSLine, remixCodeBoxLine, remixTerminalWindowLine } from '@ng-icons/remixicon';
import { CodeBlock } from '../../../../shared/components/code-block/code-block';

@Component({
  imports: [NgIcon, CodeBlock, RouterLink],
  providers: [provideIcons({ remixArrowRightSLine, remixCodeBoxLine, remixTerminalWindowLine })],
  selector: 'app-introduction',
  styleUrl: './introduction.css',
  templateUrl: './introduction.html',
})
export class Introduction {
  pseudocodeExample = [
    'Program HitungLuas',
    '',
    'Kamus:',
    '  // Deklarasi variabel diletakkan di sini',
    '',
    'Algoritma:',
    '  // Logika utama diletakkan di sini',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    // Deklarasi variabel akan berada di sini',
    '',
    '    // Logika utama akan berada di sini',
    '}',
  ].join('\n');
}
