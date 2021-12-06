package domain

type Board struct {
	Id          string `bson:"_id"`
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description,omitempty"`
	ImageUrl    string `bson:"image_url,omitempty"`
	IsAdult     bool   `bson:"is_adult"`
}

type Thread struct {
	Id        string `bson:"_id"`
	BoardId   string `bson:"board_id"`
	PostCount int64  `bson:"post_count,omitempty"`
	Title     string `bson:"title"`
	CreatedAt uint64 `bson:"created_at"`
	BumpedAt  uint64 `bson:"bumped_at"`
}

type Post struct {
	Id          string   `bson:"_id"`
	ThreadId    string   `bson:"thread_id"`
	Description string   `bson:"description,omitempty"`
	QuoteIds    []string `bson:"quote_ids,omitempty"`
	Text        string   `bson:"text,omitempty"`
	ImageUrl    string   `bson:"image_url,omitempty"`
	VideoUrl    string   `bson:"video_url,omitempty"`
	CreatedAt   uint64   `bson:"created_at"`
}
