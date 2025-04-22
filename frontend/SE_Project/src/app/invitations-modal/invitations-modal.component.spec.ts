import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { InvitationsModalComponent } from './invitations-modal.component';
import { ProjectService, Invitation } from '../services/project.service';
import { of } from 'rxjs';
import { MatButtonModule } from '@angular/material/button';

describe('InvitationsModalComponent', () => {
  let component: InvitationsModalComponent;
  let fixture: ComponentFixture<InvitationsModalComponent>;
  let mockProjectService: jasmine.SpyObj<ProjectService>;

  const mockInvitations: Invitation[] = [
    { 
        id: 1, 
        project_id: 101, 
        project_title: 'Project A',
        inviter_name: 'Researcher_A',
        email: 'researcherA@email.com',
        role: 'researcher',
        status: 'pending',
        created_at: '2024-04-20T12:00:00Z',
    },
    { 
        id: 2, 
        project_id: 102, 
        project_title: 'Project B',
        inviter_name: 'Researcher_B',
        email: 'researcherB@email.com',
        role: 'programmer',
        status: 'pending',
        created_at: '2024-04-21T12:00:00Z',
    }
  ];

  beforeEach(async () => {
    mockProjectService = jasmine.createSpyObj('ProjectService', ['getInvitations', 'respondToInvitation']);

    await TestBed.configureTestingModule({
      imports: [InvitationsModalComponent, MatButtonModule],
      providers: [{ provide: ProjectService, useValue: mockProjectService }]
    }).compileComponents();

    fixture = TestBed.createComponent(InvitationsModalComponent);
    component = fixture.componentInstance;
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });

  it('should load invitations on init', fakeAsync(() => {
    mockProjectService.getInvitations.and.returnValue(of({ invitations: mockInvitations }));

    // triggers ngOnInit
    fixture.detectChanges(); 
    tick();

    expect(component.invitations()).toEqual(mockInvitations);
    expect(component.loading()).toBeFalse();
    expect(component.error()).toBeNull();
  }));

  // If invitations are empty show no pending invitations
  it('should handle empty invitations', fakeAsync(() => {
    mockProjectService.getInvitations.and.returnValue(of({ invitations: [] }));

    fixture.detectChanges();
    tick();

    expect(component.invitations()).toEqual([]);
    expect(component.error()).toBe('No pending invitations.');
  }));

  // Test accepting an invitation
  it('should accept an invitation and reload', fakeAsync(() => {
    mockProjectService.getInvitations.and.returnValue(of({ invitations: mockInvitations }));
    mockProjectService.respondToInvitation.and.returnValue(of({status: 'success', message: 'Invitation accepted.'}));

    fixture.detectChanges();
    tick();

    component.accept(mockInvitations[0]);
    tick();

    expect(mockProjectService.respondToInvitation).toHaveBeenCalledWith(101, 1, 'accept');
  }));

  // Test rejecting an invitation 
  it('should reject an invitation and reload', fakeAsync(() => {
    mockProjectService.getInvitations.and.returnValue(of({ invitations: mockInvitations }));
    mockProjectService.respondToInvitation.and.returnValue(of({status: 'reject', message: 'Invitation rejected.'}));

    fixture.detectChanges();
    tick();

    component.reject(mockInvitations[1]);
    tick();

    expect(mockProjectService.respondToInvitation).toHaveBeenCalledWith(102, 2, 'reject');
  }));

});