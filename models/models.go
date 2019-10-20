package models

import (
	"time"
)

type Music  struct{
	Id      		uint64 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`// sql:"type:bigint;"`
	Title   		string 	`gorm:"type:varchar(255);index:title"`
	GenreId 		uint    `gorm:"NOT NULL"`
	TimedLyricsPath string 	`gorm:"NULL;type:varchar(255)"`//file, each charater has start and end time in milliseconds, {'a':{"start":45,"end":50}}
	AudioPath 		string 	`gorm:"NULL;type:varchar(255)"`//path to audio file, TODO:linux folders can contain up to 268,173,300 files
	VideoPath 		string 	`gorm:"NULL; type:varchar(255)"`//path to karakoke video
	AuthorId 		uint    `gorm:"NULL"`//to whom original song belongs to 
	UserId			uint  	`gorm:"NOT NULL"`	
	OriginalMusicId uint    `gorm:"NOT NULL; DEFAULT:0"`//0 - for karaoke music with video, all the other for user covered song
	AlbumId  		uint    `gorm:"NULL"`//to pushlish with album
	Duration		string  `gorm:"NULL; type:varchar(10)"`
	CreatedAt 		*time.Time `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
}
type Admin struct{
	ID      	uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
    UserId      uint  	`gorm:"NULL; UNIQUE; DEFAULT:0"`  
}
type User struct{
	Id      	uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	Name 		string 	`gorm:"NOT NULL;type:varchar(255)"`//Ashat Sharipov
	Phone 	  	string 	`gorm:"NOT NULL;UNIQUE; type:varchar(20)"`//
	CreatedAt 	*time.Time  `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
	//Musics      []*Music    `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"omitempty"`
	//SecondName string `gorm:"NOT NULL; type:varchar(255)"`
	//Email 	  string `gorm:"NULL;UNIQUE; type:varchar(255)"`
}
type UserRating struct{
	Id      	uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	Rating    	uint  	`gorm:"NULL; DEFAULT:0"`
	MusicId     uint    `gorm:"NOT NULL;index:musicIdInRatings"`
	UserId      uint    `gorm:"NOT NULL"`
}
type Author struct{
	Id uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	Name string `gorm:"NOT NULL; PRIMARY_KEY;type:varchar(255);UNIQUE"`//Ashat Sharipov
}
type Genre struct{
	Id      uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	Name    string 	`gorm:"NOT NULL; type:varchar(255) UNIQUE"`
	Musics  []Music `gorm:"foreignkey:GenreId;association_foreignkey:ID" json:",omitempty"`
}
type Album struct{
	Id      uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	Name    string 	`gorm:"NOT NULL; type:varchar(255)"`
	AuthorId uint   `gorm:"NOT NULL;"`
	Musics  []*Music `gorm:"-"`
}
type PlayListHelper struct{
	Id  	uint 	`gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	MusicId uint64 	`gorm:" NOT NULL" sql:"type:bigint;"`
}
type PlayList struct{
	Id 			uint `gorm:"NOT NULL; PRIMARY_KEY; AUTO_INCREMENT"`
	PlayListId 	uint `gorm:"NOT NULL"`//id of playlisthelper
	UserId 		uint `gorm:"NOT NULL"`
	Name 		string `gorm:"NOT NULL; type:varchar(255)"`  
	CreatedAt   *time.Time  `gorm:"TYPE:DATETIME;NOT NULL" schema:"-" json:",omitempty"`
	Musics      []*Music `gorm:"-"`
}