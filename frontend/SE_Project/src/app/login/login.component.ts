import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})

export class LoginComponent {
    
    constructor(private router: Router){}

    form = new FormGroup({
        email: new FormControl(),
        password: new FormControl()
    })

    getFormValue () {
        const value = this.form?.value;
        return value;
    }

    onRegister() {
        this.router.navigate(['/registration']);
    }
}