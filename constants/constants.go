package constants

const (
	BaseURL  = "https://www.premiumize.me/api/"
	AuthURL  = "https://www.premiumize.me/authorize"
	TokenURL = "https://www.premiumize.me/token"

	endpointFolder           = "folder/"
	EndpointFolderList       = BaseURL + endpointFolder + "list"
	EndpointFolderCreate     = BaseURL + endpointFolder + "create"
	EndpointFolderRename     = BaseURL + endpointFolder + "rename"
	EndpointFolderPaste      = BaseURL + endpointFolder + "paste"
	EndpointFolderDelete     = BaseURL + endpointFolder + "delete"
	EndpointFolderUploadInfo = BaseURL + endpointFolder + "uploadinfo"
	EndpointFolderSearch     = BaseURL + endpointFolder + "search"

	endpointItem        = "item/"
	EndpointItemListAll = BaseURL + endpointItem + "listall"
	EndpointItemDelete  = BaseURL + endpointItem + "delete"
	EndpointItemRename  = BaseURL + endpointItem + "rename"
	EndpointItemDetails = BaseURL + endpointItem + "details"

	endpointTransfer              = "transfer/"
	EndpointTransferCreate        = BaseURL + endpointTransfer + "create"
	EndpointTransferDirectDL      = BaseURL + endpointTransfer + "directdl"
	EndpointTransferList          = BaseURL + endpointTransfer + "list"
	EndpointTransferClearFinished = BaseURL + endpointTransfer + "clearfinished"
	EndpointTransferDelete        = BaseURL + endpointTransfer + "delete"

	endpointAccount     = "account/"
	EndpointAccountInfo = BaseURL + endpointAccount + "info"

	endpointZip         = "zip/"
	EndpointZipGenerate = BaseURL + endpointZip + "generate"

	endpointCache      = "cache/"
	EndpointCacheCheck = BaseURL + endpointCache + "check"

	endpointServices     = "services/"
	EndpointServicesList = BaseURL + endpointServices + "list"

	HeaderUserAgent       = "Dart/2.10 (dart:io)" // Official client XD
	HeaderContentTypeForm = "application/x-www-form-urlencoded;charset=utf-8"

	ClientID          = "284888557"
	TokenResponseType = "device_code"
)
