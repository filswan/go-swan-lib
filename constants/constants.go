package constants

const (
	EMPTY_STRING         = ""
	DEFAULT_SELECT_LIMIT = "100"

	STORAGE_SERVER_TYPE_WEB_SERVER  = "web server"
	STORAGE_SERVER_TYPE_IPFS_SERVER = "ipfs server"

	SWAN_API_STATUS_SUCCESS = "success"
	SWAN_API_STATUS_FAIL    = "fail"

	TASK_TYPE_VERIFIED = "verified"
	TASK_TYPE_REGULAR  = "regular"

	TASK_STATUS_ASSIGNED              = "Assigned"
	TASK_STATUS_DEAL_SENT             = "DealSent"
	TASK_STATUS_PROGRESS_WITH_FAILURE = "ProgressWithFailure"

	TASK_FAST_RETRIEVAL = 1

	TASK_BID_MODE_AUTO   = 1
	TASK_BID_MODE_MANUAL = 0

	TASK_IS_PUBLIC = 1

	EPOCH_PER_HOUR = 120

	PATH_TYPE_NOT_EXIST = 0 //this path not exists
	PATH_TYPE_FILE      = 1 //file
	PATH_TYPE_DIR       = 2 //directory
	PATH_TYPE_UNKNOWN   = 3 //unknown path type

	JSON_FILE_NAME_BY_CAR    = "car.json"
	JSON_FILE_NAME_BY_GOCAR  = "car.json"
	JSON_FILE_NAME_BY_UPLOAD = "car.json"
	JSON_FILE_NAME_BY_TASK   = "-metadata.json"
	JSON_FILE_NAME_BY_DEAL   = "-deals.json"
	JSON_FILE_NAME_BY_AUTO   = "car.json"

	CSV_FILE_NAME_BY_CAR    = "car.csv"
	CSV_FILE_NAME_BY_GOCAR  = "car.csv"
	CSV_FILE_NAME_BY_UPLOAD = "car.csv"
	CSV_FILE_NAME_BY_TASK   = "-metadata.csv"
	CSV_FILE_NAME_BY_DEAL   = "-deals.csv"
	CSV_FILE_NAME_BY_AUTO   = "car.csv"

	AuthorizationHeaderKey = "Authorization"

	TASK_SOURCE_ID_DEFAULT      = 0
	TASK_SOURCE_ID_SWAN         = 1
	TASK_SOURCE_ID_SWAN_CLIENT  = 2
	TASK_SOURCE_ID_SWAN_FS3     = 3
	TASK_SOURCE_ID_SWAN_PAYMENT = 4
	TASK_SOURCE_ID_OTHER        = 5

	LOTUS_PRICE_MULTIPLE = 1e18

	WALLET_NON_VERIFIED_MESSAGE = "Not a Verified Client"
)
