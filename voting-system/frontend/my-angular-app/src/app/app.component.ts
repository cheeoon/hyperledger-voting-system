import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';  // Import HttpClient and HttpHeaders

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  response: any; // To store the response from the server
  errorMessage: string | null = null; // To store any error message

  constructor(private http: HttpClient) {}

  submitData() {
    const url = 'http://localhost:3000/invoke?channelid=mychannel&chaincodeid=votingcc&function=InitElection&args=%5B%22Alice%22%2C%20%22Bob%22%2C%20%22Charlie%22%5D&args=2024-12-07T10%3A00%3A00Z&args=2024-12-08T10%3A00%3A00Z'
  
    const headers = new HttpHeaders({
      'Content-Type': 'application/x-www-form-urlencoded',
    });
  
    const body = new URLSearchParams();
    body.set('channelid', 'mychannel');
    body.set('chaincodeid', 'basic');
    body.set('function', 'InitElection');
    body.append('args', JSON.stringify(["Alicddde", "Bodddb", "Cdddharlie"])); // Correctly stringify args
    body.append('args', '2024-12-07T10:00:00Z');
    body.append('args', '2024-12-08T10:00:00Z');
  
    this.http.post(url, { headers }).subscribe({
      next: (res) => {
        console.log('Response:', res); // Log the response for debugging
        if (res) {
          console.log('fgdfg')
          this.response = res;
          this.errorMessage = null;
        } else {
          console.log('sgdfg')
          this.errorMessage = 'Empty response received';
        }
      },
      error: (err) => {
        console.log('sgdfg')
        console.error('Error:', err); // Log the error for debugging
        this.response = null;
      }
    });
  }
}
