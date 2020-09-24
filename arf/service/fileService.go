package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type WalkFunc func(path string, info os.FileInfo, err error) error
type GetLabelFunc func(filePath string) (string, error)


// *************************************************************************
//							   CREATE TREE
// *************************************************************************

func CreateTree(root string, recursive bool) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	if stat.IsDir() {
		rootNode.Children = getChildren(root, recursive)
	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}

func getChildren(parentPath string, recursive bool) []*model.Node {
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

		if stat.IsDir()  && recursive {
			childNode.Children = getChildren(absPath, recursive) // parallel
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
			return err
		}

		fileType, _ := getLabelFunc(path)
		size, ok := fileTypesDict[fileType]
		if ok {
			fileTypesDict[fileType] = size + info.Size() / 1000
		} else {
			fileTypesDict[fileType] = info.Size() / 1000
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
	err = walkFn(root, stat, nil)
	if err != nil {
		return err
	}

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
			err = Walk(abs, recursive, walkFn)
		} else { // If it is a file
			err = walkFn(abs, stat, nil)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
