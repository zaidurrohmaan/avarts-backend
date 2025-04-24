package constants

const (
	// Status
	StatusSuccess = "success"
	StatusFailed  = "failed"

	// User
	UserFetchSuccess  = "user data fetched"
	UserNotFound      = "user not found"
	UserUpdateSuccess = "user data updated"
	UsernameIsTaken   = "username already taken"

	// Activity & Photo
	ActivityCreateSuccess            = "activity created"
	ActivityCreateSuccessWithWarning = "activity created, but failed to fetch activity data"
	ActivityCreateFailed             = "failed to create activity"
	ActivityNotFound                 = "activity not found"
	ActivityFetchSuccess             = "activity data fetched"
	PhotoFileRequired                = "photo file is required"
	PhotoMetadataSaveFailed          = "failed to save photo metadata"
	FileUploadFailed                 = "failed to upload file"
	FileUploadSuccess                = "file uploaded"
	OpenFileFailed                   = "failed to open file"
	ReadFileFailed                   = "failed to read file"
	InvalidImage                     = "file is not a valid image. allowed extensions: .jpg, .jpeg, .png"
	FileSizeExceeded                 = "file size exceeds the maximum size"
	DeleteActivitySuccess            = "activity deleted"
	ActivityDeleteAccessDenied       = "access denied: you are not allowed to delete this activity"

	// Auth
	LoginSuccess   = "logged in"
	Unauthorized   = "unauthorized access"
	MissingIDToken = "missing ID token"

	// Like
	LikeCreateSuccess = "like created"
	LikeAlreadyExists = "like already exists"
	LikeNotFound      = "like not found"
	LikeDeleted       = "like deleted"

	// Comment
	CommentCreateSuccess            = "comment created"
	CommentCreateFailed             = "failed to create comment"
	CommentDeleteSuccess            = "comment deleted"
	CommentDeleteFailed             = "failed to delete comment"
	CommentDeleteAccessDenied       = "access denied: you are not allowed to delete this comment"
	CommentCreateSuccessWithWarning = "comment created, but failed to fetch full data"
	CommentNotFound                 = "comment not found"

	// Validation / Error
	InvalidRequestFormat = "invalid request format"
)
