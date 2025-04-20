package constants

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILED = "failed"

	// USER
	USER_FOUND_SUCCESS = "Successfully fetch user data"
	USER_NOT_FOUND = "User not found"
	USER_UPDATE_SUCCESS = "Successfully update user data"

	// ACTIVITY & PHOTO
	CREATE_ACTIVITY_SUCCESS = "Success to create activity"
	CREATE_ACTIVITY_SUCCESS_WITH_WARNING = "Success to create activity. Failed to retrieve activity data."
	CREATE_ACTIVITY_FAILED = "Failed to create activity"
	ACTIVITY_NOT_FOUND = "Activity not found"
	ACTIVITY_FOUND_SUCCESS = "Successfully fetch activity data"
	PHOTO_FILE_REQUIRED = "Photo file required"
	SAVE_PHOTO_METADATA_FAILED = "Failed to save picture metadata"
	UPLOAD_FILE_FAILED = "Failed to upload file"
	UPLOAD_FILE_SUCCESS = "Success to upload file"

	// AUTH
	LOGIN_SUCCESS = "Successfully login"
	UNAUTHORIZED = "Unauthorized"
	MISSING_ID_TOKEN = "Missing ID token"

	INVALID_TYPE_USER_ID = "Invalid type: userID"
	INVALID_TYPE_ACTIVITY_ID = "Invalid type: activityID"
	INVALID_REQUEST = "Invalid request"

)