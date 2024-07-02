package data

var ExcelPermissionMap map[int64]string = map[int64]string{}

var AdminMap map[int64]string = map[int64]string{
	1415509974: "admin_tang", // tang
}

var BannedMap map[int64]string = map[int64]string{}
var FoodyMap map[int64]string = map[int64]string{}

type RestaurantInfo struct {
	Recommender string `json:"recommender"`
	RecID       int64  `json:"recommender_id"`
	Info        string `json:"info"`
}

var RestaurantMap [3]map[string]RestaurantInfo = [3]map[string]RestaurantInfo{
	make(map[string]RestaurantInfo),
	make(map[string]RestaurantInfo),
	make(map[string]RestaurantInfo),
}
var Cache map[string]string = map[string]string{}

var GifCache map[string]string = map[string]string{}
var MusicCahce map[string]string = map[string]string{}

// Struct to hold all the maps for easy marshaling/unmarshaling
type SaveData struct {
	ExcelPermissionMap map[int64]string            `json:"excel_permission_map"`
	AdminMap           map[int64]string            `json:"admin_map"`
	BannedMap          map[int64]string            `json:"banned_map"`
	Cache              map[string]string           `json:"cache"`
	FoodyMap           map[int64]string            `json:"foody_map"`
	RestaurantMap      []map[string]RestaurantInfo `json:"restaurant_map"`
}

type AddRestaurantState struct {
	State       int
	Name        string
	PriceStr    string
	Description string
}

var UserStates = make(map[int64]*AddRestaurantState)

var LastMessageID = make(map[int64]int)
