// main.ts
import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { HttpRequest, HttpHandlerFn, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';

function authInterceptor() {
  return (
    req: HttpRequest<any>,
    next: HttpHandlerFn
  ): Observable<HttpEvent<any>> => {
    const token = sessionStorage.getItem('token');
    if (token) {
      const authReq = req.clone({
        setHeaders: { Authorization: `Bearer ${token}` },
      });
      return next(authReq);
    }
    return next(req);
  };
}

bootstrapApplication(AppComponent, {
  ...appConfig,
  providers: [
    ...(appConfig.providers || []),
    provideHttpClient(withInterceptors([authInterceptor()])),
  ],
}).catch((err) => console.error(err));
