Class {
	#name : #ArfDeletePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'emptyCheckBox',
		'creationCheckBox',
		'creationDateInput',
		'accessedCheckBox',
		'accessedDateInput',
		'submitBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfDeletePresenter class >> defaultLayout [
	^ SpBoxLayout   newVertical 
			add: '' withConstraints: [ :c | c height: 25];
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #emptyCheckBox withConstraints: [ :c | c height: 45];
			add: #creationCheckBox withConstraints: [ :c | c height: 45];
			add: #creationDateInput withConstraints: [ :c | c height: 45];
			add: #accessedCheckBox withConstraints: [ :c | c height: 45];
			add: #accessedDateInput;
			add: #submitBtn withConstraints: [ :c | c height: 40];
		   yourself.
]

{ #category : #specs }
ArfDeletePresenter class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #initialization }
ArfDeletePresenter >> initializePresenters [
	pathInput := self newTextInput autoAccept: true.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	emptyCheckBox := self newCheckBox label: 'Delete empty files and folders'.
	creationCheckBox  := self newCheckBox label: 'Delete files and folders created before:'.
	accessedCheckBox := self newCheckBox label: 'Delete files and folders not accessed after:'.
												
	creationDateInput := self instantiate: SpDatePresenter.
	creationDateInput display: [ :item | item yyyymmdd ].
	
	accessedDateInput := self instantiate: SpDatePresenter.
	accessedDateInput display: [ :item | item yyyymmdd ].
	
	submitBtn := self newButton label: 'Delete'.
	submitBtn action: [ self submitDelete ].
]

{ #category : #initialization }
ArfDeletePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Delete';
		initialExtent: 300@430;
		centered
]

{ #category : #initialization }
ArfDeletePresenter >> submitDelete [
	| enteredPath recursive empty createdBefore notAccessedAfter response content isValidPath |

	recursive := recursiveCheckBox state asString.
	empty := emptyCheckBox  state asString.
	enteredPath := pathInput text.
	
	createdBefore := '01-01-0001' .
	notAccessedAfter := '01-01-0001' .
	creationCheckBox state  ifTrue: [ createdBefore := creationDateInput date printFormat: #(1 2 3 $- 2 1)] .
	accessedCheckBox state ifTrue: [ notAccessedAfter := accessedDateInput date printFormat: #(1 2 3 $- 2 1)] .
	
	"Validate that at least one criteria is provided"
	createdBefore = '01-01-0001' ifTrue: [ 
		notAccessedAfter = '01-01-0001' ifTrue: [ 
			emptyCheckBox state ifFalse: [ 
				self inform: 'Critera not provided'.
				 ]
			 ]
		] .
	
	"Validate path"
	enteredPath = '' ifTrue: [ ArfInformationDialog new message: 'You must enter a path';
		openWithSpec. ^ self ]. 
	isValidPath := enteredPath asFileReference.
	isValidPath exists ifFalse: [ ArfInformationDialog  new message: 'The path is invalid'; 
		openWithSpec. ^ self] .
	
	response := (ZnClient new)
   		url: 'http://localhost:8080/api/deleteFiles';
		queryAt: 'path' put: enteredPath ;
		queryAt: 'recursive' put: recursive;
		queryAt: 'empty' put: empty;
		queryAt: 'createdBefore' put: createdBefore ;
		queryAt: 'notAccessedAfter' put: notAccessedAfter;
   		delete;
   		response.
   content := response contents. 
	ArfInformationDialog new message: content asString; openWithSpec.
]
