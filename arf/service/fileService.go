package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type WalkFunc func(path string, info os.FileInfo, err error) error
type GetLabelFunc func(filePath string) (string, error)


// *************************************************************************
//							   CREATE TREE
// *************************************************************************

func CreateTree(root string) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	if stat.IsDir() {
		rootNode.Children = getChildren(root)
	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}
func getChildren(parentPath string) []*model.Node {
	// ReadDir reads the directory named by dirname
	// and returns a list of directory entries sorted by filename.
	children, err := ioutil.ReadDir(parentPath)

	if err != nil {
		fmt.Println(err)
	}

	childrenNodes := make([]*model.Node, 0)

	for _, file := range children {

		absPath, _ := filepath.Abs(parentPath + string(os.PathSeparator) + file.Name() )
		stat, err := os.Stat(absPath)

		if err != nil {
			fmt.Println(err)
		}

		childNode := createFileTreeNode(absPath, stat)
		childrenNodes = append(childrenNodes, &childNode)

		if stat.IsDir()  {
			childNode.Children = getChildren(absPath) // parallel
		} else {
			childNode.Children = make([]*model.Node, 0)
		}
	}

	return childrenNodes
}

func createFileTreeNode(filePath string, fileInfo os.FileInfo) model.Node {
	// for created and accessed dates
	// they are not provided in FileInfo
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	fileType, _ := getFileContentType(filePath)
	file := model.File{Name: fileInfo.Name(), FullPath: filePath, IsDir: fileInfo.IsDir(), Size: fileInfo.Size(),
		Created: timeStat.BirthTime(), Modified: fileInfo.ModTime(), Accessed: timeStat.AccessTime(),
		FileType: fileType }
	treeNode := model.Node{ Element: file }
	return treeNode

}

func getFileContentType(filePath string) (string, error) {
	out, _ := os.Open(filePath)
	defer out.Close()
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	// fmt.Println(contentType)
	return contentType, nil
}

// *************************************************************************
//							   CHART DATA
// *************************************************************************

func GetFileChartData(root string, chartType string) map[string]int64 {

	getLabelFunc := getLabelFunc(chartType)
	// closure
	fileTypesDict := make(map[string]int64)

	Walk(root, true, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		fileType, _ := getLabelFunc(path)
		size, ok := fileTypesDict[fileType]
		if ok {
			fileTypesDict[fileType] = size + info.Size()
		} else {
			fileTypesDict[fileType] = info.Size()
		}

		return nil
	})

	return fileTypesDict
}

func getLabelFunc(chartType string) GetLabelFunc {

	if chartType == "fileType"{
		return getFileContentType
	} else if chartType == "createdDate" {
		return getCreatedDateLabel
	} else if chartType == "modifiedDate" {
		return getModifiedDateLabel
	}else { // accessed date
		return getAccessedDateLabel
	}
}

func getCreatedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	createdYear := timeStat.BirthTime().Year()

	return strconv.Itoa(createdYear), nil
}

func getModifiedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	modifiedYear := timeStat.ModTime().Year()

	return strconv.Itoa(modifiedYear), nil
}

func getAccessedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	accessedYear := timeStat.ModTime().Year()

	return strconv.Itoa(accessedYear), nil
}

func Walk(root string, recursive bool, walkFn WalkFunc) error {

	stat, err := os.Stat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}

	if !stat.IsDir() {
		return walkFn(root, stat, nil)
	}

	// if there is no error
	// first call walkFn for the root
	walkFn(root, stat, nil)

	// ReadDir reads the directory named by dirname
	// and returns a list of directory entries sorted by filename.
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {

		abs, _ := filepath.Abs(root + string(os.PathSeparator) + file.Name() )
		stat, err := os.Stat(abs)

		if err != nil {
			return err
		}

		if stat.IsDir() && recursive {
			Walk(abs, recursive, walkFn)
		} else { // If it is a file
			walkFn(abs, stat, nil)
		}
	}

	return nil
}


// *************************************************************************
//								DELETE
// *************************************************************************

func DeleteFiles(path string, recursive bool, empty bool, createdBefore time.Time, notAccessedAfter time.Time) {
	tree := CreateTree(path)

	postorderDelete(tree.Root, empty, createdBefore, notAccessedAfter)

}

func postorderDelete(node *model.Node, empty bool, createdBefore time.Time, notAccessedAfter time.Time ) {

	for _, file := range node.Children {
		postorderDelete(file, empty, createdBefore, notAccessedAfter)
	}

	file := node.Element.(model.File)
	stat, _ := os.Stat(file.FullPath)
	size := stat.Size()

	if ( empty && size == 0 ) ||
		(createdBefore.After(file.Created) && (!file.IsDir || file.Size == 0)) ||
		(notAccessedAfter.After(file.Accessed) && (!file.IsDir || file.Size == 0)){
		os.Remove(file.FullPath)
		fmt.Println("Removed: ", file.FullPath)
	}
}

// *************************************************************************
//								REORGANIZE
// *************************************************************************

func ReorganizeFiles(src string, dest string, recursive bool, fileType bool, fileSize int64, createdDate string) {


	Walk(src, recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFolderName := generateFolderName(path, fileType, fileSize, createdDate)
		newFolderName = strings.Replace(newFolderName, "/", "_", -1)

		newFolderPath := dest + string(os.PathSeparator) + newFolderName
		fmt.Println("New folder path: ", newFolderPath)

		// check if directory already exists
		if _, err := os.Open(newFolderPath); err != nil {
			err := os.Mkdir(newFolderPath, 0755)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Move file
		err = os.Rename(path, newFolderPath + string(os.PathSeparator) + info.Name())

		// TODO Write recovery data to a file

		return nil
	})

}

func generateFolderName(path string, fileType bool, fileSize int64, createdDate string) string {
	if fileType {
		return generateFolderNameByType(path)
	} else if fileSize != 0 {
		return generateFolderNameBySize(path, fileSize)
	} else {
		return generateFolderNameByDate(path, createdDate)
	}
}

func generateFolderNameByType(path string) string {
	fileType, _ := getFileContentType(path)
	return fileType
}

func generateFolderNameBySize(path string, step int64) string {
	var label string
	stat, _ := os.Stat(path)
	size := stat.Size()

	groupNum := size / step
	fmt.Println("Size: ", size)
	fmt.Println("GroupNum: ", groupNum)

	if groupNum == 0 {
		label = strconv.FormatInt(step*groupNum, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	} else {
		label = strconv.FormatInt(step*groupNum + 1, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	}

	return label
}

func generateFolderNameByDate(path string, dateType string) string {

	timeStat, _ := times.Stat(path)
	var folderName string

	if dateType == "d" {
		folderName = timeStat.BirthTime().Format("02-01-2006")
	} else if dateType == "m" {
		folderName = timeStat.BirthTime().Format("01-2006")
	} else { // year
		folderName = timeStat.BirthTime().Format("2006")
	}

	return folderName
}