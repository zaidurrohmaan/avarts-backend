package constants

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILED = "failed"

	// USER
	USER_FOUND_SUCCESS = "successfully fetch user data"
	USER_NOT_FOUND = "user not found"
	USER_UPDATE_SUCCESS = "successfully update user data"

	// ACTIVITY & PHOTO
	CREATE_ACTIVITY_SUCCESS = "success to create activity"
	CREATE_ACTIVITY_SUCCESS_WITH_WARNING = "success to create activity. Failed to retrieve activity data."
	CREATE_ACTIVITY_FAILED = "failed to create activity"
	ACTIVITY_NOT_FOUND = "activity not found"
	ACTIVITY_FOUND_SUCCESS = "successfully fetch activity data"
	PHOTO_FILE_REQUIRED = "photo file required"
	SAVE_PHOTO_METADATA_FAILED = "failed to save picture metadata"
	UPLOAD_FILE_FAILED = "failed to upload file"
	UPLOAD_FILE_SUCCESS = "success to upload file"

	// AUTH
	LOGIN_SUCCESS = "successfully login"
	UNAUTHORIZED = "unauthorized"
	MISSING_ID_TOKEN = "missing ID token"

	// LIKE
	CREATE_LIKE_SUCCESS = "successfully create like"
	LIKE_ALREADY_EXISTS = "like already exists"
	LIKE_NOT_FOUND = "like not found"
	LIKE_DELETED = "like deleted"

	// COMMENT
	CREATED_COMMENT_SUCCESS = "successfully create comment"
	CREATED_COMMENT_FAILED = "failed to create comment"

	INVALID_TYPE_USER_ID = "invalid type: userID"
	INVALID_TYPE_ACTIVITY_ID = "invalid type: activityID"
	INVALID_REQUEST = "invalid request"
)