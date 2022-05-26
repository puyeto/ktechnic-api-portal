package models

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID              uint32              `gorm:"primary_key;auto_increment" json:"id"`
	Nickname        string              `gorm:"size:255;not null;unique" json:"nickname"`
	Email           string              `gorm:"size:100;not null;unique" json:"email"`
	Phone           string              `gorm:"size:100;not null" json:"phone"`
	Password        string              `gorm:"size:100;not null;" json:"password,omitempty"`
	CompanyID       uint32              `gorm:"not null;" json:"company_id"`
	RoleID          int                 `gorm:"not null;" json:"role_id,omitempty"`
	UpdatedBy       uint32              `json:"updated_by"`
	Token           string              `gorm:"-" json:"token,omitempty"`
	RoleName        string              `gorm:"-" json:"role_name,omitempty"`
	CreatedAt       time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Company         CompanyShortDetails `gorm:"-" json:"company"`
	PermissionField []interface{}       `gorm:"-" json:"permission_field"`
	Permissions     map[string]int      `gorm:"-" json:"permissions,omitempty"`
}

// Hash ...
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword ...
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave ...
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare ...
func (u *User) Prepare() {
	u.Nickname = html.EscapeString(strings.TrimSpace(strings.ToUpper(u.Nickname)))
	u.Email = html.EscapeString(strings.TrimSpace(strings.ToLower(u.Email)))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate ...
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Phone == "" {
			return errors.New("Required Phone")
		}
		// if u.Password == "" {
		// 	return errors.New("Required Password")
		// }
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if u.Phone == "" {
			return errors.New("Required Phone")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// Count Users ...
func (u *User) CountUsers(db *gorm.DB, roleid uint32) int {
	var count int
	query := db.Debug().Model(u)

	if roleid == 1002 {
		query = query.Where("company_id = ?", u.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", u.UpdatedBy)
	}

	query.Count(&count)
	return count
}

// ListUsers Get all users...
func (u *User) ListUsers(db *gorm.DB, roleid uint32, offset, limit int) (*[]User, error) {
	var err error
	users := []User{}

	query := db.Debug().Select("users.id, users.nickname, users.email, users.phone, users.company_id, users.role_id, users.created_at, users.updated_at, users.updated_by, role_name")

	if roleid == 1002 {
		query = query.Where("company_id = ?", u.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", u.UpdatedBy)
	}

	if err = query.Offset(offset).Limit(limit).Joins("left join roles on roles.id = users.role_id").Order("users.id ASC").Find(&users).Error; err != nil {
		return &users, err
	}

	if len(users) > 0 {
		for i := range users {
			db.Debug().Table("companies").Model(&Companies{}).Where("id = ?", users[i].CompanyID).Take(&users[i].Company)
		}
	}
	return &users, err
}

// FindUserByID ...
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(u).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return u, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return u, errors.New("User Not Found")
	}
	return u, err
}

// SaveUser ...
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return u, err
	}

	err := tx.Debug().Create(&u).Error
	if err != nil {
		tx.Rollback()
		return u, err
	}

	if u.ID > 0 {
		u.Permissions["added_by"] = int(u.UpdatedBy)
		query := u.structureQuery()
		fmt.Println(query)
		if err := tx.Debug().Exec(query).Error; err != nil {
			tx.Rollback()
			return u, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return u, err
	}

	return u, nil
}

func (u *User) structureQuery() string {
	userid := strconv.Itoa(int(u.ID))
	query := "INSERT INTO `user_permissions` ("
	cols := ""
	rows := " VALUES ("
	for col, row := range u.Permissions {
		cols += col + ", "
		rows += strconv.Itoa(row) + ", "
	}
	return query + cols + "user_id, created_at)" + rows + userid + ", '" + time.Now().Format("2006-01-02 15:04:05") + "')"
}

// UpdateAUser ...
func (u *User) UpdateAUser(db *gorm.DB) (*User, error) {

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &User{}, err
	}

	// To hash the password
	// if u.Password == "" {
	// 	err := u.BeforeSave()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	if err := tx.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"nickname":   u.Nickname,
			"email":      u.Email,
			"phone":      u.Phone,
			"role_id":    u.RoleID,
			"update_at":  time.Now(),
			"updated_by": u.UpdatedBy,
		},
	).Error; err != nil {
		tx.Rollback()
		return &User{}, err
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&u).Error
	if err != nil {
		tx.Rollback()
		return &User{}, err
	}

	query := u.structureUpdateQuery()
	if err := tx.Debug().Exec(query).Error; err != nil {
		tx.Rollback()
		return &User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &User{}, err
	}

	return u, nil
}

func (u *User) structureUpdateQuery() string {
	userid := strconv.Itoa(int(u.ID))
	query := "UPDATE `user_permissions` SET "
	for col, row := range u.Permissions {
		query += col + "=" + strconv.Itoa(row) + ", "
	}

	return query + "updated_at='" + time.Now().Format("2006-01-02 15:04:05") + "' WHERE user_id=" + userid
}

// DeleteAUser ...
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
