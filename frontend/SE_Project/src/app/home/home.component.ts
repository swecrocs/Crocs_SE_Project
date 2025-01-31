import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgOptimizedImage } from '@angular/common';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, NgOptimizedImage],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent {
  showLoginModal = false;
  showSignupModal = false;

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.fetchBackendMessage();
  }

  fetchBackendMessage() {
    this.http.get<{ message: string }>('/api').subscribe({
      next: (response) => {
        console.log("Backend Message:", response.message); // Logs backend message in browser console
      },
      error: (error) => {
        console.error("Error fetching backend message:", error);
      },
      complete: () => {
        console.log("Backend request completed.");
      }
    });
  }

  openLoginModal() {
    this.showLoginModal = true;
  }

  closeLoginModal() {
    this.showLoginModal = false;
  }

  openSignupModal() {
    this.showSignupModal = true;
  }

  closeSignupModal() {
    this.showSignupModal = false;
  }
}
