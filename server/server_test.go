package main

import (
	"context"
	"testing"

	pb "github.com/akshitababel/assesment-TrainTicketingSystem/train_proto/github.com/akshitababel/assesment-TrainTicketingSystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock gRPC Server
type MockTrainService struct {
	mock.Mock
	pb.UnimplementedTrainTicketServiceServer
}

// Mock PurchaseTicket API
func (m *MockTrainService) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.Receipt), args.Error(1)
}

// Mock GetReceipt API
func (m *MockTrainService) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.Receipt, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.Receipt), args.Error(1)
}

// Mock GetUsersBySection API
func (m *MockTrainService) GetUsersBySection(ctx context.Context, req *pb.SectionRequest) (*pb.UserList, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.UserList), args.Error(1)
}

// Mock RemoveUser API
func (m *MockTrainService) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.RemoveUserResponse), args.Error(1)
}

// Mock ModifySeat API
func (m *MockTrainService) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.ModifySeatResponse), args.Error(1)
}

// 📌 Test Case: PurchaseTicket
func TestPurchaseTicket(t *testing.T) {
	mockService := new(MockTrainService)
	req := &pb.PurchaseRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	expectedResp := &pb.Receipt{
		From:       "London",
		To:         "France",
		UserName:   "John Doe",
		PricePaid:  20,
		SeatNumber: "A1",
	}

	mockService.On("PurchaseTicket", mock.Anything, req).Return(expectedResp, nil)

	resp, err := mockService.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

// 📌 Test Case: GetReceipt
func TestGetReceipt(t *testing.T) {
	mockService := new(MockTrainService)
	req := &pb.ReceiptRequest{Email: "john@example.com"}
	expectedResp := &pb.Receipt{
		From:       "London",
		To:         "France",
		UserName:   "John Doe",
		PricePaid:  20,
		SeatNumber: "A1",
	}

	mockService.On("GetReceipt", mock.Anything, req).Return(expectedResp, nil)

	resp, err := mockService.GetReceipt(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

// 📌 Test Case: RemoveUser
func TestRemoveUser(t *testing.T) {
	mockService := new(MockTrainService)
	req := &pb.RemoveUserRequest{Email: "john@example.com"}
	expectedResp := &pb.RemoveUserResponse{Message: "User removed successfully"}

	mockService.On("RemoveUser", mock.Anything, req).Return(expectedResp, nil)

	resp, err := mockService.RemoveUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

// 📌 Test Case: ModifySeat
func TestModifySeat(t *testing.T) {
	mockService := new(MockTrainService)
	req := &pb.ModifySeatRequest{
		Email:         "john@example.com",
		NewSeatNumber: "B5",
	}
	expectedResp := &pb.ModifySeatResponse{Message: "Seat modified successfully"}

	mockService.On("ModifySeat", mock.Anything, req).Return(expectedResp, nil)

	resp, err := mockService.ModifySeat(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}
