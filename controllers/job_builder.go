package controllers

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateAwsConfig() []v1.EnvVar {
	return []v1.EnvVar{
		{Name: "AWS_PROFILE", Value: "default"},
		{Name: "AWS_ENDPOINT", Value: "https://s3.fr-par.scw.cloud"},
		{Name: "AWS_BUCKET", Value: "stage-cf-worker"},
	}
}

func generateDownloadFilesContainer() v1.Container {
	return v1.Container{
		Name:            "download-files",
		Image:           "public.ecr.aws/aws-cli/aws-cli:latest",
		ImagePullPolicy: "IfNotPresent",
		Env:             generateAwsConfig(),
		VolumeMounts: []v1.VolumeMount{
			{Name: "s3-config", MountPath: "/root/.aws", ReadOnly: true},
			{Name: "context", MountPath: "/context"},
		},
		Command: []string{"aws", "--endpoint-url=$(AWS_ENDPOINT)"},
		Args:    []string{"s3", "sync", "s3://$(AWS_BUCKET)", "/context"},
	}
}

func generateGettingDockerfile() v1.Container {
	return v1.Container{
		Name:            "getting-dockerfile",
		Image:           "curlimages/curl",
		ImagePullPolicy: "IfNotPresent",
		VolumeMounts: []v1.VolumeMount{
			{Name: "context", MountPath: "/context", ReadOnly: false},
		},
		Command: []string{"curl"},
		Args:    []string{"-o", "/context/Dockerfile", "-L", "https://raw.githubusercontent.com/clementreiffers/JobBuilder/main/build-worker.Dockerfile"},
	}
}

func createJob() batchv1.Job {
	ttl := int32(3600)
	parallelism := int32(1)
	completions := int32(1)
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-go",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Parallelism:             &parallelism,
			Completions:             &completions,
			TTLSecondsAfterFinished: &ttl,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						generateDownloadFilesContainer(),
						generateGettingDockerfile(),
					},
				},
			},
		},
		Status: batchv1.JobStatus{},
	}
}
