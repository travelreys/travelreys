package images

type ImageOwnerLinks struct {
	Self   string `json:"self" bson:"self"`
	Html   string `json:"html" bson:"html"`
	Photos string `json:"photos" bson:"photos"`
	Likes  string `json:"likes" bson:"likes"`
}

type ImageOwner struct {
	ID       string          `json:"id" bson:"id"`
	Username string          `json:"username" bson:"username"`
	Name     string          `json:"name" bson:"name"`
	Links    ImageOwnerLinks `json:"links" bson:"links"`
}

type ImageURLS struct {
	Raw     string `json:"raw" bson:"raw"`
	Full    string `json:"full" bson:"full"`
	Regular string `json:"regular" bson:"regular"`
	Small   string `json:"small" bson:"small"`
	Thumb   string `json:"thumb" bson:"thumb"`
}

type ImageMetadataLinks struct {
	Self     string `json:"self" bson:"self"`
	Html     string `json:"html" bson:"html"`
	Download string `json:"download" bson:"download"`
}

type ImageMetadata struct {
	ID       string             `json:"id" bson:"id"`
	Width    float64            `json:"width" bson:"width"`
	Height   float64            `json:"height" bson:"height"`
	BlurHash string             `json:"blur_hash" bson:"blur_hash"`
	User     ImageOwner         `json:"user" bson:"user"`
	Urls     ImageURLS          `json:"urls" bson:"urls"`
	Links    ImageMetadataLinks `json:"links" bson:"links"`
}

type MetadataList []ImageMetadata

var (
	CoverStockImageList = MetadataList{
		{
			ID:       "qyAka7W5uMY",
			Width:    3423,
			Height:   2704,
			BlurHash: "LSKKyhr^8^M|Ek?btmRiMdxvROxb",
			User: ImageOwner{
				ID:       "IFcEhJqem0Q",
				Username: "anniespratt",
				Name:     "Annie Spratt",
				Links: ImageOwnerLinks{
					Self:   `https://api.unsplash.com/users/anniespratt`,
					Html:   `https://unsplash.com/@anniespratt`,
					Photos: `https://api.unsplash.com/users/anniespratt/photos`,
					Likes:  `https://api.unsplash.com/users/anniespratt/likes`,
				},
			},
			Urls: ImageURLS{
				Raw:     `https://images.unsplash.com/photo-1488646953014-85cb44e25828?ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA&ixlib=rb-4.0.3`,
				Full:    `https://images.unsplash.com/photo-1488646953014-85cb44e25828?crop=entropy&cs=tinysrgb&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA&ixlib=rb-4.0.3&q=80`,
				Regular: `https://images.unsplash.com/photo-1488646953014-85cb44e25828?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA&ixlib=rb-4.0.3&q=80&w=1080`,
				Small:   `https://images.unsplash.com/photo-1488646953014-85cb44e25828?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA&ixlib=rb-4.0.3&q=80&w=400`,
				Thumb:   `https://images.unsplash.com/photo-1488646953014-85cb44e25828?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA&ixlib=rb-4.0.3&q=80&w=200`,
			},
			Links: ImageMetadataLinks{
				Self:     `https://api.unsplash.com/photos/qyAka7W5uMY`,
				Html:     `https://unsplash.com/photos/qyAka7W5uMY`,
				Download: `https://unsplash.com/photos/qyAka7W5uMY/download?ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHx0cmF2ZWx8ZW58MHwwfHx8MTY3MzY2NjQ4MA`,
			},
		},
	}
)
