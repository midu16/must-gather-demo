// pkg/gather/pod.go
package gather

import (
	"fmt"
	"path"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *MustGatherOptions) NewPod(node, image string, hasMaster bool) *corev1.Pod {
	zero := int64(0)

	nodeSelector := map[string]string{
		corev1.LabelOSStable: "linux",
	}
	if node == "" && hasMaster {
		nodeSelector["node-role.kubernetes.io/master"] = ""
	}

	executedCommand := "/usr/bin/gather"
	if len(o.Command) > 0 {
		executedCommand = strings.Join(o.Command, " ")
	}

	cleanedSourceDir := path.Clean(o.SourceDir)
	volumeChecker := fmt.Sprintf(
		volumeUsageCheckerScript,
		cleanedSourceDir,
		cleanedSourceDir,
		o.VolumePercentage,
		o.VolumePercentage,
		executedCommand,
	)

	excludedTaints := []corev1.Taint{
		{Key: unreachableTaintKey, Effect: corev1.TaintEffectNoExecute},
		{Key: unreachableTaintKey, Effect: corev1.TaintEffectNoSchedule},
	}

	var tolerations []corev1.Toleration
	if node == "" && hasMaster {
		tolerations = append(tolerations, corev1.Toleration{
			Key:      "node-role.kubernetes.io/master",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		})
	}

	tolerations = append(tolerations, corev1.Toleration{
		Key:      "node.kubernetes.io/not-ready",
		Operator: corev1.TolerationOpExists,
		Effect:   corev1.TaintEffectNoSchedule,
	})

	filteredTolerations := make([]corev1.Toleration, 0)
TolerationLoop:
	for _, tol := range tolerations {
		for _, excluded := range excludedTaints {
			if tol.ToleratesTaint(&excluded) {
				continue TolerationLoop
			}
		}
		filteredTolerations = append(filteredTolerations, tol)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "must-gather-",
			Labels: map[string]string{
				"app": "must-gather",
			},
		},
		Spec: corev1.PodSpec{
			NodeName:          node,
			PriorityClassName: "system-cluster-critical",
			RestartPolicy:     corev1.RestartPolicyNever,
			Volumes: []corev1.Volume{
				{
					Name: "must-gather-output",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/log/must-gather",
							Type: new(corev1.HostPathType),
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            gatherContainerName,
					Image:           image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"/bin/bash", "-c", fmt.Sprintf("mkdir -p %s && %s & %s; sync", cleanedSourceDir, volumeChecker, executedCommand)},
					Env:             o.getEnvVars(),
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "must-gather-output",
							MountPath: cleanedSourceDir,
						},
					},
				},
				{
					Name:            "copy",
					Image:           image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "must-gather-output",
							MountPath: cleanedSourceDir,
						},
					},
				},
			},
			HostNetwork:                   o.HostNetwork,
			NodeSelector:                  nodeSelector,
			TerminationGracePeriodSeconds: &zero,
			Tolerations:                   filteredTolerations,
		},
	}

	if o.HostNetwork {
		pod.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"CAP_NET_RAW"},
			},
		}
	}

	return pod
}

// getEnvVars remains unchanged

func (o *MustGatherOptions) getEnvVars() []corev1.EnvVar {
	env := []corev1.EnvVar{
		{
			Name: "NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
	}

	if o.Since != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "MUST_GATHER_SINCE",
			Value: o.Since.String(),
		})
	}

	if o.SinceTime != "" {
		env = append(env, corev1.EnvVar{
			Name:  "MUST_GATHER_SINCE_TIME",
			Value: o.SinceTime,
		})
	}

	return env
}