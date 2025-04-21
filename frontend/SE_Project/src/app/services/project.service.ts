// src/app/services/project.service.ts

import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, catchError, of } from 'rxjs';

export interface Project {
  id?: number;
  title: string;
  description: string;
  status: string;
  visibility: string;
  required_skills: string[];
  owner_id?: number;
}

export interface Invitation {
  id: number;
  project_id: number;
  project_title: string;
  inviter_name: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
}

@Injectable({
  providedIn: 'root',
})
export class ProjectService {
  private baseUrl = 'http://localhost:8080/projects';

  constructor(private http: HttpClient) {}

  getUserProjects(): Observable<{ projects: Project[] }> {
    return this.http
      .get<{ projects: Project[] }>(`${this.baseUrl}/user`, {
        headers: this.getAuthHeaders(),
      })
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] getUserProjects error:', error);
          return of({ projects: [] });
        })
      );
  }

  getAllProjects(): Observable<{ projects: Project[] }> {
    return this.http
      .get<{ projects: Project[] }>(`${this.baseUrl}`, {
        headers: this.getAuthHeaders(),
      })
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] getAllProjects error:', error);
          return of({ projects: [] });
        })
      );
  }

  getProjectById(id: number): Observable<Project | null> {
    return this.http
      .get<Project>(`${this.baseUrl}/${id}`, {
        headers: this.getAuthHeaders(),
      })
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] getProjectById error:', error);
          return of(null);
        })
      );
  }

  createProject(projectData: Project): Observable<any> {
    return this.http
      .post<any>(`${this.baseUrl}`, projectData, {
        headers: this.getAuthHeaders(),
      })
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] createProject error:', error);
          return of({ error: 'Failed to create project' });
        })
      );
  }

  updateProject(id: number, project: Project): Observable<Project> {
    return this.http.put<Project>(`${this.baseUrl}/${id}`, project, {
      headers: this.getAuthHeaders(),
    });
  }

  getInvitations(): Observable<{ invitations: Invitation[] }> {
    return this.http
      .get<{ invitations: Invitation[] }>(`${this.baseUrl}/invitations`, {
        headers: this.getAuthHeaders(),
      })
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] getInvitations error:', error);
          return of({ invitations: [] });
        })
      );
  }

  inviteCollaborator(
    projectId: number,
    email: string,
    role: string
  ): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(
        `${this.baseUrl}/${projectId}/collaborators`,
        { email, role },
        { headers: this.getAuthHeaders() }
      )
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] inviteCollaborator error:', error);
          return of({ message: 'Failed to send invitation' });
        })
      );
  }

  respondToInvitation(
    projectId: number,
    invitationId: number,
    action: 'accept' | 'reject'
  ): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(
        `${this.baseUrl}/${projectId}/collaborators/invitations/${invitationId}/${action}`,
        {},
        { headers: this.getAuthHeaders() }
      )
      .pipe(
        catchError((error) => {
          console.error('[ProjectService] respondToInvitation error:', error);
          return of({ message: 'Failed to respond to invitation' });
        })
      );
  }

  private getAuthHeaders(): HttpHeaders {
    const token = sessionStorage.getItem('token') || '';
    return new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
  }
}
