import { Component, ElementRef, viewChild, afterNextRender, model, effect } from '@angular/core';
import { EditorState } from '@codemirror/state';
import { EditorView, lineNumbers, keymap } from '@codemirror/view';
import { defaultKeymap, indentWithTab } from '@codemirror/commands';

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

@Component({
  selector: 'app-code-editor',
  standalone: true,
  templateUrl: './code-editor.html',
  styleUrl: './code-editor.css',
})
export class CodeEditor {
  value = model('');
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
        lineNumbers(),
        keymap.of([indentWithTab, ...defaultKeymap]),
        smartEnter,
        EditorView.updateListener.of((update) => {
          if (update.docChanged) this.value.set(update.state.doc.toString());
        }),
        EditorView.theme({ '&': { fontSize: '14px', height: '100%' } }),
      ],
    });

    this.view = new EditorView({ state, parent: this.host().nativeElement });
  }
}
