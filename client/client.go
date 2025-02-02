package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/akshitababel/assesment-TrainTicketingSystem/train_proto/github.com/akshitababel/assesment-TrainTicketingSystem"

	"google.golang.org/grpc"
)

// Connects to the gRPC server
func connect() (*grpc.ClientConn, pb.TrainTicketServiceClient) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := pb.NewTrainTicketServiceClient(conn)
	return conn, client
}

// API 1: Purchase Ticket
func purchaseTicket(client pb.TrainTicketServiceClient, firstName, lastName, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.PurchaseRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	resp, err := client.PurchaseTicket(ctx, req)
	if err != nil {
		log.Printf("Error purchasing ticket: %v", err)
		return
	}

	fmt.Printf("Ticket Purchased: %+v\n", resp)
}

// API 2: Get Receipt
func getReceipt(client pb.TrainTicketServiceClient, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ReceiptRequest{Email: email}

	resp, err := client.GetReceipt(ctx, req)
	if err != nil {
		log.Printf("Error fetching receipt: %v", err)
		return
	}

	fmt.Printf("Receipt Details: %+v\n", resp)
}

// API 3: Get Users by Section
func getUsersBySection(client pb.TrainTicketServiceClient, section string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.SectionRequest{Section: section}

	resp, err := client.GetUsersBySection(ctx, req)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}

	fmt.Println("Users in section", section)
	for _, user := range resp.Users {
		fmt.Printf("Name: %s, Seat: %s\n", user.Name, user.SeatNumber)
	}
}

// API 4: Remove User
func removeUser(client pb.TrainTicketServiceClient, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RemoveUserRequest{Email: email}

	resp, err := client.RemoveUser(ctx, req)
	if err != nil {
		log.Printf("Error removing user: %v", err)
		return
	}

	fmt.Println("Remove User Response:", resp.Message)
}

// API 5: Modify User Seat
func modifySeat(client pb.TrainTicketServiceClient, email, newSeat string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ModifySeatRequest{
		Email:         email,
		NewSeatNumber: newSeat,
	}

	resp, err := client.ModifySeat(ctx, req)
	if err != nil {
		log.Printf("Error modifying seat: %v", err)
		return
	}

	fmt.Println("Modify Seat Response:", resp.Message)
}

// Main function to test all APIs
func main() {
	conn, client := connect()
	defer conn.Close()

	// Sample Data
	email := "john@example.com"
	firstName := "John"
	lastName := "Doe"

	// Test APIs
	purchaseTicket(client, firstName, lastName, email)
	getReceipt(client, email)
	getUsersBySection(client, "A")
	modifySeat(client, email, "B-5")
	removeUser(client, email)
}
