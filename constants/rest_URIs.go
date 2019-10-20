package constants
const (
	VarifyAuthCode = "https://graph.accountkit.com/v1.3/access_token?grant_type=authorization_code"
	//https://graph.accountkit.com/v1.3/access_token?grant_type=authorization_code&code=<authorization_code>&access_token=AA|354536525217067|8042b4b2bb00245c7b034e330f220848
	GetPhoneByToken = "https://graph.accountkit.com/v1.3/me/"
	//https://graph.accountkit.com/v1.3/me/?access_token=<access_token>
	FacebookAppID = "354536525217067"
	FacebookSecretKey = "8042b4b2bb00245c7b034e330f220848"

)
const (
	// host/api/v1
	ApiURI = "/api/v1"
	RegisterURI				= ApiURI + "/register"
	VerifyAuthCodeAndLogin 	= ApiURI + "/auth"
	UploadMediaFileForUser	= ApiURI + "/upload-media-file"
	UploadMediaFileForAdmin	= ApiURI + "/admin/upload-media-file"
	AdminGenres				= ApiURI + "/admin/genres"
	Genres                  = ApiURI + "/genres"
	GetMusics 		        = ApiURI + "/musics"
	GetMyMusics				= ApiURI + "/me/musics"
	GetMusicByID			= GetMusics + "/{" + ID + "}"
	GetNewestMusics 		= ApiURI + "/newest/musics" 
	GetScoreBoardForMusic	= GetMusics + "/{" + ID + "}" + "/rating"
	Users                   = ApiURI    + "/users"
	GetUser                 = Users     + "/{" + ID + "}"
	SeachMusicByTitle       = ApiURI    + "/search"
	Authors                 = ApiURI + "/authors"
	AuthorAlbums            = Authors  + "/{" + ID + "}"
	//personal
	// host/api/v1/register
	//FreeTimesURI = ApiURI + "/free-times"                  // params: zone-id 	//done
	//FreeTimeURI  = FreeTimesURI + "/{" + UUIDPathVar + "}" // done
)

//path variables
const (
	ID = "id"
	Title = "title"
)

const(
	Root = "."
	AudioPath = Root + "/audios"
	VideoPath = Root + "/videos"
)