import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { CodeBlock } from '../../../../shared/components/code-block/code-block';
import {
  remixCodeBoxLine,
  remixTerminalWindowLine,
  remixArrowRightSLine,
} from '@ng-icons/remixicon';


@Component({
  imports: [NgIcon, CodeBlock, RouterLink],
  providers: [provideIcons({ remixCodeBoxLine, remixTerminalWindowLine, remixArrowRightSLine })],
  selector: 'app-looping',
  styleUrl: './looping.css',
  templateUrl: './looping.html',
})
export class Looping {
  pseudocodeExample = [
    'Program CetakAngka',
    '',
    'Kamus:',
    '  i : integer',
    '  batas : integer',
    '',
    'Algoritma:',
    '  batas <- 3',
    '  ',
    '  output("--- Menggunakan For ---")',
    '  for i <- 1 to batas do',
    '    output("Angka ke-", i)',
    '  endfor',
    '',
    '  output("--- Menggunakan While ---")',
    '  i <- 1',
    '  while i <= batas do',
    '    output("Angka ke-", i)',
    '    i <- i + 1',
    '  endwhile',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    var i int',
    '    var batas int',
    '',
    '    batas = 3',
    '',
    '    fmt.Println("--- Menggunakan For ---")',
    '    for i = 1; i <= batas; i++ {',
    '        fmt.Println("Angka ke-", i)',
    '    }',
    '',
    '    fmt.Println("--- Menggunakan While ---")',
    '    i = 1',
    '    for i <= batas {',
    '        fmt.Println("Angka ke-", i)',
    '        i = i + 1',
    '    }',
    '}',
  ].join('\n');
}
