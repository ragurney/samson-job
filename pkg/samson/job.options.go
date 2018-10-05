package samson

// Option functional options for Poller
type Option func(*Job)

// WithDeployTimeout allows customization of token
func WithDeployTimeout(timeout int) Option {
	return func(j *Job) {
		j.deployTimeout = timeout
	}
}

// WithPollInterval allows customization of poll interval
func WithPollInterval(interval int) Option {
	return func(j *Job) {
		j.pollInterval = interval
	}
}

// WithProject allows customization of project
func WithProject(project string) Option {
	return func(j *Job) {
		j.project = project
	}
}

// WithReference allows customization of token
func WithReference(ref string) Option {
	return func(j *Job) {
		j.reference = ref
	}
}

// WithStage allows customization of stage
func WithStage(stage string) Option {
	return func(j *Job) {
		j.stage = stage
	}
}

// WithToken allows customization of token
func WithToken(token string) Option {
	return func(j *Job) {
		j.token = token
	}
}

// WithURL allows customization of stage
func WithURL(url string) Option {
	return func(j *Job) {
		j.url = url
	}
}
