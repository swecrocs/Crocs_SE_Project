import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-header',
  standalone: true,
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css'],
  imports: [CommonModule],
})
export class HeaderComponent {
  // Reuse signals or services to track login state
  isLoggedIn = signal<boolean>(!!localStorage.getItem('token'));
  dropdownOpen = signal<boolean>(false);

  constructor(private router: Router) {}

  goHome() {
    this.router.navigate(['/home']);
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

  logout() {
    localStorage.removeItem('token');
    this.isLoggedIn.set(false);
    this.dropdownOpen.set(false);
    this.router.navigate(['/']);
  }

  toggleDropdown() {
    this.dropdownOpen.set(!this.dropdownOpen());
  }
}
