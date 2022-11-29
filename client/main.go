package main

import (
	"context"
	"io"
	"log"

	"github.com/daluisgarcia/golang-probuffers-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:5061", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)

	// DoUnary(c)
	// DoClientStreaming(c)

	// DoServerStreaming(c)

	DoBidirectionalStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GetTest RPC: %v\n", err)
	}

	log.Printf("Response from GetTest: %v\n", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q1-t2",
			Question: "Question 1",
			Answer:   "aaafasdasd",
			TestId:   "t2",
		},
		{
			Id:       "q2-t2",
			Question: "Question 2",
			Answer:   "aaafasdasd",
			TestId:   "t2",
		},
		{
			Id:       "q1-t3",
			Question: "Question 3",
			Answer:   "aaafasdasd",
			TestId:   "t2",
		},
	}

	stream, err := c.SetQuestions(context.Background())

	if err != nil {
		log.Fatalf("Error while calling SetQuestions: %v\n", err)
	}

	for _, question := range questions {
		stream.Send(question)
	}

	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from SetQuestions: %v\n", err)
	}

	log.Printf("Response from SetQuestions: %v\n", msg)
}

func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "t1",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GetStudentsPerTest: %v\n", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}

		log.Printf("Response from GetStudentsPerTest: %v\n", msg)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := &testpb.TakeTestRequest{
		Answer: "aaafasdasd",
	}

	numberOfQuestions := 3

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())

	if err != nil {
		log.Fatalf("Error while calling TakeTest: %v\n", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while reading stream: %v\n", err)
			}

			log.Printf("Response from TakeTest: %v\n", msg)
		}

		close(waitChannel)
	}()

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(answer)
		}

		stream.CloseSend()
	}()
	<-waitChannel
}
