package server

import (
	"context"
	"io"

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

// The type of the stream receives the requests during the connection is up and this function processes the requests
func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			// The client has finished sending data
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		var question = &models.Question{
			Id:       msg.Id,
			Question: msg.Question,
			TestId:   msg.TestId,
			Answer:   msg.Answer,
		}

		err = s.repo.SetQuestion(stream.Context(), question)

		if err != nil {
			return err
		}
	}
}
