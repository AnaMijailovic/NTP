Class {
	#name : #ArfFileDetailsPresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'fileNameLabel',
		'fullPathLabel',
		'isDirLabel',
		'sizeLabel',
		'fileTypeLabel',
		'createdLabel',
		'modifiedLabel',
		'accessedLabel',
		'file'
	],
	#category : #Arf
}

{ #category : #specs }
ArfFileDetailsPresenter class >> defaultLayout [

	^ SpGridLayout  new
			beColumnNotHomogeneous;
			add: 'Name: ' at: 1@1; add: #fileNameLabel at: 2@1 span: 3@1;
			add: 'Full path: ' at: 1@2; add: #fullPathLabel at: 2@2 span: 3@2;
			add: 'Is directory: ' at: 1@3; add: #isDirLabel at: 2@3 span: 3@3;
			add: 'Size (bytes): ' at: 1@4; add: #sizeLabel at: 2@4 span: 3@4;
			add: 'FileType: ' at: 1@5; add: #fileTypeLabel at: 2@5 span: 3@5;
			add: 'Created: ' at: 1@6; add: #createdLabel at: 2@6 span: 3@6;
			add: 'Modified: ' at: 1@7; add: #modifiedLabel at: 2@7 span: 3@7;
			add: 'Accessed: ' at: 1@8; add: #accessedLabel at: 2@8 span: 3@8;
			yourself.
]

{ #category : #specs }
ArfFileDetailsPresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #copying }
ArfFileDetailsPresenter >> copyToClipboard [

	Clipboard clipboardText: (self model fullPath). 

	self inform: self model fullPath,' copied to clipboard'
]

{ #category : #'as yet unclassified' }
ArfFileDetailsPresenter >> fillPresenterWithModel [
	| cDate mDate aDate |
	
	cDate := self model created asDate printFormat: #(1 2 3 $- 1 1 ) .
	mDate := self model modified asDate printFormat: #(1 2 3 $- 1 1 ) .
	aDate := self model accessed asDate printFormat: #(1 2 3 $- 1 1 ) .
		
	fileNameLabel label: self model name .
	fullPathLabel label: self model fullPath .
	isDirLabel label: self model isDir asString .
	sizeLabel label: self model size asString .
	fileTypeLabel label: self model fileType .
	createdLabel label: cDate.
	modifiedLabel label: mDate .
	accessedLabel label: aDate .
	
	
]

{ #category : #initialization }
ArfFileDetailsPresenter >> initialize [ 

	self model: ArfFile new .
	super initialize
]

{ #category : #initialization }
ArfFileDetailsPresenter >> initializeDialogWindow: aDialog [
	aDialog
		closeOnBackdropClick: false;
		addButton: 'Copy path'
			do: [ self copyToClipboard ];
		addButton: 'Ok'
			do: [ :presenter | 
			presenter close ].
]

{ #category : #initialization }
ArfFileDetailsPresenter >> initializePresenters [ 
	
	file := ArfFile new initialize .
	fileNameLabel := self newLabel .
	fullPathLabel := self newLabel .
	isDirLabel := self newLabel .
	sizeLabel := self newLabel .
	fileTypeLabel := self newLabel .
	createdLabel := self newLabel .
	modifiedLabel := self newLabel .
	accessedLabel := self newLabel .

]

{ #category : #initialization }
ArfFileDetailsPresenter >> initializeWindow: aWindowPresenter [
	
	aWindowPresenter 
		title: 'File Details';
		initialExtent: 580@360;
		centered
]

{ #category : #initialization }
ArfFileDetailsPresenter >> modelChanged [

	file := self model.
	self fillPresenterWithModel
]
