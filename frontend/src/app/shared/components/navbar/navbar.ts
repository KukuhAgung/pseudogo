import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { UiStateService } from '../../../core/services/ui-state';
import { remixUploadLine, remixBookOpenLine, remixSunLine, remixMoonLine, remixTerminalFill } from '@ng-icons/remixicon';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [NgIcon, RouterLink, RouterLinkActive],
  providers: [provideIcons({ remixUploadLine, remixBookOpenLine, remixSunLine, remixMoonLine, remixTerminalFill })],
  templateUrl: './navbar.html',
})
export class NavbarComponent {
  ui = inject(UiStateService);
}
