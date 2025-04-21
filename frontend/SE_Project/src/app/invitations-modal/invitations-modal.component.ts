import { Component, OnInit, Output, EventEmitter, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { ProjectService, Invitation } from '../services/project.service';

@Component({
  selector: 'app-invitations-modal',
  standalone: true,
  imports: [CommonModule, MatButtonModule],
  templateUrl: './invitations-modal.component.html',
  styleUrls: ['./invitations-modal.component.css'],
})
export class InvitationsModalComponent implements OnInit {
  @Output() close = new EventEmitter<void>();

  invitations = signal<Invitation[]>([]);
  loading = signal(true);
  error = signal<string | null>(null);

  constructor(private projectService: ProjectService) {}

  ngOnInit() {
    this.loadInvitations();
  }

  loadInvitations() {
    this.loading.set(true);
    this.projectService.getInvitations().subscribe((resp) => {
      this.loading.set(false);
      const invs = resp.invitations || [];
      this.invitations.set(invs);
      this.error.set(invs.length ? null : 'No pending invitations.');
    });
  }

  accept(inv: Invitation) {
    this.projectService
      .respondToInvitation(inv.project_id, inv.id, 'accept')
      .subscribe(() => this.loadInvitations());
  }

  reject(inv: Invitation) {
    this.projectService
      .respondToInvitation(inv.project_id, inv.id, 'reject')
      .subscribe(() => this.loadInvitations());
  }

  onClose() {
    this.close.emit();
  }
}
