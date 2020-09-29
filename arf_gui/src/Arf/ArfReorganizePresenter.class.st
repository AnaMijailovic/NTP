"
Presenter for reorganize functionality.
"
Class {
	#name : #ArfReorganizePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'srcPathInput',
		'destPathInput',
		'recursiveCheckBox',
		'fileTypeRadioBtn',
		'creationDRadioBtn',
		'creationMRadioBtn',
		'creationYRadioBtn',
		'fileSizeRadioBtn',
		'sizeStepInput',
		'reorganizeTypeBtns',
		'submitBtn'
	],
	#category : #'Arf-presenters'
}

{ #category : #specs }
ArfReorganizePresenter class >> defaultLayout [

		^ SpBoxLayout   newVertical 
			add: '' withConstraints: [ :c | c height: 25];
			add: 'Source directory path: ' withConstraints: [ :c | c height: 25];
			add: #srcPathInput withConstraints: [ :c | c height: 25];
			add: 'Destination directory path: ' withConstraints: [ :c | c height: 25];
			add: #destPathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: 'Reorganize by: ' withConstraints: [ :c | c height: 25];
			add: #fileTypeRadioBtn withConstraints: [ :c | c height: 45];
			add: #creationDRadioBtn withConstraints: [ :c | c height: 45];
		   add: #creationMRadioBtn withConstraints: [ :c | c height: 45];
		   add: #creationYRadioBtn withConstraints: [ :c | c height: 45];
			add: #fileSizeRadioBtn withConstraints: [ :c | c height: 45];
			add: #sizeStepInput;
			add: #submitBtn withConstraints: [ :c | c height: 40];
		   yourself.
]

{ #category : #specs }
ArfReorganizePresenter class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #accessing }
ArfReorganizePresenter >> creationDate [
   | creationDate reorganizeBy |

	creationDate := ''.
	reorganizeBy  := ( reorganizeTypeBtns associations detect: [ :each | each value state ]) key.
	reorganizeBy = 'creationDay' ifTrue: [ creationDate := 'd'].
	reorganizeBy = 'creationMonth' ifTrue: [ creationDate := 'm'].
	reorganizeBy = 'creationYear' ifTrue: [ creationDate := 'y'].
	
	^ creationDate
]

{ #category : #accessing }
ArfReorganizePresenter >> fileSize [
   | fileSize reorganizeBy |

	fileSize := 0.
	reorganizeBy  := ( reorganizeTypeBtns associations detect: [ :each | each value state ]) key.
	reorganizeBy = 'fileSize' ifTrue: [ fileSize := sizeStepInput text asNumber ].
	fileSize < 0  ifTrue: [self error: 'File size must be a positive value'].
	
	^ fileSize 
]

{ #category : #initialization }
ArfReorganizePresenter >> initializePresenters [
	srcPathInput := ArfChooseFilePresenter new.
	destPathInput := ArfChooseFilePresenter new.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	fileTypeRadioBtn := self newRadioButton label: 'File type'.
	creationDRadioBtn := self newRadioButton label: 'Creation date - full date'.
	creationMRadioBtn := self newRadioButton label: 'Creation date - month and year'.
	creationYRadioBtn := self newRadioButton label: 'Creation date - year'.
	fileSizeRadioBtn := self newRadioButton label: 'File size'.
	fileTypeRadioBtn associatedRadioButtons: { creationDRadioBtn . fileSizeRadioBtn .
															creationMRadioBtn .  creationYRadioBtn }.
	reorganizeTypeBtns  := Dictionary newFrom: { (#fileType -> fileTypeRadioBtn  ) . 
	                                   (#creationDay -> creationDRadioBtn  ) .
	                                   (#creationMonth -> creationMRadioBtn  ) .
                                      (#creationYear -> creationYRadioBtn  ) .
												 (#fileSize -> fileSizeRadioBtn  ) }. 
	
	sizeStepInput := self newNumberInput autoAccept: true; minimum: 1.		
	sizeStepInput placeholder: 'Enter size step (MB)'; beInteger .							

	
	submitBtn := self newButton label: 'Reorgnize'.
	submitBtn action: [ self submitReorganize ].
]

{ #category : #initialization }
ArfReorganizePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Reorganize';
		initialExtent: 300@530;
		centered
]

{ #category : #submit }
ArfReorganizePresenter >> submitReorganize [
	
	| srcPath  destPath recursive fileType creationDate fileSize reorganizeBy content response|

	recursive := recursiveCheckBox state asString.
	
	reorganizeBy  := ( reorganizeTypeBtns associations detect: [ :each | each value state ]) key.
	fileType := reorganizeBy = 'fileType'. 
	
	[fileSize := self fileSize .
	 srcPath := srcPathInput path .
	 destPath := destPathInput path.
	 creationDate := self creationDate.]
  			on: Error
  			do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
	 
	[response := (ZnClient new)
   		url: 'http://localhost:8080/api/reorganize';
		queryAt: 'src' put: srcPath;
		queryAt: 'dest' put: destPath;
		queryAt: 'recursive' put: recursive;
		queryAt: 'fileType' put: fileType ;
		queryAt: 'createdDate' put: creationDate ;
		queryAt: 'fileSize' put: fileSize;
   		put;
   		response.]
		on: ConnectionTimedOut 
		do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
   content := response contents. 
	ArfInformationDialog new message: content asString; openModalWithSpec.


]
