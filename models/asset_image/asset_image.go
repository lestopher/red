package asset_image

import "time"

type AssetImage struct {
	ID              int
	Active          bool
	ContentType     string
	CreatedAt       time.Time
	FileSize        int
	Image           string
	ImageProcessing bool
	ParentName      string
	ParentMongoID   string
	ParentType      string
	MongoID         string
	Name            string
	UpdatedAt       time.Time
}
