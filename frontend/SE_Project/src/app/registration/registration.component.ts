import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import { Router } from '@angular/router';
import { AccountService } from '../account.service';
import {MatSnackBar, MatSnackBarVerticalPosition} from '@angular/material/snack-bar';


@Component({
  selector: 'app-registration',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule, 
            MatFormFieldModule, MatIconModule, MatInputModule],
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
})

export class RegistrationComponent {
    hide = signal(true);
    private registerPrompt = inject(MatSnackBar);
    verticalPosition: MatSnackBarVerticalPosition = 'top';

    constructor(
        private router: Router,
        private accountService: AccountService
    ){}
    
    form = new FormGroup({
        email: new FormControl('', [Validators.required, Validators.email]),
        password: new FormControl('', [Validators.required]),
        confirmPwd: new FormControl('', [Validators.required])
    })

    getFormValue () {
        const value = this.form?.value;
        this.handleRegister(value);
    }

    // new register
    handleRegister (data: any) {
        if (this.checkConfirmPwd(data)){
            const newValue = { email: data?.email, password: data?.password };
            this.accountService.register(newValue).subscribe({
                next: (response) => {
                    console.log('Successfully register:', response);
                    this.handlePromptOpen('Successfully register', 'close');
                },
                error: (err) => {
                    console.log(err);
                    const {message} = err;
                    console.error('Failed to register:', err);
                    this.handlePromptOpen(`Failed to register:${message}`, 'close');
                }
            });
        } else {
            console.log('Passwords do not match');
            this.handlePromptOpen('Passwords do not match', 'close');
        }
    }

    // validateFormValue (data: any) {
    //     if (!this.isNotEmptyString(data?.email)) 
    // }

    // check confirm password
    checkConfirmPwd (data: any) {
        const { password, confirmPwd } = data;
        return password === confirmPwd; 
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

    onHide (event: MouseEvent) {
        event.stopPropagation();
        this.hide.set(!this.hide());
    }

    onLogin () {
        this.router.navigate(['/login']);
    }

}