import { StreamLanguage, HighlightStyle } from '@codemirror/language';
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
  'type',
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
  'endif',
  'endfor',
  'endwhile',
  'endprocedure',
  'endfunction',
]);

export const pseudocodeLanguage = StreamLanguage.define({
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

export const fixedDarkHighlight = HighlightStyle.define([
  { tag: t.keyword, color: '#c084fc' },
  { tag: t.string, color: '#86efac' },
  { tag: t.comment, color: '#9ca3af', fontStyle: 'italic' },
  { tag: t.number, color: '#fbbf24' },
  { tag: t.variableName, color: '#e5e5e5' },
  { tag: t.operator, color: '#f472b6' },
  { tag: t.function(t.variableName), color: '#60a5fa' },
  { tag: t.typeName, color: '#60a5fa' },
]);

export const adaptiveHighlight = HighlightStyle.define([
  { tag: t.keyword, color: 'var(--color-syntax-keyword)' },
  { tag: t.string, color: 'var(--color-syntax-string)' },
  { tag: t.comment, color: 'var(--color-syntax-comment)', fontStyle: 'italic' },
  { tag: t.number, color: 'var(--color-syntax-number)' },
  { tag: t.variableName, color: 'var(--color-strong)' },
  { tag: t.operator, color: 'var(--color-syntax-operator)' },
]);
