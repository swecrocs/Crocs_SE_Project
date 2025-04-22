import { Component } from "@angular/core";
import { HeaderComponent } from "../header/header.component";
import {MatIconModule} from '@angular/material/icon';

@Component({
  selector: 'app-about-page',
  standalone: true,
  imports: [HeaderComponent, MatIconModule],
  templateUrl: './about-page.component.html',
  styleUrls: ['./about-page.component.css'],
})

export class AboutPageComponent {
    
}