import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { ProjectEditorComponent } from './projectEditor.component';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { ProjectService, Project } from '../services/project.service';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { ReactiveFormsModule } from '@angular/forms';

describe('ProjectEditorComponent', () => {
  let component: ProjectEditorComponent;
  let fixture: ComponentFixture<ProjectEditorComponent>;
  let mockProjectService: jasmine.SpyObj<ProjectService>;

  const fakeProject: Project = {
    id: 1,
    title: 'Test Project',
    description: 'Test Description',
    status: 'open',
    visibility: 'public',
    required_skills: ['Angular', 'TypeScript']
  };

  beforeEach(async () => {
    mockProjectService = jasmine.createSpyObj<ProjectService>('ProjectService', ['getProjectById']);

    await TestBed.configureTestingModule({
      imports: [
        ProjectEditorComponent,
        ReactiveFormsModule,
        MatSnackBarModule
      ],
      providers: [
        {
          provide: ProjectService,
          useValue: mockProjectService
        },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of({
              get: (key: string) => key === 'id' ? '1' : null
            })
          }
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ProjectEditorComponent);
    component = fixture.componentInstance;
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize form and load project', fakeAsync(() => {
    mockProjectService.getProjectById.and.returnValue(of(fakeProject));
    fixture.detectChanges();
    tick();

    expect(component.projectForm).toBeTruthy();
    expect(component.projectForm.get('title')?.value).toBe(fakeProject.title);
    expect(component.skills()).toEqual(fakeProject.required_skills);
  }));

  it('should add a skill to the form', () => {
    component.projectForm = component['fb'].group({
      title: [''],
      description: [''],
      status: [''],
      visibility: [''],
      required_skills: [''],
      skills: ['']
    });

    component.skills.set(['HTML']);
    const event = { value: 'CSS', chipInput: { clear: () => {} } } as any;

    component.handleAddSkill(event);
    expect(component.skills()).toContain('CSS');
  });

  it('should delete a skill from the form', () => {
    component.skills.set(['HTML', 'CSS']);
    component.projectForm = component['fb'].group({
      skills: ['']
    });

    component.handleDeleteSkill('CSS');
    expect(component.skills()).toEqual(['HTML']);
  });
});
