import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class UiStateService {
  isDarkMode = signal(false);
  sidebarOpen = signal(false);
  sidebarDocsOpen = signal(false);
  hasUnreadLog = signal(false);

  toggleDarkMode() {
    this.isDarkMode.update((val) => !val);
    if (this.isDarkMode()) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }

  toggleSidebar() {
    this.sidebarOpen.update((open) => !open);
    if (this.sidebarOpen()) {
      this.hasUnreadLog.set(false);
    }
  }

  toggleSidebarDocs() {
    this.sidebarDocsOpen.update((open) => !open);
  }

  flagUnreadIfClosed() {
    if (!this.sidebarOpen()) {
      this.hasUnreadLog.set(true);
    }
  }
}
