package server

import (
	"context"
	"time"

	"github.com/Koyo-os/form-crud-service/internal/entity"
	"github.com/Koyo-os/form-crud-service/internal/metrics"
	"github.com/Koyo-os/form-crud-service/internal/service"
	"github.com/Koyo-os/form-crud-service/pkg/api/pb"
	"github.com/Koyo-os/form-crud-service/pkg/logger"
)

type Server struct {
	pb.UnimplementedFormServiceServer
	logger  *logger.Logger
	service *service.Service
}

func NewServer(service *service.Service) *Server {
	return &Server{
		service: service,
		logger:  logger.Get(),
	}
}

func (s *Server) Create(ctx context.Context, req *pb.RequestCreate) (*pb.Response, error) {
	start := time.Now()
	defer func ()  {
		metrics.RequestCount.WithLabelValues("create").Inc()
		metrics.RequestDuration.WithLabelValues("create").Observe(time.Since(start).Seconds())
	}()

	if err := s.service.Create(
		entity.ToEntityForm(req.Form),
	); err != nil {
		return &pb.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}

	return &pb.Response{
		Ok: true,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.RequestUpdate) (*pb.Response, error) {
	start := time.Now()
	defer func ()  {
		metrics.RequestCount.WithLabelValues("update").Inc()
		metrics.RequestDuration.WithLabelValues("update").Observe(time.Since(start).Seconds())
	}()

	if err := s.service.Update(req.ID, req.Key, req.Value); err != nil {
		return &pb.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}
	return &pb.Response{
		Ok: true,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.RequestDelete) (*pb.Response, error) {
	start := time.Now()
	defer func ()  {
		metrics.RequestCount.WithLabelValues("delete").Inc()
		metrics.RequestDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if err := s.service.Delete(req.ID); err != nil {
		return &pb.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}
	return &pb.Response{
		Ok: true,
	}, nil
}

func  (s *Server) Get(ctx context.Context, req *pb.RequestGet) (*pb.GetResponse, error) {
	start := time.Now()
	defer func ()  {
		metrics.RequestCount.WithLabelValues("get").Inc()
		metrics.RequestDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	}()

	form, err := s.service.Get(req.ID)
	if err != nil{
		return &pb.GetResponse{
			Response: &pb.Response{
				Error: err.Error(),
				Ok: false,
			},
			Form: nil,
		}, nil
	}

	return &pb.GetResponse{
		Form: form.ToProtobuf(),
		Response: &pb.Response{
			Error: "",
			Ok: true,
		},
	}, nil
}

func (s *Server) GetMore(ctx context.Context, req *pb.RequestGetMore) (*pb.GetMoreResponse, error) {
	start := time.Now()
	defer func ()  {
		metrics.RequestCount.WithLabelValues("get_more").Inc()
		metrics.RequestDuration.WithLabelValues("get_more").Observe(time.Since(start).Seconds())
	}()

	forms, err := s.service.GetMore(req.Key, req.Value)
	if err != nil{
		return &pb.GetMoreResponse{
			Forms: nil,
			Response: &pb.Response{
				Error: "",
				Ok: false,
			},
		}, err
	}

	pbForms := make([]*pb.Form, len(forms))

	for i, form := range forms{
		pbForms[i] = form.ToProtobuf()
	}

	return &pb.GetMoreResponse{
		Forms: pbForms,
		Response: &pb.Response{
			Error: "",
			Ok: true,
		},
	}, nil
}