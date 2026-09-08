import { Component, ElementRef, viewChild, afterNextRender, input } from '@angular/core';
import { EditorState } from '@codemirror/state';
import { EditorView } from '@codemirror/view';
import { StreamLanguage, HighlightStyle, syntaxHighlighting } from '@codemirror/language';
import { go } from '@codemirror/lang-go';
import { tags as t } from '@lezer/highlight';

const PSEUDOCODE_KEYWORDS = new Set([
  'program',
  'kamus',
  'algoritma',
  'endprogram',
  'constant',
  'integer',
  'real',
  'boolean',
  'char',
  'string',
  'array',
  'of',
  'input',
  'output',
  'if',
  'then',
  'else',
  'end',
  'for',
  'to',
  'do',
  'while',
  'until',
  'repeat',
  'procedure',
  'function',
  'return',
  'in',
  'out',
  'inout',
  'and',
  'or',
  'not',
  'div',
  'mod',
  'true',
  'false',
]);

const pseudocodeLanguage = StreamLanguage.define({
  token(stream) {
    if (stream.match('//')) {
      stream.skipToEnd();
      return 'comment';
    }
    if (stream.match(/^"([^"\\]|\\.)*"/)) return 'string';
    if (stream.match(/^\d+(\.\d+)?/)) return 'number';
    if (stream.match(/^[A-Za-z_][A-Za-z0-9_]*/)) {
      const word = stream.current().toLowerCase();
      return PSEUDOCODE_KEYWORDS.has(word) ? 'keyword' : 'variableName';
    }
    if (stream.match('<-') || stream.match('->')) return 'operator';
    stream.next();
    return null;
  },
});

const modernHighlight = HighlightStyle.define([
  { tag: t.keyword, color: '#c084fc' },
  { tag: t.string, color: '#86efac' },
  { tag: t.comment, color: '#9ca3af', fontStyle: 'italic' },
  { tag: t.number, color: '#fbbf24' },
  { tag: t.variableName, color: '#e5e5e5' },
  { tag: t.operator, color: '#f472b6' },
  { tag: t.function(t.variableName), color: '#60a5fa' },
  { tag: t.typeName, color: '#60a5fa' },
]);

@Component({
  selector: 'app-code-block',
  standalone: true,
  template: `<div #host class="h-full overflow-auto"></div>`,
  styleUrl: './code-block.css',
})
export class CodeBlock {
  code = input.required<string>();
  language = input<'pseudocode' | 'go'>('pseudocode');

  private host = viewChild.required<ElementRef<HTMLDivElement>>('host');

  constructor() {
    afterNextRender(() => this.initEditor());
  }

  private initEditor() {
    const langExtension = this.language() === 'go' ? go() : pseudocodeLanguage;

    const state = EditorState.create({
      doc: this.code(),
      extensions: [
        langExtension,
        syntaxHighlighting(modernHighlight),
        EditorState.readOnly.of(true),
        EditorView.editable.of(false),
        EditorView.theme({
          '&': { backgroundColor: '#121212', height: '100%' },
          '.cm-content': {
            padding: '1.25rem',
            fontFamily: 'var(--font-mono)',
            fontSize: '13px',
            color: '#e5e5e5', 
          },
          '.cm-line': { padding: '0' },
        }),
      ],
    });

    new EditorView({ state, parent: this.host().nativeElement });
  }
}
