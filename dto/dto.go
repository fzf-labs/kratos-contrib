package dto

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Copy copier的封装,包含字段映射
func Copy(toValue, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    false,
		Converters: []copier.TypeConverter{
			TimeToTimeStampPb(),
			TimeToString(),
			GormDeletedAtToTimeStampPb(),
			GormDeletedAtToString(),
			SQLNullTimeToTimeStampPb(),
			SQLNullTimeToString(),
		},
	})
}

// TimeToTimeStampPb time.Time to *timestamppb.Timestamp
func TimeToTimeStampPb() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: &timestamppb.Timestamp{},
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(time.Time)
			return timestamppb.New(t), nil
		},
	}
}

// TimeToString time.Time to string
func TimeToString() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(time.Time)
			return t.Format("2006-01-02 15:04:05"), nil
		},
	}
}

// GormDeletedAtToTimeStampPb gorm.DeletedAt to *timestamppb.Timestamp
func GormDeletedAtToTimeStampPb() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: gorm.DeletedAt{},
		DstType: &timestamppb.Timestamp{},
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(gorm.DeletedAt)
			return timestamppb.New(t.Time), nil
		},
	}
}

// GormDeletedAtToString gorm.DeletedAt to *timestamppb.Timestamp
func GormDeletedAtToString() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: gorm.DeletedAt{},
		DstType: &timestamppb.Timestamp{},
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(gorm.DeletedAt)
			return t.Time.Format("2006-01-02 15:04:05"), nil
		},
	}
}

// SQLNullTimeToTimeStampPb sql.NullTime to *timestamppb.Timestamp
func SQLNullTimeToTimeStampPb() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: sql.NullTime{},
		DstType: &timestamppb.Timestamp{},
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(sql.NullTime)
			return timestamppb.New(t.Time), nil
		},
	}
}

// SQLNullTimeToString sql.NullTime to *timestamppb.Timestamp
func SQLNullTimeToString() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: sql.NullTime{},
		DstType: &timestamppb.Timestamp{},
		Fn: func(src interface{}) (interface{}, error) {
			t := src.(sql.NullTime)
			return t.Time.Format("2006-01-02 15:04:05"), nil
		},
	}
}

// NilSliceToEmptySlice 递归地将 nil 切片设置为空切片
func NilSliceToEmptySlice(inter any) any {
	// original input that can't be modified
	val := reflect.ValueOf(inter)
	switch val.Kind() {
	case reflect.Slice:
		newSlice := reflect.MakeSlice(val.Type(), 0, val.Len())
		if !val.IsZero() {
			// iterate over each element in slice
			for j := 0; j < val.Len(); j++ {
				item := val.Index(j)
				var newItem reflect.Value
				switch item.Kind() {
				case reflect.Struct:
					// recursively handle nested struct
					newItem = reflect.Indirect(reflect.ValueOf(NilSliceToEmptySlice(item.Interface())))
				default:
					newItem = item
				}
				newSlice = reflect.Append(newSlice, newItem)
			}
		}
		return newSlice.Interface()
	case reflect.Struct:
		// new struct that will be returned
		newStruct := reflect.New(reflect.TypeOf(inter))
		newVal := newStruct.Elem()
		// iterate over input's fields
		for i := 0; i < val.NumField(); i++ {
			newValField := newVal.Field(i)
			valField := val.Field(i)
			switch valField.Kind() {
			case reflect.Slice:
				// recursively handle nested slice
				newValField.Set(reflect.Indirect(reflect.ValueOf(NilSliceToEmptySlice(valField.Interface()))))
			case reflect.Struct:
				// recursively handle nested struct
				newValField.Set(reflect.Indirect(reflect.ValueOf(NilSliceToEmptySlice(valField.Interface()))))
			default:
				newValField.Set(valField)
			}
		}

		return newStruct.Interface()
	case reflect.Map:
		// new map to be returned
		newMap := reflect.MakeMap(reflect.TypeOf(inter))
		// iterate over every key value pair in input map
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			// recursively handle nested value
			newV := reflect.Indirect(reflect.ValueOf(NilSliceToEmptySlice(v.Interface())))
			newMap.SetMapIndex(k, newV)
		}
		return newMap.Interface()
	case reflect.Ptr:
		// dereference pointer
		return NilSliceToEmptySlice(val.Elem().Interface())
	default:
		return inter
	}
}
