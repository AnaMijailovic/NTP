Class {
	#name : #ArfInformationDialog,
	#superclass : #SpPresenter,
	#instVars : [
		'message'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfInformationDialog class >> defaultLayout [
	
	^ SpBoxLayout  newVertical
			add: '' withConstraints: [ :c | c height: 25];
			add: #message withConstraints: [ :c | c height: 100];
			yourself. 
]

{ #category : #specs }
ArfInformationDialog class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfInformationDialog >> initializeDialogWindow: aDialog [
	aDialog
		closeOnBackdropClick: false;
		addButton: 'Ok'
			do: [ :presenter | 
			presenter close ].
]

{ #category : #initialization }
ArfInformationDialog >> initializePresenters [

	message := self newLabel label: ''.
	

]

{ #category : #initialization }
ArfInformationDialog >> initializeWindow: aWindowPresenter [
	
	aWindowPresenter 
		title: 'Information Dialog';
		initialExtent: 300@220;
		centered
]

{ #category : #accessing }
ArfInformationDialog >> message: aMessage [

	message label: aMessage
]
