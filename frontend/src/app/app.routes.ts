import { Routes } from '@angular/router';
import { ConverterPage } from './features/converter-page/converter-page';
import { DocsPage } from './features/docs-page/docs-page';

export const routes: Routes = [
    { path: '', component: ConverterPage },
    { path: 'docs', component: DocsPage }
];
