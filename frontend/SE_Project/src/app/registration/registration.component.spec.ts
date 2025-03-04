import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RegistrationComponent } from './registration.component';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { of, throwError } from 'rxjs';
import { AccountService } from '../account.service';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';


describe ('RegistrationComponent', () => {
  let component: RegistrationComponent;
  let fixture: ComponentFixture<RegistrationComponent>;
  let accountService: AccountService;
  let snackBar: MatSnackBar;

  beforeEach (() => {
      TestBed.configureTestingModule({
          imports: [
              ReactiveFormsModule,
              RegistrationComponent
          ],
          providers: [
              { provide: AccountService, useClass: MockAccountService },
              MatSnackBar,
              provideAnimationsAsync()
          ]
      });

      fixture = TestBed.createComponent(RegistrationComponent);
      component = fixture.componentInstance;
      accountService = TestBed.inject(AccountService);
      snackBar = TestBed.inject(MatSnackBar);
  });

  class MockAccountService {
      register(data: any) {
          return of({ message: 'Success' });
      }
  }

  // component initialize
  it ('should create the component and initialize form controls', () => {
      fixture.detectChanges();
      expect(component).toBeTruthy();
      expect(component.form).toBeDefined();
      expect(component.form.controls['email']).toBeTruthy();
      expect(component.form.controls['password']).toBeTruthy();
      expect(component.form.controls['confirmPwd']).toBeTruthy();
  });

  it ('should register successfully when passwords match', () => {
      spyOn(accountService, 'register').and.returnValue(of({ message: 'Success' }));
      spyOn(snackBar, 'open');

      component.form.setValue({
          email: 'test@email.com',
          password: 'test01234',
          confirmPwd: 'test01234'
      });

      component.getFormValue();

      expect(accountService.register).toHaveBeenCalledWith({ email: 'test@email.com', password: 'test01234' });
      expect(snackBar.open).toHaveBeenCalledWith('Successfully register', 'close', jasmine.objectContaining({
          duration: 5000
      }));
  });

  // show popup error message when passwords do not match
  it('should show error when passwords do not match', () => {
      spyOn(snackBar, 'open');

      component.form.setValue({
          email: 'test@email.com',
          password: 'test01234',
          confirmPwd: 'test43210'
      });

      component.getFormValue();

      expect(snackBar.open).toHaveBeenCalledWith('Passwords do not match', 'close', jasmine.objectContaining({
          duration: 5000
      }));
  });

  it ('should return true when passwords match', () => {
      const result = component.checkConfirmPwd({ password: 'test01234', confirmPwd: 'test01234' });
      expect(result).toBeTrue();
  });

  it ('should return false when passwords do not match', () => {
      const result = component.checkConfirmPwd({ password: 'test01234', confirmPwd: 'test43210' });
      expect(result).toBeFalse();
  });

  // show popup error message when registration fails
  it ('should show error message when registration fails', () => {
      spyOn(accountService, 'register').and.returnValue(throwError({ message: 'Error occurred' }));
      spyOn(snackBar, 'open');

      component.form.setValue({
          email: 'test@email.com',
          password: 'test01234',
          confirmPwd: 'test01234'
      });

      component.getFormValue();

      expect(snackBar.open).toHaveBeenCalledWith('Failed to register:Error occurred', 'close', jasmine.objectContaining({
          duration: 5000
      }));
  });

  // check if the login button can navigate to login page correctly
  it ('should navigate to login page when onLogin is called', () => {
      const router = TestBed.inject(Router);
      spyOn(router, 'navigate');

      component.onLogin();

      expect(router.navigate).toHaveBeenCalledWith(['/login']);
  });
});
