package server

import (
	"context"
	"io"
	"log"

	"github.com/daluisgarcia/golang-probuffers-grpc/models"
	"github.com/daluisgarcia/golang-probuffers-grpc/repository"
	"github.com/daluisgarcia/golang-probuffers-grpc/studentpb"
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
		Id:   req.GetId(),
		Name: req.GetName(),
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
			Id:       msg.GetId(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
			Answer:   msg.GetAnswer(),
		}

		err = s.repo.SetQuestion(context.Background(), question)

		if err != nil {
			return err
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
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

		var enrollment = &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)

		if err != nil {
			return err
		}
	}
}

// The second parameter here is the stream that the server will use to send the data to the client
func (s *TestServer) GetStudentsPerTest(request *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), request.GetTestId())

	if err != nil {
		return err
	}

	for _, student := range students {
		err = stream.Send(&studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Uses the stream to send and receive data
func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "t1")

	if err != nil {
		return err
	}

	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}
		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Answer received: ", answer.GetAnswer())
	}
}
