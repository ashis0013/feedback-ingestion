package main

import (
	"encoding/json"
	"log"

	"github.com/ashis0013/feedback-ingestion/cron"
	"github.com/ashis0013/feedback-ingestion/models"
	"github.com/ashis0013/feedback-ingestion/repository"
	"github.com/ashis0013/feedback-ingestion/service"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	service *service.FeedbackIngestionService
}

func (s *Server) handlePostFeedback(c *fiber.Ctx) error {
	req := &models.AddFeedbackRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(401).SendString("Bad request")
	}
	err := s.service.AddFedback(req)
	if err != nil {
		return c.Status(500).SendString("Internal server error")
	}
	return c.SendStatus(200)
}

func (s *Server) handleGetFeedback(c *fiber.Ctx) error {
	var filter models.QueryFilter
	err := json.Unmarshal([]byte(c.Query("json")), &filter)
	if err != nil || filter.IsInvalid() {
		return c.Status(400).SendString("Invalid input")
	}
	res, err := s.service.GetFeedback(&filter)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(res)
}

func (s *Server) handlePostTenant(c *fiber.Ctx) error {
	req := new(models.AddTenantRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(401).SendString("Bad request")
	}
	err := s.service.OnboardTenant(req)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(201)
}

func (s *Server) route() {
	s.app.Post("/feedback", s.handlePostFeedback)
	s.app.Get("/bruh", s.handleGetFeedback)
	s.app.Post("/tenant", s.handlePostTenant)
}

func (s *Server) Start() {
	s.route()
	log.Fatal(s.app.Listen(":3000"))
}

func NewServer(svc *service.FeedbackIngestionService) *Server {
	return &Server{
		app:     fiber.New(),
		service: svc,
	}
}

func main() {
	repository := repository.NewPostgresRepository()
	repository.Init()
	defer repository.Terminate()

	server := NewServer(service.NewFeedbackIngestionService(repository, []cron.IngestionModule{}))
	server.Start()
}
