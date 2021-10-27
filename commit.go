package bitbucket

import (
	"strconv"
	"time"
)

type CommitTime time.Time

type Author struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type Committer struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type Parent struct {
	Id        string `json:"id"`
	DisplayId string `json:"displayId"`
}

type Commit struct {
	Id                 string     `json:"id"`
	DisplayId          string     `json:"displayId"`
	Author             Author     `json:"author"`
	AuthorTimestamp    CommitTime `json:"authorTimestamp"`
	Committer          Committer  `json:"committer"`
	CommitterTimestamp CommitTime `json:"committerTimestamp"`
	Message            string     `json:"message"`
	Parents            []Parent   `json:"parents"`
}

func (t CommitTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *CommitTime) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.UnixMilli(q)
	return nil
}

func (t CommitTime) Unix() int64 {
	return time.Time(t).Unix()
}

func (t CommitTime) Time() time.Time {
	return time.Time(t).UTC()
}

func (t CommitTime) String() string {
	return time.Time(t).String()
}
