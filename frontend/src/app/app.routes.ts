import { Routes } from '@angular/router';
import { ConverterPage } from './features/converter-page/converter-page';
import { DocsPage } from './features/docs-page/docs-page';
import { Introduction } from './features/docs-page/pages/introduction/introduction';
import { Variable } from './features/docs-page/pages/variable/variable';
import { Io } from './features/docs-page/pages/io/io';
import { Operator } from './features/docs-page/pages/operator/operator';
import { ConditionalStatements } from './features/docs-page/pages/conditional-statements/conditional-statements';
import { Looping } from './features/docs-page/pages/looping/looping';

export const routes: Routes = [
  { path: '', component: ConverterPage },
  {
    path: 'docs',
    component: DocsPage,
    children: [
      { path: '', redirectTo: 'introduction', pathMatch: 'full' },
      { path: 'introduction', component: Introduction },
      { path: 'variable', component: Variable },
      { path: 'input-output', component: Io },
      { path: 'operators', component: Operator },
      { path: 'conditional-statements', component: ConditionalStatements },
      { path: 'looping', component: Looping },
    ],
  },
];
