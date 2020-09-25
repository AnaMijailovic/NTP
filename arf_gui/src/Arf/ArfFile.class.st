"
Class the represents a file.
"
Class {
	#name : #ArfFile,
	#superclass : #Object,
	#instVars : [
		'name',
		'fullPath',
		'isDir',
		'size',
		'fileType',
		'created',
		'modified',
		'accessed'
	],
	#category : #'Arf-model'
}

{ #category : #accessing }
ArfFile >> accessed [
	^ accessed
]

{ #category : #accessing }
ArfFile >> accessed: aAccessed [
	accessed := aAccessed
]

{ #category : #accessing }
ArfFile >> created [ 
	^ created
]

{ #category : #accessing }
ArfFile >> created: aCreated [
	created := aCreated
]

{ #category : #accessing }
ArfFile >> fileType [
	^ fileType 
]

{ #category : #accessing }
ArfFile >> fileType: aFileType [
	fileType := aFileType 
]

{ #category : #accessing }
ArfFile >> fullPath [
	^ fullPath
]

{ #category : #accessing }
ArfFile >> fullPath: aFullPath [
	fullPath := aFullPath 
]

{ #category : #initialization }
ArfFile >> initialize [
	
	self class initializeSlots: self .
	super initialize .
	name := 'namee' .
	fullPath := '' .
	isDir := false . 
	size := 0 .
	fileType := '' .
	created := '2018-01-01' asDate .
	modified := '2018-01-01' asDate . 
	accessed := '2018-01-01' asDate .
]

{ #category : #accessing }
ArfFile >> isDir [
	^ isDir
]

{ #category : #accessing }
ArfFile >> isDir: aIsDir [
	isDir := aIsDir 
]

{ #category : #accessing }
ArfFile >> modified [ 	
	^ modified
]

{ #category : #accessing }
ArfFile >> modified: aModified [
	modified := aModified
]

{ #category : #accessing }
ArfFile >> name [
	^ name
]

{ #category : #accessing }
ArfFile >> name: aName [
	name := aName
]

{ #category : #accessing }
ArfFile >> size [ 	
	^ size
]

{ #category : #accessing }
ArfFile >> size: aSize [	
	size := aSize
]
