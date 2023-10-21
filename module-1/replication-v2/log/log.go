package log

import (
  "sync"
)

type pendingMessage struct {
  index int
  message string
}

type Log struct {
  index int
  indexMX sync.Mutex

  log []string
  pending []pendingMessage
  logMX sync.Mutex
}

func NewLog() Log {
  return Log {
    index: -1,
  }
}

func (l *Log) nextIndex() int {
  l.indexMX.Lock()
  defer l.indexMX.Unlock()
  l.index++
  return l.index
}

func (l *Log) GetLog() []string {
  return l.log
}

func (l *Log) Process(index int, m string) {
  l.logMX.Lock()
  defer l.logMX.Unlock()

  pendingMessage := pendingMessage{
    index,
    m,
  }

  // append pending message using shift sort
  for i, p := range l.pending {
    if index < p.index {
      before := l.pending[0:i]
      after := l.pending[i:len(l.pending)]
      l.pending = append(append(before, pendingMessage), after...)
      break
    }
  }

  // add pending messages to the log if they come in expected index
  lastExpectedIndex := len(l.log)
  for _, p := range l.pending {
    if p.index == lastExpectedIndex { //expectly the next message in thes index
      l.log = append(l.log, p.message)
    } else if p.index <  lastExpectedIndex{ // message duplicate but might have new value
      l.log[p.index] = p.message 
    } else {
      break
    }
  }
}
