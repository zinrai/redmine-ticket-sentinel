package sentinel

#Status: {
	id:   int
	name: string
}

#Tracker: {
	id:   int
	name: string
}

#Project: {
	id:   int
	name: string
}

#User: {
	id:    int
	login: string
}

// Ticket Structure Definition
#Issue: {
	id:         int
	subject:    string
	status:     #Status
	tracker:    #Tracker
	project:    #Project
	assignedTo: #User
	author:     #User
	children?: [...#Issue]
}

// Basic Structure of Rule Evaluation
#Rule: {
	name:     string
	message:  string
	evaluate: bool
	mention:  string
}
