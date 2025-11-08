package bucket

type bucketACLTypeUtil struct{}
type storageTypeUtil struct{}

var BucketACLTypeUtil bucketACLTypeUtil
var StorageTypeUtil storageTypeUtil

func (bucketACLTypeUtil) FromString(str string) BucketACLType {
	switch str {
	case "private":
		return BucketACLType_PRIVATE
	case "public_read":
		return BucketACLType_PUBLIC_READ
	case "public":
		return BucketACLType_PUBLIC
	}
	return -1
}

func (bucketACLTypeUtil) ToString(acl BucketACLType) string {
	switch acl {
	case BucketACLType_PRIVATE:
		return "private"
	case BucketACLType_PUBLIC_READ:
		return "public_read"
	case BucketACLType_PUBLIC:
		return "public"
	}
	return "unknown"
}

func (storageTypeUtil) FromString(str string) StorageType {
	switch str {
	case "standard":
		return StorageType_STANDARD
	case "low_freq":
		return StorageType_LOW_FREQ
	}
	return -1
}

func (storageTypeUtil) ToString(storage_t StorageType) string {
	switch storage_t {
	case StorageType_STANDARD:
		return "standard"
	case StorageType_LOW_FREQ:
		return "low_freq"
	}
	return "unknown"
}
