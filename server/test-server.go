package server

import (
	"context"

	"github.com/daluisgarcia/golang-probuffers-grpc/models"
	"github.com/daluisgarcia/golang-probuffers-grpc/repository"
	"github.com/daluisgarcia/golang-probuffers-grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, request *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {

	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}

	err := s.repo.SetTest(ctx, test)

	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}
