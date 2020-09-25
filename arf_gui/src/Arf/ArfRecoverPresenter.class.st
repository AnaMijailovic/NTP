"
Presenter for recover functionality.
"
Class {
	#name : #ArfRecoverPresenter,
	#superclass : #SpPresenter,
	#instVars : [
		'pathInput',
		'recoverBtn'
	],
	#category : #'Arf-presenters'
}

{ #category : #initialization }
ArfRecoverPresenter class >> defaultLayout [
		^ SpBoxLayout  newVertical
			add: 'Recovery file path:' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: '';
			add: #recoverBtn withConstraints: [ :c | c height: 40];
			yourself. 

]

{ #category : #specs }
ArfRecoverPresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfRecoverPresenter >> initializePresenters [
	pathInput := ArfChooseFilePresenter new.
	
	recoverBtn := self newButton label: 'Recover'.
	recoverBtn action: [ self submitRecover ].
]

{ #category : #initialization }
ArfRecoverPresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Recover';
		initialExtent: 300@160;
		centered
]

{ #category : #submit }
ArfRecoverPresenter >> submitRecover [
	
	| path response content|
	
	[path := pathInput path ]
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
