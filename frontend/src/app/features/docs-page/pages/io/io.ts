import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { CodeBlock } from '../../../../shared/components/code-block/code-block';
import { remixCodeBoxLine, remixTerminalWindowLine, remixArrowRightSLine } from '@ng-icons/remixicon';

@Component({
  imports: [CodeBlock, NgIcon, RouterLink],
  providers: [provideIcons({ remixArrowRightSLine, remixCodeBoxLine, remixTerminalWindowLine })],
  selector: 'app-io',
  styleUrl: './io.css',
  templateUrl: './io.html',
})
export class Io {
  pseudocodeExample = [
    'Program SapaPengguna',
    '',
    'Kamus:',
    '  nama : string',
    '  umur : integer',
    '',
    'Algoritma:',
    '  output("Siapa nama kamu? ")',
    '  input(nama)',
    '',
    '  output("Berapa umurmu? ")',
    '  input(umur)',
    '',
    '  output("Halo ", nama, ", umur kamu ", umur, " tahun.")',
  ].join('\n');

  golangExample = [
    'package main',
    '',
    'import "fmt"',
    '',
    'func main() {',
    '    var nama string',
    '    var umur int',
    '',
    '    fmt.Print("Siapa nama kamu? ")',
    '    fmt.Scan(&nama)',
    '',
    '    fmt.Print("Berapa umurmu? ")',
    '    fmt.Scan(&umur)',
    '',
    '    fmt.Println("Halo ", nama, ", umur kamu ", umur, " tahun.")',
    '}',
  ].join('\n');
}
