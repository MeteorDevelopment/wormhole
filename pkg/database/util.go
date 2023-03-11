package database

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

func encodeUUID(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	return w.WriteString(v.Interface().(uuid.UUID).String())
}

func decodeUUID(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	str, err := r.ReadString()
	if err != nil {
		return err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(id))
	return nil
}
