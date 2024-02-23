package modules

import (
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/component"
	taskcontroller "projects/LDmitryLD/task-service/gateway/internal/modules/task/controller"
	taskservice "projects/LDmitryLD/task-service/gateway/internal/modules/task/service"
	usercontroller "projects/LDmitryLD/task-service/gateway/internal/modules/user/controller"
	userservice "projects/LDmitryLD/task-service/gateway/internal/modules/user/service"
)

type Controllers struct {
	Task taskcontroller.Tasker
	User usercontroller.Userer
}

func NewControllers(user userservice.Userer, task taskservice.Tasker, components *component.Components) *Controllers {
	taskController := taskcontroller.NewTaskController(task, components)
	userController := usercontroller.NewUserController(user, components)

	return &Controllers{
		Task: taskController,
		User: userController,
	}
}
