package sentinel

#StoryWithoutChildren: #Rule & {
	name:    "子チケットなしストーリー"
	message: "子チケットが一つも無い In Progress のストーリーとなっています。"
	evaluate: {
		let issue = _issue
		issue.status.id == #StatusInProgress && !issue.children
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#SubTaskToStory: #Rule & {
	name:    "大規模サブタスク"
	message: "サブタクスのボリュームが大きいように思います。サブタクスからストーリーへ分離する分岐点にきていないか確認をお願いします。"
	evaluate: {
		let issue = _issue
		let hasLargeSubTask = {
			let subTasks = [for c in issue.children if c.tracker.id == #TrackerSubTask {c}]
			let large = [for t in subTasks if (len(t.children) >= #MaxSubTaskSize) && t.status.id != #StatusClosed {t}]
			len(large) > 0
		}
		hasLargeSubTask
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#ForgotUpdateStory: #Rule & {
	name:    "ストーリー更新忘れ"
	message: "ストーリーのステータスが New のまま子チケットが進行しています。"
	evaluate: {
		let issue = _issue
		issue.status.id == #StatusNew && len([for c in issue.children if c.status.id != #StatusNew {c}]) > 0
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#TooManySubTasks: #Rule & {
	name:    "過多なサブタクス"
	message: "ストーリーに対してサブタクスが多いように思います。"
	evaluate: {
		let issue = _issue
		let subTaskCount = len([for c in issue.children if c.tracker.id == #TrackerSubTask {c}])
		subTaskCount > #MaxSubTaskPerStory
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#StoryOneOwnTicket: #Rule & {
	name:    "自身の未完了タスク"
	message: "ストーリーに Resolved, Closed となっていない自身が担当のタスクが一つ残っています。"
	evaluate: {
		let issue = _issue
		let ownTasks = [for c in issue.children if c.assignedTo.id == issue.assignedTo.id && c.status.id <= #StatusInProgress {c}]
		len(ownTasks) == 1 && issue.status.id != #StatusNew
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.assignedTo.login}
		mapping.slackMention
	}
}

#StoryOneTeamMemberTicket: #Rule & {
	name:    "チームメンバーの未完了タスク"
	message: "ストーリーに Resolved, Closed となっていないチームメンバーが担当のタスクが一つ残っています。"
	evaluate: {
		let issue = _issue
		let teamTasks = [for c in issue.children if c.assignedTo.id != issue.assignedTo.id && c.status.id <= #StatusInProgress {c}]
		len(teamTasks) == 1 && issue.status.id != #StatusNew
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#SubTaskOneRemainInProgress: #Rule & {
	name:    "進行中タスク残り一つ"
	message: "サブタクスに In Progress のタスクが一つ残っています。"
	evaluate: {
		let issue = _issue
		let subTasks = [for c in issue.children if c.tracker.id == #TrackerSubTask && len(c.children) > 0 {
			let inProgressTasks = [for t in c.children if t.status.id == #StatusInProgress {t}]
			let newTasks = [for t in c.children if t.status.id == #StatusNew {t}]
			inProgress: len(inProgressTasks)
			new:        len(newTasks)
		}]
		len([for s in subTasks if s.inProgress == 1 && s.new == 0 {s}]) > 0
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#TaskChildTicket: #Rule & {
	name:    "タスクの子チケット"
	message: "タスクに子チケットがぶら下がっています。"
	evaluate: {
		let issue = _issue
		let taskWithChildren = [for c in issue.children if c.tracker.id == #TrackerTask && len(c.children) > 0 {c}]
		len(taskWithChildren) > 0
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}

#SubTaskChildTicket: #Rule & {
	name:    "サブタクスの子チケット"
	message: "サブタクスの子チケットにサブタクスがぶら下がっています。"
	evaluate: {
		let issue = _issue
		let hasSubTaskChild = [for c in issue.children if c.tracker.id == #TrackerSubTask {
			let subTaskChildren = [for t in c.children if t.tracker.id == #TrackerSubTask {t}]
			len(subTaskChildren) > 0
		}]
		len([for h in hasSubTaskChild if h {h}]) > 0
	}
	mention: {
		let mapping = #GetSlackMention & {redmineLogin: _issue.author.login}
		mapping.slackMention
	}
}
