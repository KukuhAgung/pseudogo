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
  selector: 'app-operator',
  styleUrl: './operator.css',
  templateUrl: './operator.html',
})
export class Operator {
  pseudocodeExample = [
    'Program EvaluasiNilai',
    '',
    'Kamus:',
    '  nilaiAkhir : integer',
    '  sisaBagi : integer',
    '  isLulus : boolean',
    '',
    'Algoritma:',
    '  nilaiAkhir <- 85 + 10',
    '  ',
    '  // Menggunakan operator aritmatika MOD',
    '  sisaBagi <- nilaiAkhir MOD 2',
    '',
    '  // Menggunakan operator perbandingan dan logika',
    '  isLulus <- (nilaiAkhir >= 75) AND (sisaBagi = 1)',
    '',
    '  output("Status lulus: ", isLulus)',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    var nilaiAkhir int',
    '    var sisaBagi int',
    '    var isLulus bool',
    '',
    '    nilaiAkhir = 85 + 10',
    '',
    '    // Menggunakan operator aritmatika MOD',
    '    sisaBagi = nilaiAkhir % 2',
    '',
    '    // Menggunakan operator perbandingan dan logika',
    '    isLulus = (nilaiAkhir >= 75) && (sisaBagi == 1)',
    '',
    '    fmt.Println("Status lulus:", isLulus)',
    '}',
  ].join('\n');
}
