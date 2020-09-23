Class {
	#name : #ArfRecoverPresenter,
	#superclass : #SpPresenter,
	#instVars : [
		'pathInput',
		'recoverBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfRecoverPresenter class >> defaultLayout [
		^ SpBoxLayout  newVertical
			add: '' withConstraints: [ :c | c height: 25];
			add: 'Recovery file path:' withConstraints: [ :c | c height: 25];
			add: #pathInput ;
			add: #recoverBtn withConstraints: [ :c | c height: 40];
			yourself. 

]

{ #category : #specs }
ArfRecoverPresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfRecoverPresenter >> initializePresenters [
	pathInput := self newTextInput autoAccept: true.
	
	recoverBtn := self newButton label: 'Recover'.
	recoverBtn action: [ self submitRecover ].
]

{ #category : #initialization }
ArfRecoverPresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Recover';
		initialExtent: 300@200;
		centered
]

{ #category : #accessing }
ArfRecoverPresenter >> path [ 
	|path isValidPath |

	path := pathInput text.
	path = '' ifTrue: [ self error: 'You must enter a path' ]. 
	
	isValidPath := path asFileReference.
	isValidPath exists ifFalse: [ self error: 'Path is invalid' ] .
	
	^ path
]

{ #category : #submit }
ArfRecoverPresenter >> submitRecover [
	
	| path response content|
	
	[path := self path ]
  			on: Error
  			do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
				
	[response := (ZnClient new)
   		url: 'http://localhost:8080/api/recover';
		queryAt: 'path' put: path ;
   		put;
   		response]
		on: ConnectionTimedOut 
		do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
			
   content := response contents. 
	ArfInformationDialog new message: content asString; openModalWithSpec.
				
	
	self inform: 'Submit', path
]
