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
  dropdownOpen = signal<boolean>(false);

  constructor(private router: Router) {}

  public isLoggedIn(): boolean {
    return !!sessionStorage.getItem('token');
  }

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

  goToProjects() {
    if (!this.isLoggedIn()) {
      this.router.navigate(['/login']);
    } else {
      this.router.navigate(['/projects']);
    }
  }

  logout() {
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('userId');
    this.dropdownOpen.set(false);
    this.router.navigate(['/']);
  }

  toggleDropdown() {
    this.dropdownOpen.set(!this.dropdownOpen());
  }
}
