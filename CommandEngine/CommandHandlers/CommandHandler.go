package commandhandlers

import "github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"

type ICommandHandler interface {
	HandleAsync(command commands.ICommand)
}
