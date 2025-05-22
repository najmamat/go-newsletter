# Create message
- Editor fill a title of message
- Editor creates an html content
- System parses html content to plaintext format
- System verifies that the editor is authenticated
- System checks that the editor has write access to the target newsletter

# Publish message immediately
- Editor choose to publish message immediately
- System creates a new post record with status set to "published"
- System retrieves the list of confirmed subscribers for the newsletter
- System sends the post via email to all confirmed subscribers
- System logs the sent status and timestamps (published_at)
- System returns the published post object

# TODO:
- Save message as a draft
- Schedule publishing
- Cancel scheduled publishing
- Update schedule


# Questions
- Which html elements will system support?
- Are there any html elements which cannot be parsed to plaintext / are ugly in plaintext?