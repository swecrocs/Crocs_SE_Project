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

@Injectable({
  providedIn: 'root',
})
export class ProjectService {
  private baseUrl = 'http://localhost:8080/projects';

  constructor(private http: HttpClient) {}

  getAllProjects(): Observable<Project[]> {
    const headers = this.getAuthHeaders();
    return this.http.get<Project[]>(this.baseUrl, { headers }).pipe(
      catchError((error) => {
        console.error('[ProjectService] getAllProjects error:', error);
        return of([]);
      })
    );
  }

  getProjectById(id: number): Observable<Project | null> {
    const headers = this.getAuthHeaders();
    return this.http.get<Project>(`${this.baseUrl}/${id}`, { headers }).pipe(
      catchError((error) => {
        console.error('[ProjectService] getProjectById error:', error);
        return of(null);
      })
    );
  }

  createProject(projectData: Project): Observable<any> {
    const headers = this.getAuthHeaders();
    return this.http.post<any>(this.baseUrl, projectData, { headers }).pipe(
      catchError((error) => {
        console.error('[ProjectService] createProject error:', error);
        return of({ error: 'Failed to create project' });
      })
    );
  }

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('token') || '';
    return new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });
  }
}
