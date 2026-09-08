import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  imports: [RouterLink, RouterLinkActive],
  selector: 'app-docs-sidebar',
  styleUrl: './docs-sidebar.css',
  templateUrl: './docs-sidebar.html',
})
export class DocsSidebar {}
