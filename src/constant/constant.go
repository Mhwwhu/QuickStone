package constant

const CtxUserIdKey string = "UserId"
const CtxUserNameKey string = "UserName"
const StatusCodeKey string = "status_code"
const StatusMsgKey string = "status_msg"

const MetadataVarPrefix = "metadata"

const (
	ObjectStorageExchange = "obj-storage"
	ObjectEventPrefix     = "object"

	ObjectStoredEvent          = "object.stored"
	ObjectSoftDeletedEvent     = "object.soft-deleted"
	ObjectPhysicalDeletedEvent = "object.physical-deleted"

	ObjectMetaQueue            = "obj-storage.meta.q"
	ObjectOnUploadProcessQueue = "obj-storage.on-upload-process.q"
)

const ()
