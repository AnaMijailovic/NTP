package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type WalkFunc func(path string, info os.FileInfo, err error) error
type GetLabelFunc func(filePath string) (string, error)

// *************************************************************************
//				      WALK FUNCTION - PARALLEL VERSION
// *************************************************************************

// Links that helped:
// https://blog.golang.org/pipelines
// https://stackoverflow.com/questions/38170852/is-this-an-idiomatic-worker-thread-pool-in-go/38172204#38172204
// https://codereview.stackexchange.com/questions/114376/bruteforce-md5-password-cracker

	type Walk struct {
		root string
		recursive bool
		walkFn WalkFunc
		jobs chan string
		errs chan error
		jobsWg sync.WaitGroup // for each job, not worker!
		errsWg sync.WaitGroup // wait group for collecting errors
		poolSize int
	}

// Collect errors from the errs chanel
// Necessary to avoid blockage if the errs channel buffer is filled
	func (walk *Walk) GetErrors(errs *[]error){
		defer walk.errsWg.Done()
		for err := range walk.errs {
			*errs = append(*errs, err)
		}
	}

// Workers thread pool will be created.
// Worker takes jobs form the job channel
// and calls walkPath function for each of them.
	func (walk *Walk) startWorker() {
		for job := range walk.jobs {
			walk.walkPath(job)
			walk.jobsWg.Done()
		}
	}

// Calls walk.walkFunc for the file on the path arg.
// If path is a directory, reads all its subfiles
// and adds new jobs to the job channel.
// If some error happens adds that error to the
// errors channel.
func (walk *Walk) walkPath(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		walk.errs <- err
		return
	}

	err = walk.walkFn(path, stat, nil)
	if err != nil {
		walk.errs <- err
		return
	}

	if !stat.IsDir() {
		return
	}

	// ReadDir reads the directory named by dirname
	// and returns a list of directory entries sorted by filename.
	subfiles, err := ioutil.ReadDir(path)
	if err != nil {
		walk.errs <- err
		return
	}

	if path == walk.root || walk.recursive {

		for _, subfile := range subfiles {
			subfilePath, _ := filepath.Abs(path + string(os.PathSeparator) + subfile.Name() )

			walk.jobsWg.Add(1)
			select {
			case walk.jobs <- subfilePath: // add a new job
			default: // if jobs buffer is full, call walkPath func recursively
				walk.jobsWg.Done()
				walk.walkPath(subfilePath)
			}
		}
	}
}

// Function that starts goroutines
func (walk *Walk) startWalking() []error{

	// make channels
	walk.jobs = make(chan string, 1000)
	walk.errs = make(chan error, 1000)

	errs := make([]error, 0)

	// Start a goroutine for collecting errors from the errors channel
	walk.errsWg.Add(1)
	go walk.GetErrors(&errs)

	// Create workers pool
	for i := 0; i < walk.poolSize; i++ {
		go walk.startWorker()
	}

	// Add root to the jobs channel
	walk.jobsWg.Add(1)
	walk.jobs <- walk.root

	walk.jobsWg.Wait()   // wait for all jobs to finish
	close(walk.jobs)     // all jobs are done, close their channel
	close(walk.errs)     // all worker threads are done, so it's safe to close errors channel here
	walk.errsWg.Wait()   // wait for all errors to be collected from the channel

	return errs
}

// *************************************************************************
//				        WALK FUNCTION - SERIAL VERSION
// *************************************************************************

// Serial version of walk function
func WalkS(root string, recursive bool, walkFn WalkFunc) error {

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
			err = WalkS(abs, recursive, walkFn)
		} else { // If it is a file
			err = walkFn(abs, stat, nil)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// *************************************************************************
//							   CREATE TREE
// *************************************************************************

// Parallel version
// Crate tree representing files structure
func CreateTree(root string, recursive bool) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	var childrenChan = make(chan []*model.Node, 1)

	if stat.IsDir() {
		getChildren(root, recursive, childrenChan, nil)
		rootNode.Children = <- childrenChan

	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}

func getChildren(parentPath string, recursive bool, channel chan []*model.Node, done *sync.WaitGroup)  {

	children, err := ioutil.ReadDir(parentPath)
	var wg sync.WaitGroup

	if err != nil {
		fmt.Println(err)
	}

	if done != nil {
		defer done.Done()
	}

	childrenNodes := make([]*model.Node, len(children))

	for index, file := range children {

		absPath, _ := filepath.Abs(parentPath + string(os.PathSeparator) + file.Name() )
		stat, err := os.Stat(absPath)

		if err != nil {
			fmt.Println(err)
		}

		childNode := createFileTreeNode(absPath, stat)

		// Add to a specific index -> no data race here
		childrenNodes[index] = &childNode

		if stat.IsDir()  && recursive {

			var newChannel = make(chan []*model.Node, 1)
			wg.Add(1)

			go func(absPath string, wg *sync.WaitGroup, newChannel chan []*model.Node, childNode *model.Node) {
				getChildren(absPath, recursive, newChannel, wg)
				childNode.Children = <- newChannel

			}(absPath, &wg, newChannel, &childNode)

		} else {
			childNode.Children =  make([]*model.Node, 0)
		}
	}

	wg.Wait()
	channel <- childrenNodes
}

// Serial version
// Crate tree representing files structure
func CreateTreeS(root string, recursive bool) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	if stat.IsDir() {
		rootNode.Children = getChildrenS(root, recursive)
	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}

// Serial version
func getChildrenS(parentPath string, recursive bool) []*model.Node {

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
			childNode.Children = getChildrenS(absPath, recursive)
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

// Parallel version
// Get the data needed for the charts
func GetFileChartData(root string, chartType string) map[string]int64 {

	// closure
	fileTypesDict := make(map[string]int64)

	// Create channel used for locking fileTypesDict
	lock := make(chan bool, 1)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		fileType, _ := getLabel(chartType, path)

		lock <- true // lock fileTypesDict
		size, ok := fileTypesDict[fileType]
		if ok {
			fileTypesDict[fileType] = size + info.Size() / 1000
		} else {
			fileTypesDict[fileType] = info.Size() / 1000
		}
		<- lock // unlock fileTypesDict

		return nil
	}

	// Call walk function
	walk := Walk{ root: root, recursive: true, walkFn: walkFn, poolSize: 5}
	walkErrs := walk.startWalking()

	if walkErrs != nil {
		for _, err := range walkErrs {
			fmt.Println(err)
		}
	}

	return fileTypesDict
}

// Serial version
// Get the data needed for the charts
func GetFileChartDataS(root string, chartType string) map[string]int64 {

	// closure
	fileTypesDict := make(map[string]int64)

	WalkS(root, true, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		fileType, _ := getLabel(chartType, path)
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

// Make a label depending on the chartType
func getLabel(chartType string, filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	if chartType == "fileType"{
		return getFileContentType(filePath)
	} else if chartType == "createdDate" {
		createdYear := timeStat.BirthTime().Year()
		return strconv.Itoa(createdYear), nil
	} else if chartType == "createdDateM" {
		return getMonthLabel(timeStat.BirthTime()), nil
	} else if chartType == "modifiedDate" {
		modifiedYear := timeStat.ModTime().Year()
		return strconv.Itoa(modifiedYear), nil
	} else if chartType == "modifiedDateM" {
		return getMonthLabel(timeStat.ModTime()), nil
	} else if chartType == "accessedDate"{
		accessedYear := timeStat.ModTime().Year()
		return strconv.Itoa(accessedYear), nil
	} else {
		return getMonthLabel(timeStat.AccessTime()), nil
	}
}

// Make a label that consists of month and year
func getMonthLabel(date time.Time) string {
	year := date.Year()
	month := date.Month().String()
	return month + "-" + strconv.Itoa(year)
}
