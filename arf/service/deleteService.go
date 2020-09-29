package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"os"
	"sync"
)


// Parallel version
// Create a tree that represents file structure.
// Use post-order depth-first search algorithm to delete
// all files and folders that match the selected criteria.
// Returns the number of deleted files.
func DeleteFiles(deleteData *model.DeleteData) *int {

	tree := CreateTree(deleteData.Path, deleteData.Recursive)
	filesDeleted := 0

	// Create channel used for locking filesDeleted var
	lock := make(chan bool, 1)

	filesDel := &filesDeleted
	postorderDelete(tree.Root, deleteData, &filesDeleted, lock, nil)

	return filesDel
}

// Parallel version
// Post-order depth-first search algorithm to delete files.
// Multiple criteria are connected with 'or' operator.
// Returns the number of deleted files.
func postorderDelete(node *model.Node, deleteData *model.DeleteData, filesDeleted *int, lock chan bool, done *sync.WaitGroup)  {
	// Create wait group
	var wg sync.WaitGroup

	for _, file := range node.Children {
		wg.Add(1)
		go postorderDelete(file, deleteData, filesDeleted, lock, &wg)
	}

	// Wait for all child nodes to be processed
	wg.Wait()
	
	file := node.Element.(model.File)
	stat, _ := os.Stat(file.FullPath)
	size := stat.Size()

	// Check deletion criteria
	if ( deleteData.Empty && size == 0 ) ||
		(deleteData.CreatedBefore.After(file.Created) && (!file.IsDir || file.Size == 0)) ||
		(deleteData.NotAccessedAfter.After(file.Accessed) && (!file.IsDir || file.Size == 0)){
		err := os.Remove(file.FullPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Removed: ", file.FullPath)
			lock <- true
			*filesDeleted++
			<- lock
		}
	}

	if done != nil {
		done.Done()
	}

}

// Serial version
// Create a tree that represents file structure.
// Use post-order depth-first search algorithm to delete
// all files and folders that match the selected criteria.
// Returns the number of deleted files.
func DeleteFilesS(deleteData *model.DeleteData) *int {
	tree := CreateTreeS(deleteData.Path, deleteData.Recursive)
	filesDeleted := 0
	filesDel := &filesDeleted
	filesDel = postorderDeleteS(tree.Root, deleteData, &filesDeleted)

	return filesDel
}

// Serial version
// Post-order depth-first search algorithm to delete files.
// Multiple criteria are connected with 'or' operator.
// Returns the number of deleted files.
func postorderDeleteS(node *model.Node, deleteData *model.DeleteData, filesDeleted *int) *int {

	for _, file := range node.Children {
		filesDeleted = postorderDeleteS(file, deleteData, filesDeleted)
	}

	file := node.Element.(model.File)
	stat, _ := os.Stat(file.FullPath)
	size := stat.Size()

	if ( deleteData.Empty && size == 0 ) ||
		(deleteData.CreatedBefore.After(file.Created) && (!file.IsDir || file.Size == 0)) ||
		(deleteData.NotAccessedAfter.After(file.Accessed) && (!file.IsDir || file.Size == 0)){
		err := os.Remove(file.FullPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Removed: ", file.FullPath)
			*filesDeleted++
		}
	}

	return filesDeleted
}