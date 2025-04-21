import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { InvitationsModalComponent } from '../invitations-modal/invitations-modal.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, InvitationsModalComponent],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css'],
})
export class HeaderComponent {
  // controls the profile dropdown
  dropdownOpen = signal(false);

  // controls whether the invitations modal is shown
  showInvitationsModal = signal(false);

  constructor(private router: Router) {}

  isLoggedIn(): boolean {
    return !!sessionStorage.getItem('token');
  }

  goHome() {
    this.router.navigate(['/']);
  }

  goToLogin() {
    this.router.navigate(['/login']);
  }

  goToSignup() {
    this.router.navigate(['/registration']);
  }

  goToProfile() {
    this.router.navigate(['/edit-profile']);
  }

  goToProjects() {
    if (!this.isLoggedIn()) {
      this.router.navigate(['/login']);
    } else {
      this.router.navigate(['/projects']);
    }
  }

  // called when “Invitations” is clicked in the nav
  goToInvitations() {
    if (!this.isLoggedIn()) {
      this.router.navigate(['/login']);
    } else {
      this.showInvitationsModal.set(true);
    }
  }

  // hide the invitations modal
  hideInvitationsModal() {
    this.showInvitationsModal.set(false);
  }

  logout() {
    sessionStorage.clear();
    this.dropdownOpen.set(false);
    this.router.navigate(['/']);
  }

  toggleDropdown() {
    this.dropdownOpen.set(!this.dropdownOpen());
  }
}
