import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from '../header/header.component';

@Component({
  selector: 'app-home',
  standalone: true,
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
  imports: [CommonModule, HeaderComponent],
})
export class HomeComponent {
  isLoggedIn = signal<boolean>(!!localStorage.getItem('token'));
  dropdownOpen = signal<boolean>(false);

  constructor(private router: Router) {}

  goToLogin() {
    this.router.navigate(['/login']);
  }

  goToSignup() {
    this.router.navigate(['/registration']);
  }

  goToProfile() {
    this.router.navigate(['/profile']);
  }

  logout() {
    sessionStorage.removeItem('token');
    this.isLoggedIn.set(false);
    this.dropdownOpen.set(false);
    this.router.navigate(['/']);
  }

  toggleDropdown() {
    this.dropdownOpen.set(!this.dropdownOpen());
  }
}
