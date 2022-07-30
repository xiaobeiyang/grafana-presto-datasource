package main

import (
	"database/sql"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

const (
	dateFormat      = "2006-01-02"
	dateTimeFormat1 = "2006-01-02 15:04:05"
	dateTimeFormat2 = "2006-01-02T15:04:05Z"
	dataTimeFormat3 = time.RFC3339
	dataTimeFormat4 = time.RFC3339Nano

	timeFormatNano  = "15:04:05.000000000"
	timeFormatMicro = "15:04:05.000000"
	timeFormatMilli = "15:04:05.000"
	timeFormat      = "15:04:05"
)

func Converters() []sqlutil.Converter {
	return []sqlutil.Converter{
		{

			Name:          "handle boolean",
			InputTypeName: "boolean",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableBool,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					boolean := ns.String == "true"
					return &boolean, nil
				},
			},
		},
		{
			Name:          "handle tinyint",
			InputTypeName: "tinyint",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableInt8,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseInt(ns.String, 10, 8)
					if err != nil {
						return nil, err
					}
					v8 := int8(v)
					return &v8, nil
				},
			},
		},
		{
			Name:          "handle smallint",
			InputTypeName: "smallint",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableInt16,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseInt(ns.String, 10, 16)
					if err != nil {
						return nil, err
					}
					v16 := int16(v)
					return &v16, nil
				},
			},
		},
		{
			Name:          "handle int/integer",
			InputTypeName: "integer",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableInt32,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseInt(ns.String, 10, 32)
					if err != nil {
						return nil, err
					}
					v32 := int32(v)
					return &v32, nil
				},
			},
		},
		{
			Name:          "handle bigint",
			InputTypeName: "bigint",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableInt64,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseInt(ns.String, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:          "handle real",
			InputTypeName: "real",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableFloat32,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseFloat(ns.String, 32)
					if err != nil {
						return nil, err
					}
					v32 := float32(v)
					return &v32, nil
				},
			},
		},
		{
			Name:          "handle double",
			InputTypeName: "double",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableFloat64,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseFloat(ns.String, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle decimal",
			InputTypeRegex: regexp.MustCompile("decimal.*"),
			InputScanType:  reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableFloat64,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := strconv.ParseFloat(ns.String, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:          "handle date",
			InputTypeName: "date",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableTime,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := time.Parse(dateFormat, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat1, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat2, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dataTimeFormat3, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dataTimeFormat4, ns.String)
					if err == nil {
						return &v, nil
					}

					return nil, err
				},
			},
		},
		{
			Name:          "handle time",
			InputTypeName: "time",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableTime,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := time.Parse(timeFormatNano, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(timeFormatMicro, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(timeFormatMilli, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(timeFormat, ns.String)
					if err == nil {
						return &v, nil
					}

					return nil, err
				},
			},
		},
		{
			Name:          "handle timestamp",
			InputTypeName: "timestamp",
			InputScanType: reflect.TypeOf(sql.NullString{}),
			FrameConverter: sqlutil.FrameConverter{
				FieldType: data.FieldTypeNullableTime,
				ConverterFunc: func(in interface{}) (interface{}, error) {
					ns := in.(*sql.NullString)
					if !ns.Valid {
						return nil, nil
					}
					v, err := time.Parse(dateTimeFormat1, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat2, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dataTimeFormat3, ns.String)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dataTimeFormat4, ns.String)
					if err == nil {
						return &v, nil
					}

					return nil, err
				},
			},
		},
	}
}
