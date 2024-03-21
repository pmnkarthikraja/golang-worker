# golang-worker

This project is a Go implementation of an HTTP server and worker that processes incoming requests and sends them to a webhook in a specified format

How to Use
Clone the Repository:

git clone <repository_url>

Run the Server:

go run main.go
The server will start listening on port 8080 by default.

Send Requests:
Send POST requests to the server with JSON data in the specified format using postman

Example:

{
"ev": "contact_form_submitted",
"et": "form_submit",
"id": "cl_app_id_001",
"uid": "cl_app_id_001-uid-001",
"mid": "cl_app_id_001-uid-001",
"t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
"p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
"l": "en-US",
"sc": "1920 x 1080",
"atrk1": "form_varient",
"atrv1": "red_top",
"atrt1": "string",
"atrk2": "ref",
"atrv2": "XPOWJRICW993LKJD",
"atrt2": "string",
"uatrk1": "name",
"uatrv1": "iron man",
"uatrt1": "string",
"uatrk2": "email",
"uatrv2": "ironman@avengers.com",
"uatrt2": "string",
"uatrk3": "age",
"uatrv3": "32",
"uatrt3": "integer"
}
send as body with post request

# Concepts being used

HTTP Server: The project implements an HTTP server using Go's net/http package to handle incoming requests.

Goroutines: Each incoming request is processed in a separate goroutine, allowing the server to handle multiple requests concurrently.

Channels: Go channels are used to communicate between the HTTP server and the worker goroutines. The server sends requests to the worker goroutines through channels for further processing.

JSON Encoding/Decoding: JSON encoding and decoding are used to parse incoming requests and format outgoing data. Go's encoding/json package is utilized for this purpose.

Testing: The project includes unit tests written using Go's built-in testing framework (testing package) to ensure the correctness of the implemented functionalities.

Concurrency: Concurrency is achieved through goroutines and channels, allowing for efficient handling of multiple incoming requests simultaneously.

Error Handling: Error handling is implemented using Go's error handling mechanisms

HTTP Client: Go's net/http package is used to send HTTP POST requests to the specified webhook URL.

These concepts together form the foundation of the HTTP server and worker implementation in the project.
