package entity

import (
	"strconv"
	"strings"
	"time"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleName    string         `gorm:"type:varchar(50);not null"`
	Description string         `gorm:"type:varchar(255);null"`
	IsActive    bool           `gorm:"type:bool;not null;default:TRUE"`
	Permission  datatypes.JSON `gorm:"null"`
}

func (r *Role) FromProto(role *pb.Role) error {
	if role.Id != "" {
		id, err := strconv.Atoi(role.Id)
		if err != nil {
			return err
		}
		r.ID = uint(id)
	}
	marshaller := jsonpb.Marshaler{}
	jsonString, err := marshaller.MarshalToString(role.Permission)
	if err != nil {
		return err
	}

	r.CreatedAt = time.Now()
	r.IsActive = role.IsActive
	r.RoleName = role.RoleName
	r.Description = role.Description
	r.Permission = datatypes.JSON(jsonString)

	return nil
}

func (r *Role) ToProto() *pb.Role {
	var unmarshaller jsonpb.Unmarshaler
	protoStruct := &structpb.Struct{}
	if err := unmarshaller.Unmarshal(strings.NewReader(r.Permission.String()), protoStruct); err != nil {
		protoStruct = nil
	}

	return &pb.Role{
		Id:          strconv.Itoa(int(r.ID)),
		IsActive:    r.IsActive,
		RoleName:    r.RoleName,
		Description: r.Description,
		Permission:  protoStruct,
	}
}
