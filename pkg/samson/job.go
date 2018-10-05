package samson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Job a Samson job configuration
type Job struct {
	client        *http.Client
	deployTimeout int    `required:"true"`
	pollInterval  int    `required:"true"`
	project       string `required:"true"`
	reference     string `required:"true"`
	stage         string `required:"true"`
	token         string `required:"true"`
	url           string `required:"true"`
}

type triggerDeployResponse struct {
	ID  json.Number `json:"id"`
	URL string      `json:"url"`
}

type deployStatusResponse struct {
	Deploy deploy `json:"deploy"`
}

type deploy struct {
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

var samsonSuccessTermSet = map[string]struct{}{
	"succeeded": {},
}

var samsonDoneTermSet = map[string]struct{}{
	"succeeded": {},
	"failed":    {},
	"errored":   {},
	"cancelled": {},
}

// NewJob initializes a Samson job
func NewJob(options ...Option) *Job {
	zerolog.TimeFieldFormat = ""

	j := Job{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	for i := range options {
		options[i](&j)
	}

	return &j
}

func (j *Job) triggerDeploy() (deployID string, err error) {
	url := fmt.Sprintf("%s/projects/%s/stages/%s/deploys", j.url, j.project, j.stage)
	data := []byte(fmt.Sprintf(`{"deploy": {"reference": %q}}`, j.reference))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", j.token))

	resp, err := j.client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("Non-200 recieved trying to trigger deploy")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := triggerDeployResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("JOB - SAMSON: Triggered deploy (%s)", res.URL)

	return string(res.ID), nil
}

func (j *Job) getDeployStatus(deployID string) (d deploy, err error) {
	log.Debug().Msgf("JOB - SAMSON: Fetching deploy (id: %s) status...", deployID)

	url := fmt.Sprintf("%s/projects/%s/deploys/%s", j.url, j.project, deployID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return deploy{}, errors.New("JOB - SAMSON: Error trying to fetch deploy status")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", j.token))

	resp, err := j.client.Do(req) // TODO: check response status
	if err != nil {
		return deploy{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deploy{}, err
	}

	res := deployStatusResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return deploy{}, err
	}

	return res.Deploy, nil
}

func (j *Job) pollForResult(deployID string) (deploy, error) {
	c := make(chan deploy, 1)

	ticker := time.NewTicker(time.Duration(j.pollInterval) * time.Second)
	go func() {
		for range ticker.C {
			log.Debug().Msg("JOB - SAMSON: Polling for build result...")

			d, err := j.getDeployStatus(deployID)

			if err != nil {
				log.Error().Msgf("JOB - SAMSON: %s", err.Error())
				continue
			}

			log.Debug().Msgf("JOB - SAMSON: Deploy %s", d.Summary)

			if contains(samsonDoneTermSet, d.Status) {
				c <- d
			}
		}
	}()

	select {
	case d := <-c:
		ticker.Stop()
		return d, nil
	case <-time.After(time.Duration(j.deployTimeout) * time.Minute):
		ticker.Stop()
		return deploy{}, errors.New("timed out waiting for deploy result")
	}
}

func (j *Job) reportSuccess() {
	log.Debug().Msgf("JOB - SAMSON: Reporting success for deploy.")

	// report success
	os.Exit(0)
}

func (j *Job) reportFailure() {
	log.Debug().Msgf("JOB - SAMSON: Reporting failure for deploy.")

	// report failure
	os.Exit(1)
}

func (j *Job) reportStatus(status string) {
	if contains(samsonSuccessTermSet, status) {
		j.reportSuccess()
	}
	j.reportFailure()
}

// Execute starts a Samson deploy job
func (j *Job) Execute() {
	var err error
	var deployID string
	var d deploy

	if deployID, err = j.triggerDeploy(); err == nil {
		if d, err = j.pollForResult(deployID); err == nil {
			j.reportStatus(d.Status)
		}
	}
	if err != nil {
		log.Fatal().Msgf("JOB - SAMSON: %s", err.Error())
	}
}

func contains(set map[string]struct{}, item string) bool {
	_, ok := set[item]
	return ok
}
