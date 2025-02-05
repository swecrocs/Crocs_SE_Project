import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-registration',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
})

export class RegistrationComponent {

    constructor(private router: Router){}
    
    form = new FormGroup({
        email: new FormControl(),
        password: new FormControl(),
        confirmPwd: new FormControl()
    })

    getFormValue () {
        const value = this.form?.value;
        return value;
    }

    onLogin () {
        this.router.navigate(['/login']);
    }

}