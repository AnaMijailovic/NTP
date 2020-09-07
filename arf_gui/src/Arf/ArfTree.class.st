Class {
	#name : #ArfTree,
	#superclass : #Object,
	#instVars : [
		'root'
	],
	#category : #Arf
}

{ #category : #converting }
ArfTree >> asList [
    | list addChildren |
	 list := OrderedCollection new
				add: self root;
				yourself.
	 
	 addChildren := [ :parent | 
			parent children collect: [ : child |
				list add: child.
				addChildren value: child. ]].
	addChildren value: self root.
	 
	^ list
	
]

{ #category : #accessing }
ArfTree >> root [ 
	^ root
]
