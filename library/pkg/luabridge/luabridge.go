package luabridge

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/AzureWrathCyd/gluamapper"
	lua "github.com/AzureWrathCyd/gopher-lua"
)

var typeMap = map[reflect.Kind]lua.LValueType{
	reflect.String:  lua.LTString,
	reflect.Int:     lua.LTNumber,
	reflect.Int8:    lua.LTNumber,
	reflect.Int16:   lua.LTNumber,
	reflect.Int32:   lua.LTNumber,
	reflect.Int64:   lua.LTNumber,
	reflect.Uint:    lua.LTNumber,
	reflect.Uint8:   lua.LTNumber,
	reflect.Uint16:  lua.LTNumber,
	reflect.Uint32:  lua.LTNumber,
	reflect.Uint64:  lua.LTNumber,
	reflect.Float32: lua.LTNumber,
	reflect.Float64: lua.LTNumber,
	reflect.Bool:    lua.LTBool,
	reflect.Struct:  lua.LTTable,
	//暂时只支持传结构体引用
	reflect.Ptr: lua.LTTable,
}

func toLuaNumber(goValue reflect.Value) (ret lua.LNumber) {
	ret = 0
	switch goValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret = lua.LNumber(goValue.Int())
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret = lua.LNumber(goValue.Uint())
		break
	case reflect.Float32, reflect.Float64:
		ret = lua.LNumber(goValue.Float())
		break
	}
	return ret
}

func toLuaFunc(L *lua.LState, value reflect.Value, self *reflect.Value) *lua.LFunction {
	luaFunc := func(L *lua.LState) int {
		if paramList, err := getParamListFromLua(L, value.Type(), self != nil); err != nil {
			fmt.Println("fuck param list error,err:" + err.Error())
			return 0
		} else {
			if self != nil {
				paramList = append([]reflect.Value{*self}, paramList...)
			}
			retList := value.Call(paramList)
			return sendRetListToLua(L, value.Type(), retList)
		}
	}
	return L.NewFunction(luaFunc)
}

func toLuaTable(L *lua.LState, goValue reflect.Value) *lua.LTable {
	t := &lua.LTable{}
	rType := goValue.Type()
	for i := 0; i < rType.NumMethod(); i++ {
		method := rType.Method(i)
		name := method.Name
		luaFunc := toLuaFunc(L, method.Func, &goValue)
		t.RawSet(lua.LString(name), luaFunc)
	}
	return t
}

func toLuaValue(L *lua.LState, goValue reflect.Value) (ret lua.LValue, err error) {
	ret = lua.LNil
	rKind := goValue.Kind()
	if luaType, ok := typeMap[rKind]; ok == false {
		err = errors.New("unsupport type " + goValue.Type().String())
	} else {
		switch luaType {
		case lua.LTString:
			tmpStr := goValue.String()
			ret = lua.LString(tmpStr)
			break
		case lua.LTBool:
			tmpBool := goValue.Bool()
			ret = lua.LBool(tmpBool)
			break
		case lua.LTNumber:
			ret = toLuaNumber(goValue)
			break
		case lua.LTTable:
			ret = toLuaTable(L, goValue)
			break
		}
	}

	return ret, err
}

func toGoNumber(retRef interface{}, luaNumber lua.LNumber) {
	rKind := reflect.TypeOf(retRef).Elem().Kind()
	switch rKind {
	case reflect.Int:
		*(retRef.(*int)) = int(luaNumber)
		break
	case reflect.Int8:
		*(retRef.(*int8)) = int8(luaNumber)
		break
	case reflect.Int16:
		*(retRef.(*int16)) = int16(luaNumber)
		break
	case reflect.Int32:
		*(retRef.(*int32)) = int32(luaNumber)
		break
	case reflect.Int64:
		*(retRef.(*int64)) = int64(luaNumber)
		break
	case reflect.Uint:
		*(retRef.(*uint)) = uint(luaNumber)
		break
	case reflect.Uint8:
		*(retRef.(*uint8)) = uint8(luaNumber)
		break
	case reflect.Uint16:
		*(retRef.(*uint16)) = uint16(luaNumber)
		break
	case reflect.Uint32:
		*(retRef.(*uint32)) = uint32(luaNumber)
		break
	case reflect.Uint64:
		*(retRef.(*uint64)) = uint64(luaNumber)
		break
	case reflect.Float32:
		*(retRef.(*float32)) = float32(luaNumber)
		break
	case reflect.Float64:
		*(retRef.(*float64)) = float64(luaNumber)
		break
	}
}

func toGoValue(goValue interface{}, luaValue *lua.LValue) (err error) {
	rKind := reflect.TypeOf(goValue).Elem().Kind()
	if luaType, ok := typeMap[rKind]; ok == false {
		err = errors.New("unsupport type " + reflect.TypeOf(goValue).Elem().String())
	} else if luaType != (*luaValue).Type() {
		err = errors.New("not required type, need type:" + reflect.TypeOf(goValue).Elem().String() + ",ret type:" + (*luaValue).Type().String())
	} else {
		switch luaType {
		case lua.LTString:
			*(goValue.(*string)) = lua.LVAsString(*luaValue)
			break
		case lua.LTBool:
			*(goValue.(*bool)) = lua.LVAsBool(*luaValue)
			break
		case lua.LTNumber:
			toGoNumber(goValue, lua.LVAsNumber(*luaValue))
			break
		case lua.LTTable:
			err = gluamapper.Map((*luaValue).(*lua.LTable), goValue)
			break
		}
	}

	return err
}

func toReflectValue(t reflect.Type, luaValue lua.LValue) (ret reflect.Value, err error) {
	rKind := t.Kind()
	if luaType, ok := typeMap[rKind]; ok == false {
		err = errors.New("unsupport type " + t.String())
	} else if luaType != (luaValue).Type() {
		err = errors.New("not required type, need type:" + t.String() + ",param type:" + luaValue.Type().String())
	} else {
		switch luaType {
		case lua.LTString:
			ret = reflect.ValueOf(lua.LVAsString(luaValue))
			break
		case lua.LTBool:
			ret = reflect.ValueOf(lua.LVAsBool(luaValue))
			break
		case lua.LTNumber:
			ret = reflect.ValueOf(lua.LVAsNumber(luaValue))
			break
		case lua.LTTable:
			//New return a pointer to a new zero value
			tmpValue := reflect.New(t).Interface()
			if err = gluamapper.Map(luaValue.(*lua.LTable), tmpValue); err == nil {
				ret = reflect.ValueOf(tmpValue).Elem()
			}
			break
		}
	}
	return ret, err
}

//SafeCall call lua global function with protect,
//retRef must be a pointer
func SafeCall(L *lua.LState, funcName string, retRef interface{}, args ...interface{}) (err error) {
	var rKind reflect.Kind
	if retRef == nil {
		rKind = reflect.Invalid
	} else {
		rKind = reflect.TypeOf(retRef).Kind()
	}

	if rKind != reflect.Ptr && rKind != reflect.Invalid {
		err = errors.New("param retRef must be a pointer")
		return
	}

	var luaArgs []lua.LValue
	var luaArg lua.LValue

	for _, arg := range args {
		if luaArg, err = toLuaValue(L, reflect.ValueOf(arg)); err != nil {
			return
		}
		luaArgs = append(luaArgs, luaArg)
	}
	//目前只支持无返回和返回一个值
	var nRet int
	if rKind == reflect.Invalid {
		nRet = 0
	} else {
		nRet = 1
	}
	if err = L.CallByParam(
		lua.P{
			Fn:      L.GetGlobal(funcName),
			NRet:    nRet,
			Protect: true,
		}, luaArgs...,
	); err != nil {
		return
	}
	if nRet == 1 {
		tmpRet := L.Get(-1)
		if tmpRet != lua.LNil {
			L.Pop(1)
		}
		if err = toGoValue(retRef, &tmpRet); err != nil {
			return
		}
	}

	return err
}

//SafeLoad load lua file by name with protect
func SafeLoad(fileName string) (L *lua.LState, err error) {
	L = lua.NewState()
	err = L.DoFile(fileName)
	if err != nil {
		L = nil
	}
	return L, err
}

//SafeLoad load lua str with protect
func SafeLoadString(str string) (L *lua.LState, err error) {
	L = lua.NewState()
	err = L.DoString(str)
	if err != nil {
		L = nil
	}
	return L, err
}

func getParamListFromLua(L *lua.LState, funcType reflect.Type, hasSelf bool) (paramList []reflect.Value, err error) {
	numIn := funcType.NumIn()
	var i = 0
	if hasSelf {
		i++
	}
	for ; i < numIn; i++ {
		luaParam := L.Get(-numIn + i)
		if param, tmpErr := toReflectValue(funcType.In(i), luaParam); tmpErr != nil {
			return nil, tmpErr
		} else {
			paramList = append(paramList, param)
		}
	}
	return paramList, nil
}

func sendRetListToLua(L *lua.LState, funcType reflect.Type, retList []reflect.Value) int {
	pushNum := 0
	for _, ret := range retList {
		if luaRet, err := toLuaValue(L, ret); err == nil {
			L.Push(luaRet)
			pushNum++
		} else {
			L.Pop(pushNum)
			return 0
		}
	}
	return pushNum
}

func RegGlobalFunction(L *lua.LState, f interface{}, name string) (err error) {
	funcValue := reflect.ValueOf(f)
	if funcValue.Kind() != reflect.Func {
		err = errors.New("f not a function")
		return
	}

	L.SetGlobal(name, toLuaFunc(L, funcValue, nil))
	return nil
}
