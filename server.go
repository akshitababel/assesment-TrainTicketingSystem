package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "train_proto" // Replace with your actual import path

	"google.golang.org/grpc"
)

// Seat allocation data structure
type TrainTicketService struct {
	pb.UnimplementedTrainTicketServiceServer
	mu      sync.Mutex
	seats   map[string]string     // email -> seat_number
	users   map[string]pb.Receipt // email -> Receipt details
	section map[string][]string   // section -> seat numbers
}

// Initialize the train with available seats
func NewTrainTicketService() *TrainTicketService {
	return &TrainTicketService{
		seats:   make(map[string]string),
		users:   make(map[string]pb.Receipt),
		section: map[string][]string{"A": generateSeats("A", 10), "B": generateSeats("B", 10)},
	}
}

// Generate seat numbers for sections
func generateSeats(section string, count int) []string {
	seats := make([]string, count)
	for i := 0; i < count; i++ {
		seats[i] = fmt.Sprintf("%s-%d", section, i+1)
	}
	return seats
}

// Assigns a random available seat in a section
func (s *TrainTicketService) assignSeat() (string, string, bool) {
	for sec, seats := range s.section {
		if len(seats) > 0 {
			seat := seats[0]           // Take the first available seat
			s.section[sec] = seats[1:] // Remove assigned seat from available list
			return sec, seat, true
		}
	}
	return "", "", false
}

// gRPC Method: Purchase Ticket
func (s *TrainTicketService) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the user already booked
	if _, exists := s.users[req.Email]; exists {
		return nil, fmt.Errorf("user already has a ticket")
	}

	// Assign a seat
	section, seat, ok := s.assignSeat()
	if !ok {
		return nil, fmt.Errorf("train is full")
	}

	// Create receipt
	receipt := pb.Receipt{
		From:       "London",
		To:         "France",
		UserName:   req.FirstName + " " + req.LastName,
		Email:      req.Email,
		PricePaid:  20.0,
		Section:    section,
		SeatNumber: seat,
	}

	// Store user details
	s.users[req.Email] = receipt
	s.seats[req.Email] = seat

	return &receipt, nil
}

// gRPC Method: Get Receipt
func (s *TrainTicketService) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[req.Email]
	if !exists {
		return nil, fmt.Errorf("receipt not found")
	}

	return &receipt, nil
}

// gRPC Method: Get Users by Section
func (s *TrainTicketService) GetUsersBySection(ctx context.Context, req *pb.SectionRequest) (*pb.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userList pb.UserList
	for email, receipt := range s.users {
		if receipt.Section == req.Section {
			userList.Users = append(userList.Users, &pb.UserInfo{Name: receipt.UserName, SeatNumber: receipt.SeatNumber})
		}
	}

	return &userList, nil
}

// gRPC Method: Remove User from Train
func (s *TrainTicketService) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seat, exists := s.seats[req.Email]
	if !exists {
		return &pb.RemoveUserResponse{Success: false, Message: "user not found"}, nil
	}

	// Free the seat
	section := s.users[req.Email].Section
	s.section[section] = append(s.section[section], seat)

	// Remove user
	delete(s.users, req.Email)
	delete(s.seats, req.Email)

	return &pb.RemoveUserResponse{Success: true, Message: "user removed"}, nil
}

// gRPC Method: Modify Seat
func (s *TrainTicketService) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[req.Email]
	if !exists {
		return &pb.ModifySeatResponse{Success: false, Message: "user not found"}, nil
	}

	// Free the old seat
	oldSeat := receipt.SeatNumber
	section := receipt.Section
	s.section[section] = append(s.section[section], oldSeat)

	// Assign new seat
	s.users[req.Email] = pb.Receipt{
		From:       receipt.From,
		To:         receipt.To,
		UserName:   receipt.UserName,
		Email:      receipt.Email,
		PricePaid:  receipt.PricePaid,
		Section:    section,
		SeatNumber: req.NewSeatNumber,
	}
	s.seats[req.Email] = req.NewSeatNumber

	return &pb.ModifySeatResponse{Success: true, Message: "seat updated"}, nil
}

// Start gRPC server
func main() {
	server := grpc.NewServer()
	service := NewTrainTicketService()
	pb.RegisterTrainTicketServiceServer(server, service)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
