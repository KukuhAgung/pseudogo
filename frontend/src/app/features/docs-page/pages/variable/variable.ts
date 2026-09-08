import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { remixArrowRightSLine, remixCodeBoxLine, remixTerminalWindowLine } from '@ng-icons/remixicon';
import { CodeBlock } from '../../../../shared/components/code-block/code-block';

@Component({
  imports: [NgIcon, CodeBlock, RouterLink],
  providers: [provideIcons({ remixArrowRightSLine, remixCodeBoxLine, remixTerminalWindowLine })],
  selector: 'app-variable',
  styleUrl: './variable.css',
  templateUrl: './variable.html',
})
export class Variable {
  pseudocodeExample = [
    'Program DataMahasiswa',
    '',
    'Kamus:',
    '  nama : string',
    '  umur : integer',
    '  IPK : real',
    '  isLulus : boolean',
    '  nilaiTugas, nilaiUjian : integer',
    '',
    'Algoritma:',
    '  // Mulai memproses data...',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    var nama string',
    '    var umur int',
    '    var IPK float64',
    '    var isLulus bool',
    '    var nilaiTugas, nilaiUjian int',
    '',
    '    // Mulai memproses data...',
    '}',
  ].join('\n');
}
