package consts

const (
	StatusOK                  int32 = 200
	StatusClientParamsInvalid int32 = 400 // invalid client parameters 客户端参数有问题，可能是用户操作不当
	StatusGPTError            int32 = 500 // GPT API Error, caused by network or server mostly GPT 报错。最可能的情况是 OpenAI 的接口临时性的问题（网络、服务器问题等）
	StatusRetry               int32 = 501 // service busy, try again later 服务繁忙，稍后重试
)
