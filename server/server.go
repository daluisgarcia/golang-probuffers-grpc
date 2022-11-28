package server

import (
	"context"

	"github.com/daluisgarcia/golang-probuffers-grpc/models"
	"github.com/daluisgarcia/golang-probuffers-grpc/repository"
	"github.com/daluisgarcia/golang-probuffers-grpc/studentpb"
)

type Server struct {
	repo repository.Repository
	studentpb.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetStudent(ctx context.Context, request *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	student, err := s.repo.GetStudent(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	return &studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {

	student := &models.Student{
		Id:   req.Id,
		Name: req.Name,
		Age:  req.Age,
	}

	err := s.repo.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &studentpb.SetStudentResponse{
		Student: req,
	}, nil
}

func (s *Server) StudentServiceServer() studentpb.StudentServiceServer {
	return s
}
