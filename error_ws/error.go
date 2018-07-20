package error_ws

const (
	MessageSendError = "消息发送错误"
	RoomEnterError = "房间进入错误"
	MultiEnterRoomError = "房间重复进入错误"
)

type Errorws struct {
	errorMessageMap map[string]string
}

func New() *Errorws {
	return &Errorws{errorMessageMap:map[string]string{
		"0001":MessageSendError,
		"0002":RoomEnterError,
		"0003":MultiEnterRoomError,
	}}
}
var std = New()

func Errormessagegenerate(code string) string {
	if err, ok := std.errorMessageMap[code]; ok{
		return err
	} else {
		return code
	}
}
