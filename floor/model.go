package floor

var (
	Prefix = "[哈哈]"
)

type UserInfo struct {
	Avatar   string
	Id       int
	Role     int
	Nick     string
	Age      int
	Gender   int
	Level    int
	Birthday int
}

type Post struct {
	Id       int
	Title    string
	Category string
}

type Category struct {
	Id          int
	Title       string
	IsSubscribe int
}

type Comment struct {
	Id       int
	Text     string
	UserInfo UserInfo
}
