import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { Project, ProjectService } from '../services/project.service';
import { HeaderComponent } from '../header/header.component';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';

@Component({
  selector: 'app-projects-home',
  standalone: true,
  imports: [CommonModule, HeaderComponent, ReactiveFormsModule],
  templateUrl: './projects-home.component.html',
  styleUrls: ['./projects-home.component.css'],
})
export class ProjectsHomeComponent implements OnInit {
  projects = signal<Project[]>([]);
  loading = signal<boolean>(true);
  errorMessage = signal<string | null>(null);
  showCreateModal = signal<boolean>(false);
  createForm!: FormGroup;

  constructor(
    private projectService: ProjectService,
    private router: Router,
    private fb: FormBuilder
  ) {}

  ngOnInit() {
    this.loadProjects();

    this.createForm = this.fb.group({
      title: ['', Validators.required],
      description: [''],
      status: ['open', Validators.required],
      visibility: ['private', Validators.required],
      required_skills: [''],
    });
  }

  loadProjects() {
    this.loading.set(true);
    this.projectService.getAllProjects().subscribe((response) => {
      this.loading.set(false);
      const data = response.projects || [];

      if (data.length > 0) {
        this.projects.set(data);
        this.errorMessage.set(null);
      } else {
        this.projects.set([]);
        this.errorMessage.set('No projects found.');
      }
    });
  }

  onCreateProject() {
    this.showCreateModal.set(true);
  }

  closeModal() {
    this.showCreateModal.set(false);
    this.createForm.reset({
      title: '',
      description: '',
      status: 'open',
      visibility: 'private',
      required_skills: '',
    });
  }

  submitCreateForm() {
    if (this.createForm.valid) {
      const formValue = this.createForm.value;

      const skillsArray = formValue.required_skills
        ? formValue.required_skills.split(',').map((s: string) => s.trim())
        : [];

      const newProject: Project = {
        title: formValue.title,
        description: formValue.description || '',
        status: formValue.status,
        visibility: formValue.visibility,
        required_skills: skillsArray,
      };

      this.projectService.createProject(newProject).subscribe((response) => {
        if (response?.message) {
          alert('Project created successfully!');
          this.closeModal();
          this.loadProjects();
        } else {
          alert('Failed to create project.');
        }
      });
    } else {
      this.createForm.markAllAsTouched();
    }
  }

  openProject(id: number | undefined) {
    if (!id) return;
    this.router.navigate(['/projects', id]);
  }
}
