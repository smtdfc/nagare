package message

import message2 "github.com/smtdfc/nagare/shared/messages"

type ReadOnlyChannel <-chan message2.Message
type WriteOnlyChannel chan<- message2.Message
type Channel chan message2.Message
