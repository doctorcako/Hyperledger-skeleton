package models

import (
	"reflect"
	"strconv"
)

// DefaultTag - Server model
func (a HttpHandler) DefaultTag() HttpHandler {
	aAux := HttpHandler{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Type"); exist {
		aAux.Type = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Url"); exist {
		aAux.Url = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Port"); exist {
		aAux.Port = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("NetworkType"); exist {
		aAux.NetworkType = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Swagger"); exist {
		aAux.Swagger = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("LogLevel"); exist {
		aAux.LogLevel = f.Tag.Get("default")
	}

	return aAux
}

// DefaultTag - Micro model
func (a Micro) DefaultTag() Micro {

	aAux := Micro{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Name"); exist {
		aAux.Name = f.Tag.Get("default")
	}

	return aAux
}

// DefaultTag - DB config
func (a Db) DefaultTag() Db {
	aAux := Db{}

	aAux.Postgres = aAux.Postgres.DefaultTag()
	aAux.Redis = aAux.Redis.DefaultTag()

	return aAux
}

// DefaultTag - Postgres config
func (p Postgres) DefaultTag() Postgres {
	aAux := Postgres{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Host"); exist {
		aAux.Host = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Port"); exist {
		aAux.Port = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("User"); exist {
		aAux.User = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Password"); exist {
		aAux.Password = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Database"); exist {
		aAux.Database = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Schema"); exist {
		aAux.Schema = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("MaxOpenConns"); exist {
		if value, err := strconv.Atoi(f.Tag.Get("default")); err == nil {
			aAux.MaxOpenConns = value
		}
	}

	if f, exist := typ.FieldByName("MaxIdleConns"); exist {
		if value, err := strconv.Atoi(f.Tag.Get("default")); err == nil {
			aAux.MaxIdleConns = value
		}
	}

	return aAux
}

// DefaultTag - Redis config
func (p Redis) DefaultTag() Redis {
	aAux := Redis{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Host"); exist {
		aAux.Host = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Port"); exist {
		aAux.Port = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Database"); exist {
		if value, err := strconv.Atoi(f.Tag.Get("default")); err == nil {
			aAux.Database = value
		}
	}

	if f, exist := typ.FieldByName("MaxOpenConns"); exist {
		if value, err := strconv.Atoi(f.Tag.Get("default")); err == nil {
			aAux.MaxOpenConns = value
		}
	}

	if f, exist := typ.FieldByName("MaxIdleConns"); exist {
		if value, err := strconv.Atoi(f.Tag.Get("default")); err == nil {
			aAux.MaxIdleConns = value
		}
	}

	return aAux
}
