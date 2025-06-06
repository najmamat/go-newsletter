package scheduler

import (
	"context"
	"go-newsletter/internal/services"
	"log/slog"
	"time"
)

// PostPublisher is a service for automatically publishing scheduled posts
type PostPublisher struct {
	postService *services.PostService
	interval    time.Duration
	shutdownCh  chan struct{}
	logger      *slog.Logger
}

// NewPostPublisher creates a new instance of PostPublisher
func NewPostPublisher(postService *services.PostService, logger *slog.Logger) *PostPublisher {
	return &PostPublisher{
		postService: postService,
		interval:    time.Minute, // Check every minute
		shutdownCh:  make(chan struct{}),
		logger:      logger,
	}
}

// Start begins the background publishing process
func (p *PostPublisher) Start() {
	p.logger.Info("Starting scheduled post publisher service")
	go p.run()
}

// Stop terminates the publishing process
func (p *PostPublisher) Stop() {
	p.logger.Info("Stopping scheduled post publisher service")
	close(p.shutdownCh)
}

// run is the main loop for checking and publishing scheduled posts
func (p *PostPublisher) run() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Check immediately upon starting
	p.publishScheduledPosts()

	for {
		select {
		case <-ticker.C:
			p.publishScheduledPosts()
		case <-p.shutdownCh:
			p.logger.Info("Scheduled post publisher service stopped")
			return
		}
	}
}

// publishScheduledPosts finds and publishes all posts whose publication time has arrived
func (p *PostPublisher) publishScheduledPosts() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	p.logger.InfoContext(ctx, "Checking for scheduled posts to publish")

	now := time.Now()
	p.logger.InfoContext(ctx, "Current time", "time", now.Format(time.RFC3339))

	// Get all scheduled posts that are due for publication
	posts, err := p.postService.GetPostsDueForPublication(ctx, now)
	if err != nil {
		p.logger.ErrorContext(ctx, "Error fetching scheduled posts", "error", err)
		return
	}

	if len(posts) == 0 {
		p.logger.InfoContext(ctx, "No posts scheduled for publication at this time")
		return
	}

	p.logger.InfoContext(ctx, "Found posts to publish", "count", len(posts))

	// Publish each post
	successCount := 0
	failureCount := 0

	for _, post := range posts {
		if post == nil || post.Id == nil {
			p.logger.WarnContext(ctx, "Invalid post data received from database")
			continue
		}

		p.logger.InfoContext(ctx, "Publishing post",
			"postId", post.Id,
			"title", post.Title,
			"scheduledAt", post.ScheduledAt.Format(time.RFC3339))

		err := p.postService.PublishPost(ctx, *post.Id)
		if err != nil {
			p.logger.ErrorContext(ctx, "Error publishing post", "postId", post.Id, "error", err)
			failureCount++
			continue
		}

		successCount++
		p.logger.InfoContext(ctx, "Post published successfully", "postId", post.Id, "title", post.Title)
	}

	p.logger.InfoContext(ctx, "Post publishing completed", "successCount", successCount, "failureCount", failureCount)
}
