package store

import "time"

type BaseModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Approve struct {
	BaseModel
	Approver   string `json:"approver"`
	Reason     string `json:"reason"`
	IsApproved bool   `json:"is_approved"`
}

type Alarm struct {
	BaseModel
	Type       string  `json:"Type"`      //垃圾类型
	FaceUser   string  `json:"face_user"` //AI 识别用户
	Message    string  `json:"message"`   //通知信息
	Address    string  `json:"address"`   //监控地点
	City       string  `json:"city"`
	Similarity float32 `json:"similarity"` //相似度
}

type Device struct {
	DeviceType string `json:"device_type"` //设备型号
	City       string `json:"city"`
	Total      int    `json:"total"`
}

type AlarmVideo struct {
	BaseModel
	UID   string    `json:"uid"` //视频标识符？
	Begin time.Time `json:"begin"`
	End   time.Time `json:"end"`
}
