package services

import (
	"context"
	"go-newsletter/internal/config"
	"go-newsletter/internal/models"
	"log/slog"
	"time"

	"github.com/resend/resend-go/v2"
)

type MailingService struct {
	cfg    *config.ResendConfig
	logger *slog.Logger
}

func NewMailingService(cfg *config.ResendConfig, logger *slog.Logger) *MailingService {
	return &MailingService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *MailingService) SendMail(to []string, subject string, html string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := resend.NewClient(s.cfg.ApiKey)

	params := &resend.SendEmailRequest{
		From:    s.cfg.Sender,
		To:      to,
		Subject: subject,
		Html:    html,
	}

	_, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		s.logger.ErrorContext(ctx, "Error when sending mail", "error", err)
		return models.NewInternalServerError("Failed to send email")
	}
	s.logger.Info("Email sent")
	return nil
}
