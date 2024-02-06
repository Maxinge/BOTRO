package main

import(
	"os"
	// "os/exec"
	"path/filepath"
	// "path"
	"io/ioutil"
	// "strings"
	"fmt"
	// "io"
	"encoding/json"
	"runtime"
	// "sort"
	"reflect"
	"strconv"
	// "net"
	// "time"
	// "golang.org/x/net/html/charset"
	// "golang.org/x/text/transform"
)


func Slh() string{
	if runtime.GOOS == "windows" { return "\\" }
	return "/"
}
func CurDir() string {
	ex, err := os.Executable()
	if err != nil { panic(err) }
	return filepath.Dir(ex)+Slh()
}


func Itos(i int) string { return strconv.Itoa(i) }
func Stoi(s string) int { ss,_:=strconv.Atoi(s); return ss }

// func F64toS(f float64) string { return fmt.Sprintf("%f", f) }
// func StoF64(s string) float64 { ss,_:=strconv.ParseFloat(s, 64); return ss }

// func RemoveStrSlice(slice []string, s int) []string {
//     return append(slice[:s], slice[s+1:]...)
// }

// func power(base, exponent int) int {
// 	result := 1
// 	for i := 0; i < exponent; i++ { result *= base }
// 	return result
// }

func byteArrayToUInt64(b []byte) uint64 {
	if len(b) < 4 { return 0 }
	result := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
              uint64(b[3])<<32 | uint64(b[3])<<40 | uint64(b[3])<<48 | uint64(b[3])<<56
	return result
}

func byteArrayToUInt32(b []byte) uint32 {
	if len(b) < 4 { return 0 }
	result := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return result
}

func byteArrayToUInt16(b []byte) uint16 {
	if len(b) < 2 { return 0 }
	result := uint16(b[0]) | uint16(b[1])<<8
	return result
}

func inArrayByte(bb [][]byte,search []byte) bool{
    for _,v := range bb {
        if sliceEqual(v,search) { return true }
    }
    return false
}

func sliceEqual(slice1 []byte, slice2 []byte) bool {
	if len(slice1) != len(slice2) { return false }
	for i := range slice1 {
		if slice1[i] != slice2[i] { return false }
	}
	return true
}

func splitBitsArrayMulti(bb []byte ,search [][]byte) [][]byte {
    slen := len(search[0])
	for _,vv := range search {
		if len(vv) != slen{
			if err != nil { fmt.Printf("err splitBitsArray : all search need same len \n", err) }
			return [][]byte{}
		}
	}
    var res [][]byte
    j := 0
    for i := 0; i <= len(bb)-slen ; i++ {
        if inArrayByte(search,bb[i:i+slen]) {
            res = append(res,bb[j:i])
            i += slen
            j = i
        }
    }
    res = append(res,bb[j:len(bb)])
    return res
}

func splitBitsArray(bb []byte ,search []byte) [][]byte {
    slen := len(search)
    var res [][]byte
    j := 0
    for i := 0; i <= len(bb)-slen ; i++ {
        if sliceEqual(bb[i:i+slen], search) {
            res = append(res,bb[j:i])
            i += slen
            j = i
        }
    }
    res = append(res,bb[j:len(bb)])
    return res
}

func dumpArray(arr []byte) string{
    ss := "{"
	for i, b := range arr {
		ss += fmt.Sprintf("%d", b)
		if i < len(arr)-1 { ss += "," }
	}
    return ss + "}"
}

func dumpArrayTab(arr []byte) string{
    ss := "{"
	for i, b := range arr {
		ss += fmt.Sprintf("%d", b)
		if i < len(arr)-1 { ss += ",\t" }
	}
    return ss + "}"
}

func printStruct(s interface{}) string{
    str := ""
    v := reflect.ValueOf(s)
    for i := 0; i < v.NumField(); i++ {
        fieldName := v.Type().Field(i).Name
        fieldValue := v.Field(i).Interface()
        str += fmt.Sprintf("%s", fieldName) + " = "
        str += fmt.Sprintf("%v", fieldValue) + "\n"
    }
    return str
}

// func SortedKeysStr(m map[string]string) ([]string) {
//     keys := make([]string, len(m))
//     i := 0
//     for k := range m {
//         keys[i] = k
//         i++
//     }
//     sort.Strings(keys)
//     return keys
// }
//
// func ExecCMD(cmd string,args... string) {
// 	err := exec.Command(cmd,args...).Run()
// 	if err != nil { fmt.Printf("err ExecCMD -- %v -- \n", err) }
// }
//
// func GetOutboundIP() string {
//     conn, err := net.Dial("udp", "8.8.8.8:80")
//     if err != nil { fmt.Printf("err -- %v -- \n", err) }
//     defer conn.Close()
//     return conn.LocalAddr().String()
// }
//
// func SortedKeysInt(m map[int]interface{}) ([]int) {
//     keys := make([]int, len(m))
//     i := 0
//     for k := range m {
//         keys[i] = k
//         i++
//     }
//     sort.Ints(keys)
//     return keys
// }
//
// func IndexOf(element string, data []string) (int) {
//    for k, v := range data {
//        if element == v {
//            return k
//        }
//    }
//    return -1    //not found.
// }
//
// func SortedKeysIntTab(m map[int][]string) ([]int) {
//     keys := make([]int, len(m))
//     i := 0
//     for k := range m {
//         keys[i] = k ; i++
//     }
//     sort.Ints(keys)
//     return keys
// }
//
// func StringInSlice(a string, list []string) bool {
//     for _, b := range list {
//         if b == a {
//             return true
//         }
//     }
//     return false
// }
//
// func ClearTerminal(){
// 	cmd := exec.Command("clear")
// 	if runtime.GOOS == "windows" {
// 		cmd = exec.Command("cmd", "/c", "cls")
// 	}
// 	cmd.Stdout = os.Stdout
// 	cmd.Run()
// }

func StrInArray(str string,arr []string) bool{
    for _, v := range arr { if str == v { return true } }
    return false
}
//
//
func intInArray(i int,ii []int) bool{
    for _, v := range ii {
        if i == v { return true }
    }
    return false
}

func GetBetween(s string,s1 string,s2 string) []string {
	s1len := len(s1)
	s2len := len(s2)
	var res []string
    for i := 0; i <= len(s)-s1len ; i++ {
		if s[i:i+s1len] == s1 {
			j := i+s1len
			for {
				j++
				if j > len(s)-s2len { break }
				if s[j:j+s2len] == s2 {
					res = append(res,s[i+s1len:j])
					i = j + s2len
					break
				}
			}
		}
    }
	return res
}
//
//
// func ExtLower(path string) string{
// 	return strings.ToLower(path[strings.LastIndex(string(path), ".")+1:len(path)])
// }
//
// func GetType(f interface{}) reflect.Type{
// 	return reflect.ValueOf(f).Type()
// }
//
// func DecodeJSON(s string) interface{}{
// 	var f interface{}
//     err := json.Unmarshal([]byte(s), &f)
// 	if err != nil { fmt.Printf("err -- %v -- \n", err) }
// 	return f
// }
//
//
// func EncodeJSON(f interface{}) string{
// 	b, err := json.Marshal(f)
// 	if err != nil { fmt.Printf("err -- %v -- \n", err) }
// 	return string(b)
// }
//
// func Exists(name string) bool {
//     if _, err := os.Stat(name); err != nil {
//     	if os.IsNotExist(err) { return false }
//     }
//     return true
// }
//
func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}
//
// func CopyFile(src, dst string) (err error) {
// 	f_src, _ := os.Open(src)
// 	os.MkdirAll(path.Dir(dst),0777)
// 	f_dest, err := os.Create(dst)
// 	if err != nil {
// 		fmt.Printf("cannot save -- %v -- \n", dst)
// 		return
// 	}
// 	io.Copy(f_dest, f_src);
// 	f_src.Close()
// 	f_dest.Close()
// 	return
// }
//
func readFileString(file_path string) string {
	data, err := ioutil.ReadFile(file_path)
	if err != nil { fmt.Printf("err readFileString-- %v -- \n", err) }
	return string(data)
}
func writeFileString(file_path string,content string) {
	f, err := os.Create(file_path)
    if err != nil { fmt.Printf("err writeFileString-- %v -- \n", err) }
    f.WriteString(content)
    f.Close();
}
//
// func readRecursDir(dirname string)  {
//     files, _ := ioutil.ReadDir(dirname)
// 	// subdir := strings.Replace(dirname+"/", root_dir, "", 1)
//     // subdir = strings.Replace(subdir, "/", "", 1)
//     for _, f := range files {
//         if !f.IsDir() {
// 			// do
//         }else{
// 			// do
//             readRecursDir(dirname+"/"+f.Name())
//         }
//     }
// }
//
// func DecodeToUTF8(s string, from string) string {
// 	enc, _ := charset.Lookup(from)
// 	r := transform.NewReader(strings.NewReader(s), enc.NewDecoder())
// 	result, _ := ioutil.ReadAll(r)
// 	return string(result)
// }
//
// func diffDate(a, b time.Time) (year, month, day, hour, min, sec int) {
//     if a.Location() != b.Location() {
//         b = b.In(a.Location())
//     }
//     if a.After(b) {
//         a, b = b, a
//     }
//     y1, M1, d1 := a.Date()
//     y2, M2, d2 := b.Date()
//
//     h1, m1, s1 := a.Clock()
//     h2, m2, s2 := b.Clock()
//
//     year = int(y2 - y1)
//     month = int(M2 - M1)
//     day = int(d2 - d1)
//     hour = int(h2 - h1)
//     min = int(m2 - m1)
//     sec = int(s2 - s1)
//
//     // Normalize negative values
//     if sec < 0 {
//         sec += 60
//         min--
//     }
//     if min < 0 {
//         min += 60
//         hour--
//     }
//     if hour < 0 {
//         hour += 24
//         day--
//     }
//     if day < 0 {
//         // days in month:
//         t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
//         day += 32 - t.Day()
//         month--
//     }
//     if month < 0 {
//         month += 12
//         year--
//     }
//
//     return
// }
