package repository

import (
	"context"
	"go-newsletter/internal/models"
	"go-newsletter/internal/models/enums"
	"log/slog"
	"time"

	"go-newsletter/pkg/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostRepository(db *pgxpool.Pool, logger *slog.Logger) *PostRepository {
	return &PostRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostRepository) GetPostsByNewsletterId(ctx context.Context, newsletterID uuid.UUID, published bool) ([]*generated.PublishedPost, error) {
	query := `
		SELECT id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
		FROM published_posts
		WHERE newsletter_id = $1`

	if published {
		query += ` AND published_at IS NOT NULL`
	} else {
		query += ` AND published_at IS NULL`
	}

	rows, err := r.db.Query(ctx, query, newsletterID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query posts", "error", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*generated.PublishedPost
	for rows.Next() {
		s := &generated.PublishedPost{}
		err := rows.Scan(
			&s.Id,
			&s.NewsletterId,
			&s.EditorId,
			&s.Title,
			&s.ContentHtml,
			&s.ContentText,
			&s.Status,
			&s.ScheduledAt,
			&s.PublishedAt,
			&s.CreatedAt,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan post row", "error", err)
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating post rows", "error", err)
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostById(ctx context.Context, postId uuid.UUID) (*generated.PublishedPost, error) {
	query := `
		SELECT id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
		FROM published_posts
		WHERE id = $1`

	post := &generated.PublishedPost{}
	err := r.db.QueryRow(ctx, query, postId).Scan(
		&post.Id,
		&post.NewsletterId,
		&post.EditorId,
		&post.Title,
		&post.ContentHtml,
		&post.ContentText,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
	)

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query post", "error", err)
		return nil, err
	}

	return post, nil
}

// GetPostsDueForPublication returns all scheduled posts that are due for publication
func (r *PostRepository) GetPostsDueForPublication(ctx context.Context, currentTime time.Time) ([]*generated.PublishedPost, error) {
	query := `
		SELECT id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
		FROM published_posts
		WHERE status = $1
		AND scheduled_at <= $2
		AND published_at IS NULL
	`

	rows, err := r.db.Query(ctx, query, enums.Scheduled.String(), currentTime)
	if err != nil {
		r.logger.ErrorContext(ctx, "Error loading posts for publication", "error", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*generated.PublishedPost
	for rows.Next() {
		s := &generated.PublishedPost{}
		err := rows.Scan(
			&s.Id,
			&s.NewsletterId,
			&s.EditorId,
			&s.Title,
			&s.ContentHtml,
			&s.ContentText,
			&s.Status,
			&s.ScheduledAt,
			&s.PublishedAt,
			&s.CreatedAt,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Error reading post row", "error", err)
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating results", "error", err)
		return nil, err
	}

	return posts, nil
}

// PublishPost updates the status of a post to published
func (r *PostRepository) PublishPost(ctx context.Context, postId uuid.UUID) error {
	query := `
		UPDATE published_posts
		SET status = $2, published_at = $3
		WHERE id = $1
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, postId, enums.Posted.String(), now)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: error publishing post", "id", postId, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.ErrorContext(ctx, "REPO: Post not found for publishing", "id", postId)
		return models.NewNotFoundError("Post not found")
	}

	return nil
}

func (r *PostRepository) DeletePostById(ctx context.Context, postId uuid.UUID) error {
	query := `
		DELETE
		FROM published_posts
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, postId)
	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: failed to delete post", "id", postId, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.ErrorContext(ctx, "REPO: Post not found for deletion", "id", postId)
		return models.NewNotFoundError("Post not found")
	}

	return nil
}

func (r *PostRepository) CreatePost(ctx context.Context, userId uuid.UUID, createPost *generated.PublishPostRequest, newsletterId uuid.UUID) (*generated.PublishedPost, error) {
	query := `
	INSERT INTO published_posts (id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
	`

	id := uuid.New()
	now := time.Now()

	status := enums.Scheduled
	var publishedAt *time.Time

	if createPost.ScheduledAt != nil && createPost.ScheduledAt.Before(now) || createPost.ScheduledAt.Equal(now) {
		status = enums.Posted
		publishedAt = &now
	}

	post := &generated.PublishedPost{}
	err := r.db.QueryRow(ctx, query,
		id,
		newsletterId,
		userId,
		createPost.Title,
		createPost.ContentHtml,
		createPost.ContentText,
		status.String(),
		createPost.ScheduledAt,
		publishedAt,
		now,
	).Scan(
		&post.Id,
		&post.NewsletterId,
		&post.EditorId,
		&post.Title,
		&post.ContentHtml,
		&post.ContentText,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
	)

	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: failed to create post", "error", err)
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, postId uuid.UUID, updatePost *generated.PublishPostRequest) (*generated.PublishedPost, error) {
	query := `
	UPDATE published_posts 
	SET title = $2, content_html = $3, content_text = $4, status = $5, scheduled_at = $6, published_at = $7
	WHERE id = $1
	RETURNING id, newsletter_id, editor_id, title, content_html, content_text, status, scheduled_at, published_at, created_at
	`

	now := time.Now()

	originalPost, err := r.GetPostById(ctx, postId)
	if err != nil {
		return nil, err
	}

	status := enums.Scheduled
	var publishedAt *time.Time

	if originalPost.PublishedAt != nil {
		status = enums.Posted
		publishedAt = originalPost.PublishedAt
	} else if updatePost.ScheduledAt != nil && (updatePost.ScheduledAt.Before(now) || updatePost.ScheduledAt.Equal(now)) {
		status = enums.Posted
		publishedAt = &now
	}

	post := &generated.PublishedPost{}
	err = r.db.QueryRow(ctx, query,
		postId,
		updatePost.Title,
		updatePost.ContentHtml,
		updatePost.ContentText,
		status.String(),
		updatePost.ScheduledAt,
		publishedAt,
	).Scan(
		&post.Id,
		&post.NewsletterId,
		&post.EditorId,
		&post.Title,
		&post.ContentHtml,
		&post.ContentText,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
	)

	if err != nil {
		r.logger.ErrorContext(ctx, "REPO: failed to update post", "error", err)
		return nil, err
	}

	return post, nil
}
