package logservices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Handler struct {
	metricsClient *metricsv.Clientset
	clientSet     kubernetes.Interface
}

func NewHandler(metricsClient *metricsv.Clientset, clientSet kubernetes.Interface) *Handler {
	return &Handler{
		metricsClient: metricsClient,
		clientSet:     clientSet,
	}
}

func (h *Handler) GetMetrics(c echo.Context) error {
	userID := c.Get("user_id").(string)

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(200)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := GetUserWorkspaceMetrics(c.Request().Context(), h.metricsClient, userID)
			if err != nil {
				fmt.Fprintf(c.Response(), "event: error\ndata: %s\n\n", err.Error())
				c.Response().Flush()
				continue
			}

			data, _ := json.Marshal(metrics)
			fmt.Fprintf(c.Response(), "event: metrics\ndata: %s\n\n", data)
			c.Response().Flush()

		case <-c.Request().Context().Done():
			return nil
		}
	}
}
func (h *Handler) StreamStorageMetrics(c echo.Context) error {
	userID := c.Get("user_id").(string)

	namespace := fmt.Sprintf("ws-%s", userID)

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			podList, err := h.clientSet.CoreV1().Pods(namespace).List(
				c.Request().Context(),
				metav1.ListOptions{},
			)
			if err != nil {
				fmt.Fprintf(c.Response(), "event: error\ndata: %s\n\n", err.Error())
				c.Response().Flush()
				continue
			}

			for _, pod := range podList.Items {
				metrics, err := GetPodStorageFromKubelet(
					c.Request().Context(),
					h.clientSet,
					namespace,
					pod.Name,
				)
				if err != nil {
					continue
				}

				data, _ := json.Marshal(metrics)
				fmt.Fprintf(c.Response(), "event: storage\ndata: %s\n\n", data)
				c.Response().Flush()
			}
		}
	}
}
