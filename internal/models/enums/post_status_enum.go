package enums

type PostStatus string

const (
	Posted    PostStatus = "POSTED"
	Scheduled PostStatus = "SCHEDULED"
)

func (s PostStatus) String() string {
	return string(s)
}
