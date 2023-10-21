package log

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestAddMessage(t *testing.T) {
  l := NewLog()

  l.Process(1, "second")
  l.Process(2, "third")
  l.Process(0, "first")

  assert.Equal(t, []string{"first", "second", "third"}, l.GetAll())
}
