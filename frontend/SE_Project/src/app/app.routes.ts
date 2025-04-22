import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { RegistrationComponent } from './registration/registration.component';
import { ProfileComponent } from './edit-profile/edit-profile.component';
import { ProjectsHomeComponent } from './projects-home/projects-home.component';
import { ProjectEditorComponent } from './edit-project/projectEditor.component';
import { AboutPageComponent } from './about-page/about-page.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'registration', component: RegistrationComponent },
  { path: 'edit-profile', component: ProfileComponent },
  { path: 'projects', component: ProjectsHomeComponent },
  { path: 'about', component: AboutPageComponent},
  { path: 'projects/:id', component: ProjectEditorComponent },
  { path: '**', redirectTo: '', pathMatch: 'full' },
];
