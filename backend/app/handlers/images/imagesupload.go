package images

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

// i have to add messages folder too
//also add the jpeg, png, gif validation

//image directory
const baseUploadDir = "./public/images"

// ImageType defines which kind of image is being uploaded
type ImageType string

const (
	PostsImage    ImageType = "post"
	AvatarImage   ImageType = "avatar"
	MessagesImage ImageType = "message"
	CommentsImage ImageType = "comment"
)

// ImageUploadHandler handles image uploads for posts, avatars, messages, and comments
func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	
	imageTypeStr := r.FormValue("image_type")
	var imageType ImageType
	switch imageTypeStr {
	case "post":
		imageType = PostsImage
	case "avatar":
		imageType = AvatarImage
	case "message":
		imageType = MessagesImage
	case "comment":
		imageType = CommentsImage
	default:
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Basic validation
	if !strings.HasPrefix(handler.Header.Get("Content-Type"), "image/") {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Map type → folder name (with 's' at the end)
	folderMap := map[ImageType]string{
		PostsImage:    "posts",
		AvatarImage:   "avatars",
		MessagesImage: "messages",
		CommentsImage: "comments",
	}

	uploadDir := path.Join(baseUploadDir, folderMap[imageType])
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0o755); err != nil {
			http.Error(w, "Could not create directory", http.StatusInternalServerError)
			return
		}
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s_%s", uuid.New().String(), handler.Filename)
	filepath := path.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	// Return a JSON with relative path for frontend
	imageURL := fmt.Sprintf("/images/%s/%s", folderMap[imageType], filename)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"image_url": imageURL})
}
