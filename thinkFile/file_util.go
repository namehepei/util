package thinkFile

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"util/think"
)

//权限说明：
//O_RDONLY int = syscall.O_RDONLY // 只读
//O_WRONLY int = syscall.O_WRONLY // 只写
//O_RDWR int = syscall.O_RDWR // 读写
//O_APPEND int = syscall.O_APPEND // 在文件末尾追加，打开后cursor在文件结尾位置
//O_CREATE int = syscall.O_CREAT // 如果不存在则创建
//O_EXCL int = syscall.O_EXCL //与O_CREATE一起用，构成一个新建文件的功能，它要求文件必须不存在
//O_SYNC int = syscall.O_SYNC // 同步方式打开，没有缓存，这样写入内容直接写入硬盘，系统掉电文件内容有一定保证
//O_TRUNC int = syscall.O_TRUNC // 打开并清空文件

// fileMode说明
// The defined file mode bits are the most significant bits of the FileMode.
// The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// The values of these bits should be considered part of the public API and
// may be used in wire protocols or disk representations: they must not be
// changed, although new bits might be added.
//const (
//	// The single letters are the abbreviations
//	// used by the String method's formatting.
//	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
//	ModeAppend                                     // a: append-only
//	ModeExclusive                                  // l: exclusive use
//	ModeTemporary                                  // T: temporary file; Plan 9 only
//	ModeSymlink                                    // L: symbolic link
//	ModeDevice                                     // D: device file
//	ModeNamedPipe                                  // p: named pipe (FIFO)
//	ModeSocket                                     // S: Unix domain socket
//	ModeSetuid                                     // u: setuid
//	ModeSetgid                                     // g: setgid
//	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
//	ModeSticky                                     // t: sticky
//
//	// Mask for the type bits. For regular files, none will be set.
//	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
//
//	ModePerm FileMode = 0777 // Unix permission bits
//)
//chmod [-R] xyz filename|dirname  -R:表示递归修改
//
//Linux文件的基本权限有9个,
// 分别是owner,group,others三种身份各自的read,write,execute权限,3个一组.
// owner/group/others 即 拥有者/群组/其他
// read/write/execute 即 可读/可写/可执行
// 可以用数字代表各个权限:
// x:1
// w:2
// r:4
//利用2进制表示,1代表有此权限,0表示没有此权限:
//---: 000 => 0
//--x: 001 => 1
//-w-: 010 => 2
//-wx: 011 => 3
//r--: 100 => 4
//r-x: 101 => 5
//rw-: 110 => 6
//rwx: 111 => 7
//
//因此:若将文件的权限修改为rwxrwx---,则对应的数字为 770
func OpenFile(filePath string, fileName string, flag int) *os.File {
	CreatePath(filePath)
	fileFullName := filePath + fileName
	// os.O_WRONLY|os.O_CREATE|os.O_APPEND
	// 以只写方式打开文件
	// 如果不存在，则创建
	// 在文件末尾追加
	file, err := os.OpenFile(fileFullName, flag, 0766)
	think.IsNil(err)
	// 本次打开的文件如果close,log日志无法写入
	// file.Close()
	return file
}

func CreatePath(path string) {
	// 如果不存在,则创建目录
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			think.IsNil(err)
		} else {
			// 创建文件夹
			err := os.MkdirAll(path, os.ModePerm)
			think.IsNil(err)
		}
	}
}

// basePath: ./XXX/YYY
// return: /AAA/BBB/XXX/YYY
func GetAbsPathWith(basePath string) string {
	basePath = strings.Replace(basePath, "/", string(os.PathSeparator), -1)
	path, err := filepath.Abs(basePath)
	think.IsNil(err)
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path
	} else {
		return path + string(os.PathSeparator)
	}
}

func LS(dir string) ([]string ,error){
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{info.Name()}, nil
	}
	infos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return getNameFromInfo(infos), nil
}

func getNameFromInfo(infos []os.FileInfo) []string {
	names := make([]string, 0)
	for _, info := range infos {
		names = append(names, info.Name())
	}
	return names
}

// find ./path -name '*.suffix'
func ListFile(path string, suffix string) []string {
	allFile := make([]string, 0)
	filePath := GetAbsPathWith(path)
	// 遍历filePath下的所有文件以及目录,ls .sql 文件
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			allFile = append(allFile, path)
			return nil
		} else {
			return nil
		}
		return nil
	})

	return allFile
}

func CopyFile(dst, src string) error {
	file, err := os.OpenFile(src, os.O_RDONLY, 0440)
	if err != nil {
		return err
	}
	defer file.Close()

	fileNew, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer fileNew.Close()

	_, err = io.Copy(fileNew, file)
	return err
}

func ReadLargeFile(fileName string, cacheSize int, ft func(bs []byte)) error{
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		bs := make([]byte, cacheSize)
		n, err := f.Read(bs)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		ft(bs[:n])
	}
	return nil
}