package services

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func (k *K8sClient) StopScalling(c context.Context,namespace string,pods string) error{
	scale,err := k.Client.AppsV1().Deployments(namespace).GetScale(c,pods,metav1.GetOptions{})

	if err != nil {
		log.Printf("failed to scalling pods : %s",err)
		return fmt.Errorf("failed to scalling pods")
	}

	scale.Spec.Replicas = 0

	_ , err =  k.Client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), pods, scale, metav1.UpdateOptions{})
    return err
}


func (k *K8sClient) StartScalling(c context.Context,namespace string,pods string) error{
	scale,err := k.Client.AppsV1().Deployments(namespace).GetScale(c,pods,metav1.GetOptions{})

	if err != nil {
		log.Printf("failed to scalling pods : %s",err)
		return fmt.Errorf("failed to scalling pods")
	}

	scale.Spec.Replicas = 1

	_ , err =  k.Client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), pods, scale, metav1.UpdateOptions{})
    return err
}