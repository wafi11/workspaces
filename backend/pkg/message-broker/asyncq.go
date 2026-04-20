package messagebroker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type TaskSchedulling struct {
	WorkspaceID string `json:"workspace_id"`
	Type string `json:"type"`
	Duration int `json:"duration"`
    TypeTimeDuration time.Duration `json:"type_time_duration"`
	Status string `json:"status"`
}

func TaskStartAndStopWorkspace(payload string, client *asynq.Client) error {
    var taskReq TaskSchedulling
    if err := json.Unmarshal([]byte(payload), &taskReq); err != nil {
        return fmt.Errorf("failed to process task")
    }

    duration := time.Duration(taskReq.Duration) * taskReq.TypeTimeDuration
    if duration < taskReq.TypeTimeDuration {
        duration = taskReq.TypeTimeDuration
    }

    task := asynq.NewTask(string(EventStopWorkspace), []byte(payload))
    _, err := client.Enqueue(task,
        asynq.ProcessIn(duration),
        asynq.TaskID(fmt.Sprintf("autostop:%s", taskReq.WorkspaceID)),
        asynq.Queue("default"),
    )
    if err != nil && !errors.Is(err, asynq.ErrTaskIDConflict) {
        return err
    }
    return nil
}

func CancelAutoStop(workspaceID string, inspector *asynq.Inspector) {
    taskID := fmt.Sprintf("autostop:%s", workspaceID)
    
    if err := inspector.DeleteTask("default", taskID); err != nil {
        log.Printf("[scheduler] failed cancel autostop task %s: %v", workspaceID, err)
    }
}

// func TaskStartAndStopWorkspace(payload string, client *asynq.Client) error {
// 	task := asynq.NewTask(string(EventStopWorkspace), []byte(payload))

// 	var taskReq TaskSchedulling
// 	if err := json.Unmarshal(task.Payload(), &taskReq); err != nil {
// 		log.Printf("failed to process task %s", err.Error())
// 		return fmt.Errorf("failed to process task")
// 	}

// 	taskID := fmt.Sprintf("stop-workspace:%s", taskReq.WorkspaceID)

// 	_, err := client.Enqueue(
// 		task,
// 		asynq.ProcessIn(time.Duration(taskReq.Duration)*time.Minute),
// 		asynq.TaskID(taskID),
// 		asynq.Retention(1*time.Hour),
// 	)
// 	if err != nil && !isTaskIDConflict(err) {
// 		return fmt.Errorf("failed to enqueue task: %w", err)
// 	}

// 	return nil
// }

// func isTaskIDConflict(err error) bool {
// 	return err != nil && err.Error() == asynq.ErrTaskIDConflict.Error()
// }