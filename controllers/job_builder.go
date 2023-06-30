package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
apiVersion: v1
kind: Pod
metadata:
  name: kaniko
spec:
  initContainers:
    - name: download-files
      image: public.ecr.aws/aws-cli/aws-cli:latest
      env:
        - name: AWS_PROFILE
          value: default
        - name: AWS_ENDPOINT
          value: https://s3.fr-par.scw.cloud
        - name: AWS_BUCKET
          value: stage-cf-worker
      volumeMounts:
        - name: s3-config
          mountPath: /root/.aws
          readOnly: true
        - name: context
          mountPath: /context
      command:
        - aws
        - "--endpoint-url=$(AWS_ENDPOINT)"
      args:
        - s3
        - sync
        - "s3://$(AWS_BUCKET)"
        - "/context"
        - "--debug"
  containers:
    - name: kaniko
      image: gcr.io/kaniko-project/executor:latest
      args: ["--dockerfile=Dockerfile",
             "--context=/context",
             "--destination=clementreiffers/artist-worker:latest"]
      volumeMounts:
        - name: registry-credentials
          mountPath: /kaniko/.docker/
          readOnly: true
        - name: context
          mountPath: /context
  volumes:
    - name: registry-credentials
      projected:
        sources:
          - secret:
              name: docker-hub
              items:
                - key: .dockerconfigjson
                  path: config.json
    - name: s3-config
      projected:
        sources:
          - secret:
              name: s3-credentials
              items:
                - key: credentials
                  path: credentials
          - configMap:
              name: aws-config
              items:
                - key: config
                  path: config
    - name: context
      emptyDir: {}
  restartPolicy: Never
*/

func create_job_builder() corev1.Pod {
	env := []corev1.EnvVar{
		{
			Name:  "AWS_PROFILE",
			Value: "default",
		},
		{
			Name:  "AWS_ENDPOINT",
			Value: "https://s3.fr-par.scw.cloud",
		},
		{
			Name:  "AWS_BUCKET",
			Value: "stage-cf-worker",
		},
	}
	volumeMounts := []corev1.VolumeMount{{Name: "s3-config", MountPath: "/root/.aws", ReadOnly: true}}

	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "kaniko"},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:         "download files",
					Image:        "public.ecr.aws/aws-cli/aws-cli:latest",
					Env:          env,
					VolumeMounts: volumeMounts,
					Command:      []string{"aws", "--endpoint=$(AWS_ENDPOINT)"},
					Args:         []string{"s3", "sync", "s3://$(AWS_BUCKET)", "/context"},
				},
				{
					Name:         "Generate Capnp",
					Image:        "node",
					Env:          env,
					VolumeMounts: volumeMounts,
					Command:      []string{"git", "clone", ""},
				},
			},
			Containers: nil,
		},
		Status: corev1.JobStatus{},
	}
}
