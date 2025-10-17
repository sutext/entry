package broker

const (
	EntryTopicPrefix = "entry:topic:"
	EntryUserPrefix  = "entry:user:"
)

func (b *broker) topicKey(topic string) string {
	return EntryTopicPrefix + topic
}

func (b *broker) userKey(uid int64) string {
	return EntryUserPrefix + string(uid)
}
