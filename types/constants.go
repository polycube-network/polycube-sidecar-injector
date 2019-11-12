package types

import (
	corev1 "k8s.io/api/core/v1"
)

var (
	SidecarStatusKey      string = "polycube.network/sidecar-status"
	SidecarStatusInjected string = "injected"
	SidecarAnnotationKey  string = "polycube.network/sidecar"
	SidecarAnnotationVal  string = "enabled"
	PolycubeContainer     corev1.Container
	PolycubeVolumes       []corev1.Volume
)

func init() {
	_true := true

	// The followings could be more easily be inserted as a JSON string,
	// but this is very clean and more readable. So, for now, let's
	// keep it like this.
	PolycubeContainer = corev1.Container{
		Name:            "polycubed",
		Image:           "polycubenetwork/polycube:latest",
		ImagePullPolicy: "Always",
		Command:         []string{"polycubed", "--loglevel=DEBUG", "--addr=0.0.0.0", "--logfile=/host/var/log/pcn_k8s"},
		VolumeMounts: []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      "lib-modules",
				MountPath: "/lib/modules",
			},
			corev1.VolumeMount{
				Name:      "usr-src",
				MountPath: "/usr/src",
			},
			corev1.VolumeMount{
				Name:      "cni-path",
				MountPath: "/host/opt/cni/bin",
			},
			corev1.VolumeMount{
				Name:      "etc-cni-netd",
				MountPath: "/host/etc/cni/net.d",
			},
			corev1.VolumeMount{
				Name:      "var-log",
				MountPath: "/host/var/log",
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Privileged: &_true,
		},
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "polycubed",
				ContainerPort: 9000,
			},
		},
		TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
	}

	PolycubeVolumes = []corev1.Volume{
		corev1.Volume{
			Name: "lib-modules",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/lib/modules",
				},
			},
		},
		corev1.Volume{
			Name: "usr-src",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/src",
				},
			},
		},
		corev1.Volume{
			Name: "cni-path",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/opt/cni/bin",
				},
			},
		},
		corev1.Volume{
			Name: "etc-cni-netd",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/etc/cni/net.d",
				},
			},
		},
		corev1.Volume{
			Name: "var-log",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log",
				},
			},
		},
		corev1.Volume{
			Name: "netns",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/run/netns",
				},
			},
		},
		corev1.Volume{
			Name: "proc",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/proc/",
				},
			},
		},
	}
}
