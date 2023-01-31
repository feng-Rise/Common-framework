package serialize

//用不同的序列化协议实现body数据的编码,支持不同序列化协议,没实现这个功能
type Serializer interface {
	// Code 协议里面对应字段的值
	Code() byte
	// Encode 编码
	Encode(val interface{}) ([]byte, error)
	// Decode 解码
	Decode(data []byte, val interface{}) error
}
