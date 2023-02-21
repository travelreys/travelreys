package images

type ImageOwnerLinks struct {
	Self   string `json:"self"`
	Html   string `json:"html"`
	Photos string `json:"photos"`
	Likes  string `json:"likes"`
}

type ImageOwner struct {
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Name     string          `json:"name"`
	Links    ImageOwnerLinks `json:"links"`
}

type ImageURLS struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

type ImageMetadataLinks struct {
	Self     string `json:"self"`
	Html     string `json:"html"`
	Download string `json:"download"`
}

type ImageMetadata struct {
	ID       string             `json:"id"`
	Width    float64            `json:"width"`
	Height   float64            `json:"height"`
	BlurHash string             `json:"blur_hash"`
	User     ImageOwner         `json:"user"`
	Urls     ImageURLS          `json:"urls"`
	Links    ImageMetadataLinks `json:"links"`
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
