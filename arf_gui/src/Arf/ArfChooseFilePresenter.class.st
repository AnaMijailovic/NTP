"
Presenter for choosing a file.
"
Class {
	#name : #ArfChooseFilePresenter,
	#superclass : #SpPresenter,
	#instVars : [
		'pathInput',
		'chooseFileBtn'
	],
	#category : #'Arf-presenters'
}

{ #category : #initialization }
ArfChooseFilePresenter class >> defaultLayout [
 
 ^SpBoxLayout newHorizontal 
		   				add: #pathInput expand: true fill: true padding: 2;
						add: #chooseFileBtn withConstraints: [ :c | c width: 70; height: 25].
]

{ #category : #specs }
ArfChooseFilePresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfChooseFilePresenter >> initializePresenters [

	pathInput := self newTextInput autoAccept: true.
	chooseFileBtn := self newButton label: 'Choose';
											 action: [ self openChooseFileDialog ].
]

{ #category : #'open dialog' }
ArfChooseFilePresenter >> openChooseFileDialog [
	
	|fd selected |

	fd := FileDialogWindow basicNew
        previewType: nil;
        initialize ;
        title: 'Choose file';
        answerPathName.
	 fd cancel.
	
	selected := ''.
	(self currentWorld openModal: fd) answer.
	[selected := fd selectedDirectory absolutePath asUrl asString] 
		on: Error 
		do: [ ^ selected ] .

	selected := selected allButFirst: 8.  "remove file:///"
	fd selectedFileName ifNotNil: [  
   			selected := selected , '/' , fd selectedFileName asString. "add file name"].
   selected := ZnResourceMetaUtils decodePercent: selected. "remove %"

   pathInput text: selected.
]

{ #category : #accessing }
ArfChooseFilePresenter >> path [ 

	|path isValidPath |

	path := pathInput text.
	path = '' ifTrue: [ self error: 'You must enter a path' ]. 
	
	isValidPath := path asFileReference.
	isValidPath exists ifFalse: [ self error: 'Path is invalid' ] .
	
	^ path
]
