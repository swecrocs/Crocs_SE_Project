import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, NgForm, ReactiveFormsModule, Validators } from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {MatSnackBar, MatSnackBarVerticalPosition} from '@angular/material/snack-bar';
import { HeaderComponent } from "../header/header.component";
import {MatChipInputEvent, MatChipsModule} from '@angular/material/chips';
import {COMMA, ENTER} from '@angular/cdk/keycodes';
import {MatSelectModule} from '@angular/material/select';
import { Project, ProjectService } from '../services/project.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-project_editor',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule,MatFormFieldModule, 
    MatIconModule, MatInputModule, HeaderComponent, MatChipsModule, MatSelectModule],
  templateUrl: './projectEditor.component.html',
  styleUrls: ['./projectEditor.component.css'],
})


export class ProjectEditorComponent {
    private registerPrompt = inject(MatSnackBar);
    verticalPosition: MatSnackBarVerticalPosition = 'top';

    projectForm!: FormGroup;
    project: Project | null = null;
    projectId: number = 0;

    readonly addOnBlur = true;
    readonly separatorKeysCodes = [ENTER, COMMA] as const;
    skills = signal<string[]>([]);

    statusOptions = ['open', 'closed'];
    visibilityOptions = ['private', 'public'];

    constructor(
        private fb: FormBuilder,
        private projectService: ProjectService,
        private route: ActivatedRoute
    ){}

    ngOnInit() {
        this.projectForm = this.fb.group({
            title: ['', Validators.required],
            description: [''],
            status: ['', Validators.required],
            visibility: ['', Validators.required],
            required_skills: [''],
        });
        this.loadProject();
    }

    loadProject() {
        this.route.paramMap.subscribe(params => {
            const idParam = params.get('id');
            this.projectId = idParam ? +idParam : 0;
      
            if (this.projectId) {
              this.projectService.getProjectById(this.projectId).subscribe(data => {
                if (data !== null) {
                    this.projectForm.patchValue(data);
                    if (data.required_skills.length > 0)
                    this.skills.set(data.required_skills);
                  } else {
                    this.handlePromptOpen('Project not found', 'close');
                  }
              });
            }
        });
    }

    handleSave () {
        if (this.projectForm.invalid) {
            this.handlePromptOpen('Please fill out required fields.', 'close');
            return;
        }

        // Get the form values
        const formValue = this.projectForm.value;

        // Add project id to form value
        const updatedProject: Project = {
            ...formValue,
            id: this.projectId
        };

        this.projectService.updateProject(this.projectId, updatedProject).subscribe({
            next: () => {
                this.handlePromptOpen('Project saved successfully!', 'OK');
            },
            error: (err) => {
                console.error('Failed to save project:', err);
                this.handlePromptOpen('Failed to save project.', 'close');
            }
        });
    }

    handleAddSkill(event: MatChipInputEvent): void {
        const value = (event.value || '').trim();

        // Add skill
        if (value) {
            this.skills.update((skills) => [...skills, value]);
            // Update form
            this.projectForm.get('required_skills')?.setValue([...this.skills()]);
        }
    
        // Clear the input value
        event.chipInput!.clear();
    }

    handleDeleteSkill (skill: string) {
        this.skills.update(skills => {
            const index = skills.indexOf(skill);
            if (index < 0) {
              return skills;
            }
      
            skills.splice(index, 1);
            return [...skills];
        });
        // Update form
        this.projectForm.get('required_skills')?.setValue([...this.skills()]);
    }

    // show feedback prompt
    handlePromptOpen (message: string, action: string) {
        this.registerPrompt.open(message, action, {
            verticalPosition: this.verticalPosition,
            duration: 5000
        });
    }

    isNotEmptyString (str: string) {
        return str && str.length > 0;
    }

}