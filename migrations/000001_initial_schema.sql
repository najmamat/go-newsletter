-- Create profiles table
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    full_name TEXT,
    avatar_url TEXT,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE profiles IS 'Stores additional information for editors (including admin status), linked to Supabase''s auth.users table.';
COMMENT ON COLUMN profiles.id IS 'User ID from Supabase Auth.';
COMMENT ON COLUMN profiles.full_name IS 'Full name of the editor.';
COMMENT ON COLUMN profiles.avatar_url IS 'URL to the editor''s avatar image.';
COMMENT ON COLUMN profiles.is_admin IS 'Flag to indicate if the user is an administrator.';
COMMENT ON COLUMN profiles.created_at IS 'Timestamp of profile creation.';
COMMENT ON COLUMN profiles.updated_at IS 'Timestamp of last profile update.';

-- Create newsletters table
CREATE TABLE IF NOT EXISTS newsletters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    editor_id UUID NOT NULL REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE newsletters IS 'Contains details about each newsletter.';
COMMENT ON COLUMN newsletters.id IS 'Unique identifier for the newsletter.';
COMMENT ON COLUMN newsletters.name IS 'Name of the newsletter.';
COMMENT ON COLUMN newsletters.description IS 'Optional description of the newsletter.';
COMMENT ON COLUMN newsletters.editor_id IS 'ID of the editor who owns this newsletter';
COMMENT ON COLUMN newsletters.created_at IS 'Timestamp of newsletter creation.';
COMMENT ON COLUMN newsletters.updated_at IS 'Timestamp of last newsletter update.';

-- Create newsletter_editors table
CREATE TABLE IF NOT EXISTS newsletter_editors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    newsletter_id UUID NOT NULL REFERENCES newsletters(id),
    editor_id UUID NOT NULL REFERENCES auth.users(id),
    role TEXT CHECK (role = ANY (ARRAY['owner'::text, 'editor'::text, 'viewer'::text])) DEFAULT 'editor',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE newsletter_editors IS 'Junction table managing many-to-many relationship between newsletters and editors with role-based access';
COMMENT ON COLUMN newsletter_editors.newsletter_id IS 'References the newsletter being edited';
COMMENT ON COLUMN newsletter_editors.editor_id IS 'References the user/editor with access to the newsletter';
COMMENT ON COLUMN newsletter_editors.role IS 'Role of the editor: owner (full control), editor (can edit content), viewer (read-only)';

-- Create subscribers table
CREATE TABLE IF NOT EXISTS subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    newsletter_id UUID NOT NULL REFERENCES newsletters(id),
    email TEXT NOT NULL,
    subscribed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unsubscribed_at TIMESTAMPTZ,
    unsubscribe_token TEXT NOT NULL UNIQUE,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirmation_token TEXT UNIQUE
);

COMMENT ON TABLE subscribers IS 'Manages email subscriptions to newsletters.';
COMMENT ON COLUMN subscribers.id IS 'Unique identifier for the subscription record.';
COMMENT ON COLUMN subscribers.newsletter_id IS 'ID of the newsletter being subscribed to.';
COMMENT ON COLUMN subscribers.email IS 'Email address of the subscriber.';
COMMENT ON COLUMN subscribers.subscribed_at IS 'Timestamp when the subscription was made.';
COMMENT ON COLUMN subscribers.unsubscribed_at IS 'Timestamp when the user unsubscribed, if applicable.';
COMMENT ON COLUMN subscribers.unsubscribe_token IS 'Unique token for one-click unsubscription.';
COMMENT ON COLUMN subscribers.is_confirmed IS 'Flag to indicate if the subscription has been confirmed via email.';
COMMENT ON COLUMN subscribers.confirmation_token IS 'Unique token for email confirmation.';

-- Create published_posts table
CREATE TABLE IF NOT EXISTS published_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    newsletter_id UUID NOT NULL REFERENCES newsletters(id),
    editor_id UUID NOT NULL REFERENCES auth.users(id),
    title TEXT NOT NULL,
    content_html TEXT NOT NULL,
    content_text TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    scheduled_at TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE published_posts IS 'Archives posts published to newsletters.';
COMMENT ON COLUMN published_posts.id IS 'Unique identifier for the published post.';
COMMENT ON COLUMN published_posts.newsletter_id IS 'ID of the newsletter this post belongs to.';
COMMENT ON COLUMN published_posts.editor_id IS 'ID of the editor who published this post.';
COMMENT ON COLUMN published_posts.title IS 'Title of the post.';
COMMENT ON COLUMN published_posts.content_html IS 'HTML content of the post.';
COMMENT ON COLUMN published_posts.content_text IS 'Plain text version of the post content.';
COMMENT ON COLUMN published_posts.status IS 'Status of the post (e.g., ''draft'', ''scheduled'', ''publishing'', ''published'', ''failed'').';
COMMENT ON COLUMN published_posts.scheduled_at IS 'If status is ''scheduled'', this is the time it will be published (UTC).';
COMMENT ON COLUMN published_posts.published_at IS 'Timestamp when the post was actually published. Becomes non-null when status is ''published''.';
COMMENT ON COLUMN published_posts.created_at IS 'Timestamp of post record creation.';

-- Add unique constraint to prevent duplicate subscriptions
ALTER TABLE subscribers ADD CONSTRAINT unique_newsletter_subscriber UNIQUE (newsletter_id, email); 