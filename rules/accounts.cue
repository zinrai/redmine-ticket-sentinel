package sentinel

// Account Mapping Type Definition
#AccountMapping: {
	[RedmineLogin=string]: {
		slack: string // Account name for Slack mentions ( atmark none )
	}
}

accounts: #AccountMapping & {
	"tanaka": {
		slack: "tanaka.taro"
	}
	"yamada": {
		slack: "yamada.hanako"
	}
	"suzuki": {
		slack: "suzuki.ichiro"
	}
	"sato": {
		slack: "sato.saburo"
	}
}

// Generating a Mention String
#GetSlackMention: {
	redmineLogin: string
	slackMention: accounts[redmineLogin].slack | redmineLogin // If there is no mapping, the Redmine login name is used as is.
}
