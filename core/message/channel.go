package message

import message2 "github.com/smtdfc/nagare/pkgs/messages"

type ReadOnlyChannel <-chan message2.Message
type WriteOnlyChannel chan<- message2.Message
type Channel chan message2.Message
