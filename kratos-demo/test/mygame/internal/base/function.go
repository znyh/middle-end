package base

import (
    "encoding/json"
    "fmt"
    "math"
    "math/rand"
    "net"
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/go-kratos/kratos/v2/log"
)

func init() {
    rand.Seed(time.Now().Unix())
}

//GetTick 毫秒级时间戳
func GetTick() int64 {
    return time.Now().UnixNano() / 1e6
}

//获取当前时间
func GetCurSec() int64 {
    return time.Now().Unix()
}

//Rand [0, max)
func Rand(max int) int {
    if max <= 0 {
        return 0
    }
    return rand.Intn(max)
}

//RandRange [min, max)
func RandRange(min int, max int) int {
    if min < 0 || max <= 0 {
        return 0
    }
    if min >= max {
        return min
    }
    return min + rand.Intn(max-min)
}

func IsHit(max int) bool {
    return Rand(100) < max
}

//ToJSON json string
func ToJSON(v interface{}) string {
    j, err := json.Marshal(v)
    if err != nil {
        return err.Error()
    }
    return string(j)
}

//string转到int
func StrToInt(src string) int {
    if len(src) <= 0 {
        return 0
    }

    //log.Printf("StrToInt:%s.", src)
    dst, err := strconv.Atoi(src)
    if err != nil {
        log.Warn("str to int error(%v).", err)
        return 0
    }
    return dst
}

//string转到int64
func StrToInt64(src string) int64 {
    //	log.Print("src:%s.", src)
    dst, err := strconv.ParseInt(src, 10, 64)
    if err != nil {
        log.Warn("str to int64 error(%v).", err)
        return 0
    }
    return dst
}

//int转到string
func IntToStr(src int) string {
    return strconv.Itoa(src)
}

//int64转到string
func Int64ToStr(src int64) string {
    return strconv.FormatInt(src, 10)
}

//string转到int32
func StrToInt32(src string) int32 {
    dst, err := strconv.Atoi(src)
    if err != nil {
        log.Warn("str to int32 error(%v).", err)
        return 0
    }
    return int32(dst)
}

//int转到string
func Int32ToStr(src int32) string {
    return strconv.Itoa(int(src))
}

//ToString .
func ToString(v interface{}) string {
    return fmt.Sprintf("%v", v)
}

func SliceShuffle(slice []int32) []int32 {
    for i := len(slice) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        slice[i], slice[j] = slice[j], slice[i]
    }
    return slice
}

func SliceContain(slice []int32, values ...int32) bool {
    if len(slice) == 0 || len(values) == 0 {
        return false
    }

    for _, val := range values {
        if !sliceContain(slice, val) {
            return false
        }
    }

    return true
}

func SliceContainAll(slice []int32, values ...int32) bool {
    if len(slice) == 0 || len(values) == 0 {
        return false
    }

    cp := SliceCopy(slice)
    for _, v := range values {
        flag := false
        for i := 0; i < len(cp); i++ {
            if v == cp[i] {
                cp[i] = -1
                flag = true
                break
            }
        }
        if flag == false {
            return false
        }
    }

    return true
}

func sliceContain(slice []int32, value int32) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}

func SliceIndex(slice []int32, value int32) int {
    for k, v := range slice {
        if v == value {
            return k
        }
    }
    return -1
}

func SliceCopy(slice []int32) []int32 {
    dst := make([]int32, len(slice))
    copy(dst, slice)
    return dst
}

func SliceCopyAndDel(slice []int32, del ...int32) []int32 {
    dst := make([]int32, len(slice))
    copy(dst, slice)
    return SliceDel(dst, del...)
}

func SliceDel(slice []int32, values ...int32) []int32 {
    if len(slice) == 0 || len(values) == 0 {
        return slice
    }

    for _, val := range values {
        slice = sliceDel(slice, val)
    }
    return slice
}

func sliceDel(slice []int32, val int32) []int32 {
    for k, v := range slice {
        if v == val {
            return append(append([]int32{}, slice[:k]...), slice[k+1:]...)
        }
    }
    return slice
}

func SliceUnique(slice []int32) []int32 {
    var (
        unique = []int32(nil)
        m      = make(map[int32]bool)
    )
    for _, v := range slice {
        if _, exist := m[v]; !exist {
            m[v] = true
            unique = append(unique, v)
        }
    }
    return unique
}

//升序
func SliceSortRise(slice []int32) []int32 {
    sort.SliceStable(slice, func(i, j int) bool {
        return slice[i] < slice[j]
    })
    return slice
}

//降序
func SliceSortDesc(slice []int32) []int32 {
    sort.SliceStable(slice, func(i, j int) bool {
        return slice[i] > slice[j]
    })
    return slice
}

//两个切片排序后的元素是否全都一样
func SliceSortEqual(slice1, slice2 []int32) bool {
    if len(slice1) != len(slice2) {
        return false
    }
    s1 := SliceSortDesc(SliceCopy(slice1))
    s2 := SliceSortDesc(SliceCopy(slice2))
    for k, v := range s1 {
        if s2[k] != v {
            return false
        }
    }
    return true
}

func SliceReverse(slice []int32) []int32 {
    size := len(slice)
    for i := 0; i < size/2; i++ {
        slice[i], slice[size-i-1] = slice[size-i-1], slice[i]
    }
    return slice
}

//取slice[start,start+n]上的元素
func SliceNSequence(slice []int32, start, n int) []int32 {
    if start+n >= len(slice) {
        return make([]int32, n)
    }
    return slice[start : start+n]
}

func SliceValueCnt(slice []int32, val int32) int {
    var cnt = 0
    for _, v := range slice {
        if v == val {
            cnt++
        }
    }
    return cnt
}

func MInInt(s ...int) int {
    if len(s) <= 0 {
        return -1
    }
    var min = s[0]
    for i := 0; i < len(s); i++ {
        if s[i] < min {
            min = s[i]
        }
    }
    return min
}

func MInInt32(s ...int32) int32 {
    if len(s) <= 0 {
        return 0
    }
    var min = s[0]
    for i := 0; i < len(s); i++ {
        if s[i] < min {
            min = s[i]
        }
    }
    return min
}

func MinInt64(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func MinFloat64(s ...float64) float64 {
    if len(s) <= 0 {
        return 0
    }
    var min = s[0]
    for i := 0; i < len(s); i++ {
        if s[i] < min {
            min = s[i]
        }
    }
    return min
}

func MaxInt32(s ...int32) int32 {
    if len(s) <= 0 {
        return 0
    }
    max := s[0]
    for i := 0; i < len(s); i++ {
        if s[i] > max {
            max = s[i]
        }
    }
    return max
}

func MaxInt64(s ...int64) int64 {
    if len(s) <= 0 {
        return 0
    }
    max := s[0]
    for i := 0; i < len(s); i++ {
        if s[i] > max {
            max = s[i]
        }
    }
    return max
}

func GetGameNo(roomId, tableId int32) int64 {
    t := GetTick() / 1000
    return t<<16 + int64(roomId)<<8 + int64(tableId)
}

//Copy deep copy struct 注意：只对结构体大写变量生效!
func Copy(src, des interface{}) error {
    j, err := json.Marshal(src)
    if err != nil {
        return err
    }

    return json.Unmarshal(j, des)
}

//数组值都是唯一的
func IsUnique(a []int32) bool {
    m := map[int32]bool{}
    for _, v := range a {
        if m[v] {
            return false
        }

        m[v] = true
    }

    return true
}

//指定概率的随机事件 返回落到probability对应的概率的下标
func RandEvent(probability []int32) int {
    r := rand.Int31n(99)
    start, end := int32(0), int32(0)

    for k, v := range probability {
        end = start + v - 1
        if r >= start && r <= end {
            return k

        }
        start += v
    }
    return 0
}

//字符串表示的16进制转化为整形
func HexStr2Int32(src string) (dst []int32) {
    if len(src) <= 0 {
        return
    }

    tmp := strings.TrimSpace(src)
    s := strings.Split(tmp, ",")

    for _, v := range s {
        vv := strings.TrimSpace(v)
        if val, err := strconv.ParseInt(vv, 0, 0); err == nil {
            dst = append(dst, int32(val))
        } else {
            return []int32{}
        }
    }

    return dst
}

//在a中取m个元素的所有组合
func Dfs(a []int32, k int) [][]int32 {
    if len(a) < k {
        return nil
    }
    var (
        n   = len(a)
        vis = []int32(nil)
        re  = []int32(nil)
        dst = [][]int32(nil)
    )

    for i := 0; i < n; i++ {
        vis = append(vis, 0)
        re = append(re, 0)
    }

    dfs(a, vis, re, &dst, n, k, 0, 0)

    return dst
}

//参数step代表选取第几个数字，参数start代表从集合的第几个开始选
func dfs(a, vis, re []int32, dst *[][]int32, n, k, start, step int) {
    //如果选够了k个就输出
    if step == k {
        var c = make([]int32, len(re))
        copy(c, re)
        *dst = append(*dst, c[:k])
    }

    for i := start; i < n; i++ {
        if vis[i] == 1 {
            continue
        }

        vis[i] = 1
        re[step] = a[i]

        dfs(a, vis, re, dst, n, k, i+1, step+1)
        vis[i] = 0
    }

    return
}

//全排列
func Permute(nums []int32) [][]int32 {
    var res [][]int32
    if len(nums) == 0 {
        return res
    }
    var tmp []int32
    var visited = make([]bool, len(nums))
    backtracking(nums, &res, tmp, visited)
    return res
}

func backtracking(nums []int32, res *[][]int32, tmp []int32, visited []bool) {
    // 成功找到一组
    if len(tmp) == len(nums) {
        var c = make([]int32, len(tmp))
        copy(c, tmp)
        *res = append(*res, c)
        return
    }
    // 回溯
    for i := 0; i < len(nums); i++ {
        if !visited[i] {
            visited[i] = true
            tmp = append(tmp, nums[i])
            backtracking(nums, res, tmp, visited)
            tmp = tmp[:len(tmp)-1]
            visited[i] = false
        }
    }
}

//组合总和
func combinationSum2(candidates []int, target int) [][]int {
    var res [][]int
    if len(candidates) == 0 {
        return res
    }
    sort.Ints(candidates)
    var tmp []int
    var visited = make([]bool, len(candidates))
    //backtracking2(candidates, &res, target, tmp, 0, visited)
    backtracking2(candidates, &res, target, tmp, 1, visited)
    return res
}

func backtracking2(candidates []int, res *[][]int, target int, tmp []int, index int, visited []bool) {
    if target < 0 {
        return
    }
    if target == 0 {
        var c = make([]int, len(tmp))
        copy(c, tmp)
        *res = append(*res, c)
        return
    }
    for i := index; i < len(candidates); i++ {
        if i > 0 && candidates[i] == candidates[i-1] && !visited[i-1] {
            continue
        }
        visited[i] = true
        //tmp = append(tmp, candidates[i])  //记录元素值
        tmp = append(tmp, i) //记录元素的下标
        backtracking2(candidates, res, target-candidates[i], tmp, i+1, visited)
        tmp = tmp[:len(tmp)-1]
        visited[i] = false
    }
}

func GetIP() string {

    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "127.0.0.1"
    }

    for _, address := range addrs {
        // 检查ip地址判断是否回环地址
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }

    return "127.0.0.1"
}

func Abs(x int32) int32 {
    return int32(math.Abs(float64(x)))
}

func AExceptB(a, b []int32) []int32 {
    ret := make([]int32, 0)
    for _, v := range a {
        if !isIn(v, b) {
            ret = append(ret, v)
        }
    }

    return ret
}

func BAllInA(a, b []int32) bool {
    //for _, v := range b {
    //	if !isIn(v, a) {
    //		return false
    //	}
    //}

    m1, m2 := map[int32]int{}, map[int32]int{}
    for _, v := range a {
        m1[v]++
    }
    for _, v := range b {
        m2[v]++
    }
    for k, v := range m2 {
        if m1[k] < v {
            return false
        }
    }
    return true
}

func isIn(e int32, ems []int32) bool {
    for _, v := range ems {
        if e == v {
            return true
        }
    }
    return false
}

func GetTimeByStr(d string, t string) (time.Time, bool) {

    if d == "" || t == "" {
        return time.Now(), false
    }

    st, err := time.ParseInLocation("2006-01-02 15:04:05", d+" "+t, time.Local)

    if err != nil {
        log.Info("GetTimeByStr Fail.%+v", err)
        return time.Now(), false
    }

    return st, true
}

func IsBetweenTime(st, et, ct time.Time) bool {

    //比较日期
    if ct.Unix() < st.Unix() || ct.Unix() > et.Unix() {
        return false
    }

    //比较时间段
    stD := st.Hour()*3600 + st.Minute()*60 + st.Second()
    etD := et.Hour()*3600 + et.Minute()*60 + et.Second()
    ctD := ct.Hour()*3600 + ct.Minute()*60 + ct.Second()
    if ctD < stD || ctD > etD {
        return false
    }

    return true
}

// 随机生成IP
func RandIP() string {

    // 中国国内网段
    ipList := [10][2]int32{
        {607649792, 608174079},     //36.56.0.0-36.63.255.255
        {1038614528, 1039007743},   //61.232.0.0-61.237.255.255
        {1783627776, 1784676351},   //106.80.0.0-106.95.255.255
        {2035023872, 2035154943},   //121.76.0.0-121.77.255.255
        {2078801920, 2079064063},   //123.232.0.0-123.235.255.255
        {-1950089216, -1948778497}, //139.196.0.0-139.215.255.255
        {-1425539072, -1425014785}, //171.8.0.0-171.15.255.255
        {-1236271104, -1235419137}, //182.80.0.0-182.92.255.255
        {-770113536, -768606209},   //210.25.0.0-210.47.255.255
        {-569376768, -564133889},   //222.16.0.0-222.95.255.255
    }

    idx := Rand(10)
    ipV := Rand(int(ipList[idx][1] - ipList[idx][0]))
    ip := ipList[idx][0] + int32(ipV)

    a1 := (int)((ip >> 24) & 0xff)
    a2 := (int)((ip >> 16) & 0xff)
    a3 := (int)((ip >> 8) & 0xff)
    a4 := (int)(ip & 0xff)

    return fmt.Sprintf("%d.%d.%d.%d", a1, a2, a3, a4)
}
