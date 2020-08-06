# 小工具
### 安装
```shell
go get github.com/heyhip/frog 
```

### 单个字符串，驼峰转带下划线，常见json格式
```go
s1 := "UserId"
s1Res := frog.Camel2Case(s1)
fmt.Println("s1Res = " + s1Res)
// s1Res = user_id
```

### 单个字符串，带下划线转驼峰
```go
s2 := "user_id"
s2Res := frog.Case2Camel(s2)
fmt.Println("s2Res = " + s2Res)
// s2Res = UserId
```

### struct转map，支持组合
```go
type Stu struct {
    UserId   int64
    UserName string
    Type     int8
}

var stu1 Stu
stu1.UserId = 12345
stu1.UserName = "张三"
stu1.Type = 1

// struct转map，支持嵌套
m1 := frog.StructToMap(stu1)
fmt.Println(m1)
// map[Type:1 UserId:12345 UserName:张三]
```

### map转换为json，再使用json.Unmarshal()绑定时，如果类型不同，字段绑定不了，以下解决这个问题

### map转struct，支持嵌套，中间先转换为json，后可绑定
```go
type StuSub struct {
    Stu
    Sex int8
}
var stu2 StuSub

maps1 := make(map[string]interface{})
maps1["UserId"] = 12345
maps1["UserName"] = "张三"
maps1["Type"] = 1
maps1["Type"] = 1
maps1["Sex"] = 1
b, e, ok := frog.MapToStruct(maps1, stu2)
fmt.Println(ok, e)
// 此处为json
fmt.Println(string(b))
// true <nil>
// {"Sex":1,"Type":1,"UserId":12345,"UserName":"张三"}

// 绑定struct
json.Unmarshal(b, &stu2)
fmt.Println(stu2.UserName)
fmt.Println(stu2.Sex)
// 张三
// 1
```

### 问题：在redis获取全部数据时，返回的是map[string]string，转换为json然后绑定到struct时，无法绑定到struct
redis.HGetAll(key).Val()
报错：cannot unmarshal string into Go struct field userInfo.fans of type int
### map[string]string转struct，下面方式解决这个问题
```go
maps2 := make(map[string]string)
maps2["UserId"] = "12345"
maps2["UserName"] = "张三"
maps2["Type"] = "1"
maps2["Type"] = "1"
maps2["Sex"] = "1"

b, _, _ = frog.MapStringToStruct(maps2, stu2)
fmt.Println(string(b))
// {"Sex":1,"Type":1,"UserId":12345,"UserName":"张三"}
json.Unmarshal(b, &stu2)
fmt.Println(stu2.UserName)
fmt.Println(stu2.Sex)
// 张三
// 1
```

map格式，驼峰转为常见json下划线样式，支持嵌套
```go
maps3 := make(map[string]interface{})
maps3["AccountId"] = "22222"
maps3["PassWord"] = "88888"
maps1["m3"] = maps3
//m1 := frog.MapCase2Camel(maps1)
m1 := frog.MapCamel2Case(maps1)
j1, _ := json.Marshal(m1)
fmt.Println(string(j1))
// {"m3":{"account_id":"22222","pass_word":"88888"},"sex":1,"type":1,"user_id":12345,"user_name":"张三"}
map格式，json带下划线格式，转换为驼峰，支持嵌套
```

```go
m2 := frog.MapCase2Camel(m1)
j2, _ := json.Marshal(m2)
fmt.Println(string(j2))
// {"M3":{"AccountId":"22222","PassWord":"88888"},"Sex":1,"Type":1,"UserId":12345,"UserName":"张三"}
```



