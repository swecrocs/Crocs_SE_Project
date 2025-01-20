import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  standalone: true,
  template: `<h1>{{ message }}</h1>`,
  styles: [],
})
export class AppComponent {
  message: string = 'Loading...';

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    // Fetch data from the Golang backend
    this.http.get<{ message: string }>('/api').subscribe({
      next: (response) => {
        this.message = response.message;
      },
      error: (error) => {
        this.message = 'Error connecting to backend.';
      },
      complete: () => {
        console.log('Request completed.');
      }
    });
  }
}
