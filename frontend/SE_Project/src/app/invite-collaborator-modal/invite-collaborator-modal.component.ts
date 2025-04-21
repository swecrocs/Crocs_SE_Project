import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { ProjectService } from '../services/project.service';

@Component({
  selector: 'app-invite-collaborator-modal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule],
  templateUrl: './invite-collaborator-modal.component.html',
  styleUrls: ['./invite-collaborator-modal.component.css'],
})
export class InviteCollaboratorModalComponent implements OnInit {
  @Input() projectId!: number;
  @Output() close = new EventEmitter<void>();

  form!: FormGroup; // declare without initializer
  loading = false;
  error: string | null = null;
  success: string | null = null;

  constructor(
    private fb: FormBuilder,
    private projectService: ProjectService
  ) {}

  ngOnInit() {
    // now fb is ready—initialize the form here
    this.form = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      role: ['programmer', Validators.required],
    });
  }

  submit() {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.loading = true;
    this.error = this.success = null;

    const { email, role } = this.form.value;
    this.projectService
      .inviteCollaborator(this.projectId, email!, role!)
      .subscribe({
        next: (res) => {
          this.loading = false;
          this.success = res.message || 'Invitation sent!';
          setTimeout(() => this.close.emit(), 1200);
        },
        error: () => {
          this.loading = false;
          this.error = 'Failed to send invitation.';
        },
      });
  }

  onClose() {
    this.close.emit();
  }
}
