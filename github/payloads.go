package github

type User struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
}

// Funnily enough, this also sometimes means a Pull Request.
type Issue struct {
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
}

type PullRequest struct {
	Url    string `json:"html_url"`
	Id     int    `json:"id"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Title  string `json:"title"`
	User   User   `json:"user"`
	Sender User   `json:"sender"`
}

type Comment struct {
	User User   `json:"user"`
	Body string `json:"body"`
}

type PayloadPullRequestOpened struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
}

type PayloadIssueCommentCreated struct {
	Action  string  `json:"action"`
	Comment Comment `json:"comment"`
	Sender  User    `json:"sender"`
	Issue   Issue   `json:"issue"`
}

type GenericPayload struct {
	Action string `json:"action"`
}
