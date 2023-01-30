package model

var role = map[string]int{
	RoleUnidentified.String(): 0,
	RoleNormal.String():       1,
	RoleAdmin.String():        9,
}

var RoleReverse = map[int]Role{
	0: RoleUnidentified,
	1: RoleNormal,
	9: RoleAdmin,
}

func (r Role) Int() int {
	return role[r.String()]
}

var MediaTypeReverse = map[string]MediaType{
	"photo": MediaTypePhoto,
	"video": MediaTypeVideo,
}

var MediaTypeToString = map[MediaType]string{
	MediaTypePhoto: "photo",
	MediaTypeVideo: "video",
}

var ActionUpdateAlbumMediasReverse = map[string]ActionUpdateAlbumMedias{
	"add":    ActionUpdateAlbumMediasAdd,
	"remove": ActionUpdateAlbumMediasRemove,
}

var ActionUpdateAlbumMediasToString = map[ActionUpdateAlbumMedias]string{
	ActionUpdateAlbumMediasAdd:    "add",
	ActionUpdateAlbumMediasRemove: "remove",
}
