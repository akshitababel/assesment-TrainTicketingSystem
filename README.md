Train Ticketing gRPC API Documentation
 Overview
This gRPC service allows users to:
✅ Purchase train tickets
✅ View their receipt
✅ Check seat allocations
✅ Modify seat assignments
✅ Remove users from the train

API Endpoints & Details
 1. Purchase Ticket
 Description: Allows a user to purchase a ticket and get seat allocation.
Request:
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com"
}

Response:
{
  "receipt": {
    "from": "London",
    "to": "France",
    "user": "John Doe",
    "price_paid": 20,
    "seat_number": "A1"
  }
}

2. Get Receipt
Description: Fetches the receipt for a user’s ticket purchase.
Request:
{
  "email": "john@example.com"
}
Response:
{
  "from": "London",
  "to": "France",
  "user": "John Doe",
  "price_paid": 20,
  "seat_number": "A1"
}

3. Get Users by Section
 Description: Retrieves all users seated in a specific train section.
 Request:
{
  "section": "A"
}
Response:
{
  "users": [
    { "name": "John Doe", "seat_number": "A1" },
    { "name": "Jane Smith", "seat_number": "A2" }
  ]
}

4. Remove User
Description: Removes a user from the train and frees up their seat.
Request:
{
  "email": "john@example.com"
}

Response:
{
  "message": "User removed successfully"
}

5. Modify Seat
Description: Changes a user’s seat allocation.
Request:
{
  "email": "john@example.com",
  "new_seat_number": "B5"
}
Response:
{
  "message": "Seat modified successfully"
}


gRPC Server Implementation
Run gRPC Server: go run server.go
Run gRPC Client:go run client.go
Generate gRPC Code from .proto:protoc --go_out=. --go-grpc_out=. train.proto
