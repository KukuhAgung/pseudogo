import { Component, ElementRef, viewChild, afterNextRender, input } from '@angular/core';
import { EditorState } from '@codemirror/state';
import { EditorView } from '@codemirror/view';
import { syntaxHighlighting } from '@codemirror/language';
import { go } from '@codemirror/lang-go';
import { pseudocodeLanguage, fixedDarkHighlight } from '../../code-mirror/pseudocode-language';


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
        syntaxHighlighting(fixedDarkHighlight),
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
