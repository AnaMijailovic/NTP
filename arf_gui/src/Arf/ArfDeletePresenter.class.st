Class {
	#name : #ArfDeletePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'emptyRadioBtn',
		'creationRadioBtn',
		'creationDateInput',
		'accessedRadioBtn',
		'accessedDateInput',
		'deletionTypeBtns',
		'submitBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfDeletePresenter class >> defaultLayout [
	^ SpBoxLayout   newVertical 
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #emptyRadioBtn withConstraints: [ :c | c height: 45];
			add: #creationRadioBtn withConstraints: [ :c | c height: 45];
			add: #creationDateInput withConstraints: [ :c | c height: 45];
			add: #accessedRadioBtn withConstraints: [ :c | c height: 45];
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
	
	emptyRadioBtn := self newRadioButton label: 'Empty'.
	creationRadioBtn := self newRadioButton label: 'Created before:'.
	accessedRadioBtn := self newRadioButton label: 'Not accessed after:'.
	emptyRadioBtn associatedRadioButtons: { creationRadioBtn  . accessedRadioBtn }.
	deletionTypeBtns := Dictionary newFrom: { (#emptyRadioBtn -> emptyRadioBtn  ) . 
	                                   (#creationRadioBtn -> creationRadioBtn  ) .
												 (#accessedRadioBtn -> accessedRadioBtn  ) }. 
												
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
	self inform: 'Submit delete'
]
