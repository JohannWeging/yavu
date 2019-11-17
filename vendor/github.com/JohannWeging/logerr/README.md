# logerr
Annotate errors with log fields.
logerr is intended to bridge the gap between `logrus.WithFields` and error wrapping libraries.

## Example
```go
package main

import (
	"github.com/JohannWeging/logerr"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := job()
	if err != nil {
		fields := logerr.GetFields(err)
		log.WithFields(fields).Errorf("job failed: %s", err)
	}
}

func job() error {
	err := task()
	if err != nil {
		err = logerr.WithField(err, "jobID", "0")
		err = errors.Annotate(err, "task failed")
		return err
	}
	return nil
}

func task() error {
	return errors.New("cause")
}

```
```
ERRO[0000] job failed: task failed: cause                jobID=0
```

## Implemented Interfaces

`logerr.Error` implements 2 interfaces:
```go

// github.com/juju/errors
type wrapper interface {
	Underlying() error
}

// github.com/pkg/errors
type causer interface {
	Cause() error
}
```
