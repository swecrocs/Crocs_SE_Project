import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs/internal/Observable";

@Injectable({providedIn: 'root'})

export class AccountService {
    private apiUrl = 'http://localhost:8080';

    constructor(private router: Router, private http: HttpClient){}

    register (data:any): Observable<any>{
        return this.http.post(`${this.apiUrl}/auth/register`, data);
    }
}