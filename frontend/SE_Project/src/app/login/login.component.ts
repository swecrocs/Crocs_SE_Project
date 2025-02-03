import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})

export class LoginComponent {
    
    form = new FormGroup({
        username: new FormControl(),
        password: new FormControl()
    })

    getFormValue () {
        const value = this.form?.value;
        return value;
    }

    onRegister() {

    }
}