package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	c "main/constants"
	contr "main/controllers"
	"main/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	//"runtime"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"

func main() {

	log.Println("karaoke main server. Version 1.0.0")
	log.Println("Loading environment variables.")
	e := godotenv.Load(".env") //Load .env file
	if e != nil {
		log.Fatal(e)
	}

	port := os.Getenv("app_port")
	if port == "" {
		log.Fatal("$app_port not set")
	} else {
		log.Println(fmt.Sprintf("$app_port: %s", port))
	}

	_, err2 := database.Connect() //Load database
	if err2 != nil {
		log.Fatal(err2)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	os.Mkdir(dir+c.VideoPath[1:], os.ModePerm)
	os.Mkdir(dir+c.AudioPath[1:], os.ModePerm)
	
	log.Println("Setting router and controllers.")
	router := mux.NewRouter()
	router.PathPrefix(c.VideoPath[1:]+"/").Handler(http.StripPrefix(c.VideoPath[1:]+"/", http.FileServer(http.Dir(c.VideoPath))))
	router.PathPrefix(c.AudioPath[1:]+"/").Handler(http.StripPrefix(c.AudioPath[1:]+"/", http.FileServer(http.Dir(c.AudioPath))))

	router.HandleFunc(c.RegisterURI, contr.Register).Methods(http.MethodPost)
	router.HandleFunc(c.VerifyAuthCodeAndLogin, contr.VerifyAuthCodeAndLogin).Methods(http.MethodPost)
	router.HandleFunc(c.UploadMediaFileForUser, contr.UploadMusicForUser).Methods(http.MethodPost)
	router.HandleFunc(c.UploadMediaFileForAdmin, contr.UploadMediaFilesForAdmin).Methods(http.MethodPost)
	router.HandleFunc(c.AdminGenres, contr.AddGenre).Methods(http.MethodPost)
	router.HandleFunc(c.Genres, contr.GetGenres).Methods(http.MethodGet)
	router.HandleFunc(c.GetMusics, contr.GetMusicByGenreAndTitle).Methods(http.MethodGet)
	router.HandleFunc(c.GetMusicByID, contr.GetMusicByID).Methods(http.MethodGet)
	router.HandleFunc(c.GetMyMusics, contr.GetMyMusics).Methods(http.MethodGet)
	router.HandleFunc(c.GetNewestMusics, contr.GetNewestMusics).Methods(http.MethodGet)
	router.HandleFunc(c.GetScoreBoardForMusic, contr.GetScoreBoardForMusic).Methods(http.MethodGet)
	router.HandleFunc(c.GetUser, contr.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc(c.SeachMusicByTitle, contr.SeachMusicByTitle).Methods(http.MethodGet)
	router.HandleFunc(c.Authors, contr.GetAuthors).Methods(http.MethodGet)
	router.HandleFunc(c.AuthorAlbums, contr.GetAlbumsByAuthor).Methods(http.MethodGet)
	log.Println("Router and controllers set successfully.")
	log.Println("Maintaining web application...")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  15 * time.Second}
	log.Fatal(server.ListenAndServe())
}
