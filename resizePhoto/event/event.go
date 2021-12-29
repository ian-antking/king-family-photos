package event

type Bucket struct {
	Name string `json:"name,omitempty"`
}

type S3 struct {
	Bucket Bucket `json:"bucket,omitempty"`
	Object Object `json:"object,omitempty"`
}

type Object struct {
	Key string `json:"key,omitempty"`
}

type Record struct {
	S3 S3 `json:"s3,omitempty"`
}

type Message struct {
	Records []Record `json:"Records,omitempty"`
}
