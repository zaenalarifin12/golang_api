package entity

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"type:varchar(255)" json:"email"`
	Name     string	`gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string	`gorm:"->;<-;not null" json:"-"`
	Token    string `gorm:"-" json:"token,omitempty"`
}
