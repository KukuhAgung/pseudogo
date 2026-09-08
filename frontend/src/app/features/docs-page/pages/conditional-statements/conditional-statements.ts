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
  selector: 'app-conditional-statements',
  styleUrl: './conditional-statements.css',
  templateUrl: './conditional-statements.html',
})
export class ConditionalStatements {
  pseudocodeExample = [
    'Program CekKelulusan',
    '',
    'Kamus:',
    '  nilai : integer',
    '',
    'Algoritma:',
    '  output("Masukkan nilai ujian: ")',
    '  input(nilai)',
    '',
    '  if nilai >= 80 then',
    '    output("Grade A: Lulus dengan sangat baik")',
    '  else if nilai >= 60 then',
    '    output("Grade B: Lulus")',
    '  else',
    '    output("Grade C: Mengulang")',
    '  endif',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    var nilai int',
    '',
    '    fmt.Print("Masukkan nilai ujian: ")',
    '    fmt.Scan(&nilai)',
    '',
    '    if nilai >= 80 {',
    '        fmt.Println("Grade A: Lulus dengan sangat baik")',
    '    } else if nilai >= 60 {',
    '        fmt.Println("Grade B: Lulus")',
    '    } else {',
    '        fmt.Println("Grade C: Mengulang")',
    '    }',
    '}',
  ].join('\n');
}
