package types

var (
	// SidecarStatusKey is the key of the annotation injected
	SidecarStatusKey string = "polycube.network/sidecar-status"
	// SidecarStatusInjected is the value of the annotation
	SidecarStatusInjected string = "injected"
	// SidecarAnnotationKey is the key of the annotation that pods must have
	SidecarAnnotationKey string = "polycube.network/sidecar"
	// SidecarAnnotationVal is the value of annotation key above
	SidecarAnnotationVal string = "enabled"
	// PolycubeImageLatest is the docker image of polycube (:latest)
	PolycubeImageLatest string = "polycubenetwork/polycube:latest"
)
