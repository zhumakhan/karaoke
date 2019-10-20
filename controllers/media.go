package controllers
import (
	"main/database"
	m "main/models"
	c "main/constants"
	"github.com/gorilla/mux"
	"bytes"
	"net/http"
	"path/filepath"
	"os"
	"io"
	"time"
	"math/rand"
	"container/list"
	"github.com/dhowden/tag"
	"log"
	"fmt"
	"strconv"
	"errors"
	"os/exec"
	"regexp"

)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const OK = "OK"
var newMusics = list.New()
var newMusicsSize = 50

func UploadMusicForUser(w http.ResponseWriter, r *http.Request){
	err := r.ParseMultipartForm(256 << 20)//, ~100 MByte
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return	
	}
	user, err := Authorize(r.Header.Get("Authorization"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"You are not authorized!")
		return
	}
	music := m.Music{}
	i64, err := strconv.ParseInt(r.Form.Get("GenreID"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	music.GenreId = uint(i64)
	if  music.GenreId == 0{
		respondWithError(w, http.StatusBadRequest, "Genre is not defined")
		return	
	}
	i64, err = strconv.ParseInt(r.Form.Get("OriginalMusicId"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	music.OriginalMusicId = uint(i64)
	if  music.OriginalMusicId == 0{
		respondWithError(w, http.StatusBadRequest, "Original music id is not defined")
		return	
	}
	filePath, err:= processMediaFile(r, "audio")
 	if err != nil{
 		respondWithError(w, http.StatusBadRequest, err.Error())
		return
 	}
 	filePath, durationStr, _, _, err := Mp3ConverterAndDuration(filePath,false)
 	if err != nil{
 		respondWithError(w, http.StatusBadRequest, err.Error() + " line : 65")
		return
 	}
 	if r.Form.Get("Title") == ""{
		music.Title = "Untitled"
	}else{
		music.Title = r.Form.Get("Title")
	}
 	music.AudioPath = filePath[2:]
 	music.Duration 	=  durationStr[14:len(durationStr)-2]
 	music.UserId 	= user.Id 
	now 			:= time.Now().UTC()
	music.CreatedAt = &now
	if err = database.GetDB().Create(&music).Error; err!= nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return	
	}
	//TODO:Process music for rating
	//return user rating for this music
	//Example below
	var rating = m.UserRating{UserId : user.Id, MusicId : music.OriginalMusicId}
	var ratings []m.UserRating
	rand.Seed(time.Now().UnixNano())
	if database.GetDB().Where(&rating).Find(&rating).RecordNotFound(){
		log.Println(rating)
		rating.Rating = uint(rand.Intn(1014))
		database.GetDB().Save(&rating)
	}else{
		log.Println(rating)
		rating.Rating = uint(rand.Intn(1014))
		database.GetDB().Model(&m.UserRating{}).Where("id = ?", rating.Id).Updates(&m.UserRating{Rating: rating.Rating})
	}
	database.GetDB().Where(&m.UserRating{MusicId : music.OriginalMusicId}).Find(&ratings)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"ratings":ratings})
}
func UploadMediaFilesForAdmin(w http.ResponseWriter, r *http.Request){
	err := r.ParseMultipartForm(256 << 20)//, ~100 MByte
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return	
	}
	user, err := Authorize(r.Header.Get("Authorization"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"You are not authorized!")
		return
	}
	if !isAdmin(user) {
		respondWithError(w, http.StatusBadRequest,"You cannot add genres, it requires admin priviliges")
		return	
	}
	music := m.Music{OriginalMusicId : 0}
	i64, err := strconv.ParseInt(r.Form.Get("GenreID"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	music.GenreId = uint(i64)
	if  music.GenreId == 0{
		respondWithError(w, http.StatusBadRequest, "Genre is not defined")
		return	
	}
	var title,author = false,false
	if r.Form.Get("Title") == ""{
		title = true
	}else{
		music.Title = r.Form.Get("Title") 
	}
	if r.Form.Get("Author") == "" || r.Form.Get("Author") == "0"{
		author = true
	}else{
		if i64, err = strconv.ParseInt(r.Form.Get("Author"), 10, 64); err != nil{
			respondWithError(w, http.StatusBadRequest, "Author is not defined")
			return	
		}
		music.AuthorId = uint(i64)
	}
	filePath, err:= processMediaFile(r, "audio")
 	if err != nil{
 		respondWithError(w, http.StatusBadRequest, err.Error())
		return
 	}
 	filePath, durationStr, authorId, titleStr, err := Mp3ConverterAndDuration(filePath,author)
 	if err != nil{
 		respondWithError(w, http.StatusBadRequest, err.Error())
		return
 	}
	videoFilePath, err:= processMediaFile(r, "video")
 	if err != nil{
 		respondWithError(w, http.StatusBadRequest, err.Error())
		return
 	}
 	music.AudioPath = filePath[2:]
 	music.Duration 	=  durationStr[14:len(durationStr)-2]
 	music.UserId 	= user.Id 
	now 			:= time.Now().UTC()
	music.CreatedAt = &now
	if title{
		music.Title = titleStr
	}
	if author{
		music.AuthorId = authorId
	}
	music.VideoPath = videoFilePath[2:]
	if err = database.GetDB().Create(&music).Error; err!= nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return	
	}
	if newMusics.Len() > newMusicsSize{
		newMusics.Remove(newMusics.Front())
	}
	newMusics.PushBack(music)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"music":music})
}
func AddGenre(w http.ResponseWriter, r *http.Request){
	user, err := Authorize(r.Header.Get("Authorization"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"You are not authorized!")
		return
	}
	if !isAdmin(user) {
		respondWithError(w, http.StatusBadRequest,"You cannot add genres, it requires admin priviliges")
		return	
	}
	genre :=  m.Genre{Name : r.URL.Query().Get("name")}
	if genre.Name == ""{
		respondWithError(w, http.StatusBadRequest, "Please specify genre")
		return	
	}
	if database.GetDB().Where(&genre).Find(&genre).RecordNotFound(){
		database.GetDB().Save(&genre)
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message": OK,"genre":genre})
}
func GetGenres(w http.ResponseWriter, r *http.Request){
	var genres []m.Genre
	database.GetDB()/*.Preload("Musics")*/.Find(&genres)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"genres":genres})	
}
//to get original music for karaoke, not user performed!
func GetMusicByGenreAndTitle(w http.ResponseWriter, r *http.Request){
	var musics []m.Music
	var tag string
	if r.URL.Query().Get("genre_id") != ""{
		tag = "genre_id"
	}else if r.URL.Query().Get("title") != ""{
		tag = "title"
	}else{
		database.GetDB().Where("original_music_id = 0",).Find(&musics)	
		respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"musics":musics})
		return
	}
	database.GetDB().Where(tag + " = ? AND original_music_id = 0",r.URL.Query().Get(tag)).Find(&musics)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"musics":musics})
}
//to get any music in db, user performed, original and etc
func GetMusicByID(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var music m.Music
	database.GetDB().Where("id = ?",params[c.ID]).Find(&music)
	if music.Id != 0{
		respondWithJSON(w, http.StatusOK, map[string]interface{}{"message":OK,"music":music})
	}else{
		respondWithError(w, http.StatusBadRequest, "RecordNotFound")
	}
}
//get list of music saved in list newMusics
func GetNewestMusics(w http.ResponseWriter, r *http.Request){
	var musics []m.Music
	for e := newMusics.Front(); e != nil; e = e.Next() {
		 musics = append(musics, e.Value.(m.Music))
	}
	respondWithJSON(w,http.StatusOK,map[string]interface{}{"message":OK, "musics":musics})
}
//get all music that user has added (as karaoke)
func GetMyMusics(w http.ResponseWriter, r *http.Request){
	user, err := Authorize(r.Header.Get("Authorization"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"You are not authorized!")
		return
	}
	var musics []m.Music
	database.GetDB().Where("user_id = ? AND original_music_id != 0",user.Id).Find(&musics)
	respondWithJSON(w,http.StatusOK,map[string]interface{}{"message":OK,"musics":musics})	
}
func GetScoreBoardForMusic(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var userRatings []m.UserRating
	database.GetDB().Where("music_id = ?",params[c.ID]).Find(&userRatings)
	respondWithJSON(w,http.StatusOK, map[string]interface{}{"message":OK,"rating":userRatings})
}
func GetUserByID(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var user  m.User
	database.GetDB().Where("id = ?",params[c.ID]).Take(&user)
	respondWithJSON(w,http.StatusOK,map[string]interface{}{"message":OK,"user":user})	
}
func SeachMusicByTitle(w http.ResponseWriter, r *http.Request){
	title := r.URL.Query().Get("title")
	var titles []string
	//rows, err := database.GetDB().Raw("SELECT title, MATCH (title) AGAINST ('" + title + "') FROM musics;").Rows() // (*sql.Rows, error)
	rows, err := database.GetDB().Raw("SELECT title FROM musics Where title LIKE '%" + title + "%' and original_music_id = 0;").Rows()
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t string
  		rows.Scan(&t)
  		titles = append(titles, t)	
	}
	respondWithJSON(w,http.StatusOK, map[string]interface{}{"message":OK,"titles":titles})
}
func GetAuthors(w http.ResponseWriter, r *http.Request){
	var authors []m.Author
	database.GetDB().Find(&authors)
	respondWithJSON(w,http.StatusOK, map[string]interface{}{"message":OK, "authors":authors})
}
func GetAlbumsByAuthor(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var albums []m.Album
	database.GetDB().Where("author_id = ?",params[c.ID]).Find(&albums)
	for i := 0; i < len(albums); i++{
		database.GetDB().Where(&m.Music{AlbumId : albums[i].Id}).Find(albums[i].Musics)
	}
	respondWithJSON(w,http.StatusOK, map[string]interface{}{"message":OK, "albums":albums})
}
func isAudio(ext string)bool{
	if ext == ".MP3" || ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".gsm" || ext == ".dct" || ext == ".flac" || ext == ".au" || ext == ".aiff" || ext == ".vox" || ext == ".raw"{
		return true
	}
	return false
}
func isVideo(ext string)bool{
	if ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".flv" || ext == ".wmv"{
		return true
	}
	return false
}
func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)
	randomness := make([]byte, n)
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)
	for pos := range output {
		random := uint8(randomness[pos])
		randomPos := random % uint8(l)
		output[pos] = letterBytes[randomPos]
	}
	return output
}
func Mp3ConverterAndDuration(filePath string, author bool)(string, string, uint, string, error){
	var err error
	err = nil
	if filepath.Ext(filePath) != ".mp3"{
 		cmd := exec.Command("ffmpeg","-i",filePath,"-acodec", "libmp3lame", filePath + ".mp3")     //ffmpeg -i besame.wav -acodec libmp3lame besame.wav.mp3	
		if err = cmd.Run(); err != nil{
			log.Println(err.Error() + "309")
			return "","",0,"",err
		}
		filePath = filePath + ".mp3"
 	}
 	cmd := exec.Command("ffmpeg","-i",filePath,"-f","null","-")     
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb	
	if err = cmd.Run(); err != nil{
		log.Println(err.Error() + "319")
		return "","", 0, "",err
	}
	re := regexp.MustCompile(`Duration: .([0-9]|:)*`)
	durationStr := fmt.Sprintf("%q\n", re.Find([]byte(errb.String())))//TODO:	check if durationStr has correct form "Duration: 00:23:23.23"
	f,err := os.Open(filePath)
	if err != nil{
		return "","", 0, "",err	
	}
	meta, err := tag.ReadFrom(f)
	if err != nil {
		log.Println("ReadMetadata error")
		log.Println(err.Error())
		log.Println(err.Error() + "332")
		return "","", 0, "",err	
	}
	var authorObj m.Author
	if author == true{
		authorObj = m.Author{Name : meta.Artist()}
		if database.GetDB().Where(&authorObj).Find(&authorObj).RecordNotFound(){
			database.GetDB().Save(&authorObj)
		}
	}
	return filePath, durationStr,authorObj.Id, meta.Title(), nil
}
func processMediaFile(r *http.Request, fileType string)(string,error){
	var maxFileSize uint
	_, fileHeader, err := r.FormFile(fileType)
	if err != nil{
		return "",err
	}
	var filePath string
	if fileType == "audio"{
		 if isAudio(filepath.Ext(fileHeader.Filename)){
			filePath = c.AudioPath + "/" + string(RandASCIIBytes(36)) + filepath.Ext(fileHeader.Filename)
			maxFileSize = 3*10E7//30MB
		}else{
			return "",errors.New("File format is not supported, extension is "+filepath.Ext(fileHeader.Filename))
		}	
	}else if fileType == "video"{
		if isVideo(filepath.Ext(fileHeader.Filename)){
			filePath = c.VideoPath + "/" + string(RandASCIIBytes(36)) + filepath.Ext(fileHeader.Filename)
			maxFileSize = 6*10E8//600MB
		}else{
			return "",errors.New("File format is not supported, extension is "+filepath.Ext(fileHeader.Filename))
		}	
	}else{
		return "",errors.New("file type: " + fileType)
	}
	if uint(fileHeader.Size) > maxFileSize{//30MB
		log.Println("FILE SIZE")
		log.Println(uint(fileHeader.Size))
		return "",errors.New("file size limit exceeded")
	}
	file, err1 := fileHeader.Open()
 	defer file.Close()
 	out,  err2 := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
 	defer out.Close()
 	_, err3 := io.Copy(out, file) 
 	if err1 != nil || err2 != nil || err3 != nil{
		return	"", errors.New("Error in processing " + fileType + " file,try again!")
 	}
 	return filePath, nil
}