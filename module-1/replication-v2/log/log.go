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

func (l *Log) NextIndex() int {
  l.indexMX.Lock()
  defer l.indexMX.Unlock()
  l.index++
  return l.index
}

func (l *Log) GetAll() []string {
  return l.log
}

func (l *Log) shiftInsert(pm pendingMessage) {
  i := 0
  // look for a place where to insert new pending message
  for i < len(l.pending) && l.pending[i].index < pm.index {
    i++
  }


  before := l.pending[0:i]
  after := l.pending[i:len(l.pending)]
  
  newPending := append([]pendingMessage{}, before...)
  newPending = append(newPending, pm)
  l.pending = append(newPending, after...)

}

func (l *Log) Process(index int, m string) {
  l.logMX.Lock()
  defer l.logMX.Unlock()

  pendingMessage := pendingMessage{
    index,
    m,
  }
  l.shiftInsert(pendingMessage)

  // add pending messages to the log if they come in expected index
  lastExpectedIndex := len(l.log)
  for _, p := range l.pending {
    if p.index == lastExpectedIndex { //expectly the next message in thes index
      l.log = append(l.log, p.message)
      lastExpectedIndex++
    } else if p.index <  lastExpectedIndex{ // message duplicate but might have new value
      l.log[p.index] = p.message 
    } else { // remove merged messages
      //copy(l.pending, l.pending[i:]) // remove all the merged elements
      break
    }
  }
}
