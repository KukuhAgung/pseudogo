import { Component, ElementRef, viewChild, afterNextRender, model, effect, input } from '@angular/core';
import { EditorState } from '@codemirror/state';
import { EditorView, lineNumbers, keymap, gutter, GutterMarker } from '@codemirror/view';
import { defaultKeymap, indentWithTab } from '@codemirror/commands';
import { syntaxHighlighting } from '@codemirror/language';
import { go } from '@codemirror/lang-go'
import { pseudocodeLanguage, codeEditorHighlight } from '../../code-mirror/pseudocode-language';

const INDENT_UNIT = '    ';

function getIndent(line: string): string {
  return line.match(/^[ \t]*/)?.[0] ?? '';
}

function computeNewLineIndent(prevLineText: string): string {
  const trimmed = prevLineText.trim();
  const currentIndent = getIndent(prevLineText);
  const opensBlock =
    trimmed.endsWith('then') ||
    trimmed.endsWith('do') ||
    trimmed === 'Kamus:' ||
    trimmed === 'Algoritma:' ||
    trimmed === 'repeat';
  return opensBlock ? currentIndent + INDENT_UNIT : currentIndent;
}

class ActiveLineDotMarker extends GutterMarker {
  override toDOM() {
    const dot = document.createElement('div');
    dot.style.width = '6px';
    dot.style.height = '6px';
    dot.style.borderRadius = '50%';
    dot.style.backgroundColor = 'var(--color-soft)';
    return dot;
  }
}

const activeLineDotMarker = new ActiveLineDotMarker();

const activeLineDotGutter = gutter({
  class: 'cm-active-line-dot-gutter',
  lineMarker(view, line) {
    const activeLine = view.state.doc.lineAt(view.state.selection.main.head);
    return line.from === activeLine.from ? activeLineDotMarker : null;
  },
  lineMarkerChange(update) {
    return update.selectionSet;
  },
  initialSpacer: () => activeLineDotMarker,
});


const modernTheme = EditorView.theme({
  '&': {
    backgroundColor: 'var(--color-bg-0)',
    color: 'var(--color-strong)',
    height: '100%',
  },
  '.cm-content': {
    caretColor: 'var(--color-strong)',
    paddingTop: '1.25rem',
    paddingLeft: '0.4rem'
  },
  '.cm-gutters': {
    backgroundColor: 'var(--color-bg-0)',
    color: 'var(--color-soft)',
    border: 'none',
    borderRight: '1px solid var(--color-line)',
  },
  '.cm-active-line-dot-gutter': {
    minWidth: '14px',
  },
  '.cm-active-line-dot-gutter .cm-gutterElement': {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

@Component({
  selector: 'app-code-editor',
  standalone: true,
  templateUrl: './code-editor.html',
  styleUrl: './code-editor.css',
})
export class CodeEditor {
  value = model('');
  language = input<'pseudocode' | 'go'>('pseudocode');
  private host = viewChild.required<ElementRef<HTMLDivElement>>('host');
  private view?: EditorView;

  constructor() {
    afterNextRender(() => this.initEditor());
    effect(() => {
      const newValue = this.value();
      if (this.view && newValue !== this.view.state.doc.toString()) {
        this.view.dispatch({
          changes: { from: 0, to: this.view.state.doc.length, insert: newValue },
        });
      }
    });
  }

  private initEditor() {
    const langExtension = this.language() === 'go' ? go() : pseudocodeLanguage;
    const smartEnter = keymap.of([
      {
        key: 'Enter',
        run: (view) => {
          const pos = view.state.selection.main.head;
          const line = view.state.doc.lineAt(pos);
          const newIndent = computeNewLineIndent(line.text);
          view.dispatch({
            changes: { from: pos, insert: '\n' + newIndent },
            selection: { anchor: pos + 1 + newIndent.length },
          });
          return true;
        },
      },
    ]);

    const state = EditorState.create({
      doc: this.value(),
      extensions: [
        langExtension,
        activeLineDotGutter,
        lineNumbers(),
        syntaxHighlighting(codeEditorHighlight),
        keymap.of([indentWithTab, ...defaultKeymap]),
        smartEnter,
        modernTheme,
        EditorView.updateListener.of((update) => {
          if (update.docChanged) this.value.set(update.state.doc.toString());
        }),
      ],
    });

    this.view = new EditorView({ state, parent: this.host().nativeElement });
  }
}
