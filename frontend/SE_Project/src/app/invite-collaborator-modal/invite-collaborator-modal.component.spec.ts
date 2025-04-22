import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { InviteCollaboratorModalComponent } from './invite-collaborator-modal.component';
import { ProjectService } from '../services/project.service';
import { of, throwError } from 'rxjs';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';

describe('InviteCollaboratorModalComponent', () => {
  let component: InviteCollaboratorModalComponent;
  let fixture: ComponentFixture<InviteCollaboratorModalComponent>;
  let mockProjectService: jasmine.SpyObj<ProjectService>;

  beforeEach(async () => {
    mockProjectService = jasmine.createSpyObj('ProjectService', ['inviteCollaborator']);

    await TestBed.configureTestingModule({
      imports: [InviteCollaboratorModalComponent, ReactiveFormsModule, MatButtonModule],
      providers: [
        { provide: ProjectService, useValue: mockProjectService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(InviteCollaboratorModalComponent);
    component = fixture.componentInstance;
    component.projectId = 103;

    // ngOnInit triggers here
    fixture.detectChanges();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
    expect(component.form).toBeDefined();
  });

  // Check required fields
  it('should return error if required fields are missing', () => {
    component.form.setValue({ email: '', role: '' });
    component.submit();
    expect(component.form.invalid).toBeTrue();
  });

  // Check email
  it('should return error if not a valid email', () => {
    component.form.setValue({ email: 'invalid-email', role: 'researcher' });
    component.submit();
    expect(component.form.get('email')?.invalid).toBeTrue();
  });

  // Test submit invite collaborator
  it('should call projectService and show success message on valid submission', fakeAsync(() => {
    mockProjectService.inviteCollaborator.and.returnValue(of({ message: 'Success!' }));
    spyOn(component.close, 'emit');

    component.projectId = 103;
    component.form.setValue({ email: 'test_1@email.com', role: 'programmer' });
    component.submit();

    // simulates async time
    tick(); 

    expect(mockProjectService.inviteCollaborator).toHaveBeenCalledWith(103, 'test_1@email.com', 'programmer');
    expect(component.success).toBe('Success!');

    // Simulate 1200ms setTimeout delay before emit
    tick(1200);
    expect(component.close.emit).toHaveBeenCalled();
  }));

  // Show error message if failed to send invitation
  it('should show error message when service fails', fakeAsync(() => {
    mockProjectService.inviteCollaborator.and.returnValue(throwError(() => new Error('Server error')));

    component.form.setValue({ email: 'test_2@email.com', role: 'programmer' });
    component.submit();

    tick(); 
    expect(component.error).toBe('Failed to send invitation.');
    expect(component.loading).toBeFalse();
  }));

});
