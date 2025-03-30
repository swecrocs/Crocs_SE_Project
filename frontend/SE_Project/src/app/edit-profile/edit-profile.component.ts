import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

import { HeaderComponent } from '../header/header.component';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, HeaderComponent, ReactiveFormsModule],
  templateUrl: './edit-profile.component.html',
  styleUrls: ['./edit-profile.component.css'],
})
export class ProfileComponent implements OnInit {
  profileForm!: FormGroup;

  userId = sessionStorage.getItem('userId') || '';

  successMessage = signal<string>('');
  errorMessage = signal<string>('');

  constructor(
    private fb: FormBuilder,
    private http: HttpClient,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.profileForm = this.fb.group({
      full_name: ['', Validators.required],
      affiliation: [''],
      bio: [''],
      email: [{ value: '', disabled: true }],
    });

    this.getProfile();
  }

  getProfile() {
    this.http
      .get<any>(`http://localhost:8080/users/${this.userId}/profile`)
      .subscribe({
        next: (data) => {
          this.profileForm.patchValue({
            full_name: data.full_name || '',
            affiliation: data.affiliation || '',
            bio: data.bio || '',
            email: data.email || '',
          });
        },
        error: (err) => {
          console.error('Error fetching profile:', err);
          this.errorMessage.set('Failed to load profile. Please try again.');
        },
      });
  }

  onSubmit() {
    if (this.profileForm.valid) {
      const body = {
        full_name: this.profileForm.value.full_name,
        affiliation: this.profileForm.value.affiliation,
        bio: this.profileForm.value.bio,
      };

      this.http
        .put<any>(`http://localhost:8080/users/${this.userId}/profile`, body)
        .subscribe({
          next: (res) => {
            console.log('Profile updated:', res);
            this.successMessage.set(res.message || 'Profile updated!');
            this.errorMessage.set('');
          },
          error: (err) => {
            console.error('Error updating profile:', err);
            this.errorMessage.set(
              'Failed to update profile. Please try again.'
            );
            this.successMessage.set('');
          },
        });
    } else {
      this.errorMessage.set('Please fill all required fields.');
    }
  }
}
