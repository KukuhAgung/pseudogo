import { Component, inject } from '@angular/core';
import { UiStateService } from '../../../../core/services/ui-state';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { remixCloseFill } from '@ng-icons/remixicon'

@Component({
  imports: [RouterLink, RouterLinkActive, NgIcon],
  providers: [provideIcons({ remixCloseFill })],
  selector: 'app-docs-sidebar',
  standalone: true,
  styleUrl: './docs-sidebar.css',
  templateUrl: './docs-sidebar.html',
})
export class DocsSidebar {
  ui = inject(UiStateService);
}
