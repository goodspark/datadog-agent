// Code generated by go generate; DO NOT EDIT.

package lookup

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/network/go/bininspect"
	"github.com/go-delve/delve/pkg/goversion"
)

var SupportedArchitectures = []string{"amd64", "arm64"}

var MinGoVersion = goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}

// GetWriteParams gets the parameter metadata (positions/types) for crypto/tls.(*Conn).Write
func GetWriteParams(version goversion.GoVersion, goarch string) ([]bininspect.ParameterMetadata, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 17, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: true, StackOffset: 0, Register: 3}, {Size: 8, InReg: true, StackOffset: 0, Register: 2}, {Size: 8, InReg: true, StackOffset: 0, Register: 5}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 8, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}, {Size: 8, InReg: false, StackOffset: 24, Register: 0}, {Size: 8, InReg: false, StackOffset: 32, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 18, Rev: -1}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: true, StackOffset: 0, Register: 1}, {Size: 8, InReg: true, StackOffset: 0, Register: 2}, {Size: 8, InReg: true, StackOffset: 0, Register: 3}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 24, Register: 0}, {Size: 8, InReg: false, StackOffset: 32, Register: 0}, {Size: 8, InReg: false, StackOffset: 40, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return nil, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetReadParams gets the parameter metadata (positions/types) for crypto/tls.(*Conn).Read
func GetReadParams(version goversion.GoVersion, goarch string) ([]bininspect.ParameterMetadata, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 18, Rev: -1}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: true, StackOffset: 0, Register: 3}, {Size: 8, InReg: true, StackOffset: 0, Register: 2}, {Size: 8, InReg: true, StackOffset: 0, Register: 5}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 17, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: true, StackOffset: 0, Register: 3}, {Size: 8, InReg: false, StackOffset: 24, Register: 0}, {Size: 8, InReg: true, StackOffset: 0, Register: 5}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 16, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 8, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 8, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}, {Size: 8, InReg: false, StackOffset: 24, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 18, Rev: -1}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: true, StackOffset: 0, Register: 1}, {Size: 8, InReg: true, StackOffset: 0, Register: 2}, {Size: 8, InReg: true, StackOffset: 0, Register: 3}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 17, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 24, Register: 0}, {Size: 8, InReg: false, StackOffset: 32, Register: 0}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 16, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}}}, {TotalSize: 24, Kind: 0x17, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 24, Register: 0}, {Size: 8, InReg: false, StackOffset: 32, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return nil, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetWriteParams gets the parameter metadata (positions/types) for crypto/tls.(*Conn).Close
func GetCloseParams(version goversion.GoVersion, goarch string) ([]bininspect.ParameterMetadata, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 17, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 8, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 18, Rev: -1}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 0, InReg: true, StackOffset: 0, Register: 0}}}}, nil
		}
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return []bininspect.ParameterMetadata{{TotalSize: 8, Kind: 0x16, Pieces: []bininspect.ParameterPiece{{Size: 8, InReg: false, StackOffset: 16, Register: 0}}}}, nil
		}
		return nil, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return nil, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetTLSConnInnerConnOffset gets the offset of the "conn" field in the "crypto/tls.Conn" struct
func GetTLSConnInnerConnOffset(version goversion.GoVersion, goarch string) (uint64, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return 0, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetTCPConnInnerConnOffset gets the offset of the "conn" field in the "net.TCPConn" struct
func GetTCPConnInnerConnOffset(version goversion.GoVersion, goarch string) (uint64, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return 0, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetConnFDOffset gets the offset of the "fd" field in the "net.conn" struct
func GetConnFDOffset(version goversion.GoVersion, goarch string) (uint64, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return 0, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetNetFD_PFDOffset gets the offset of the "pfd" field in the "net.netFD" struct
func GetNetFD_PFDOffset(version goversion.GoVersion, goarch string) (uint64, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x0, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return 0, fmt.Errorf("unsupported architecture %q", goarch)
	}
}

// GetFD_SysfdOffset gets the offset of the "Sysfd" field in the "internal/poll.FD" struct
func GetFD_SysfdOffset(version goversion.GoVersion, goarch string) (uint64, error) {
	switch goarch {
	case "amd64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x10, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	case "arm64":
		if version.AfterOrEqual(goversion.GoVersion{Major: 1, Minor: 13, Rev: 0}) {
			return 0x10, nil
		}
		return 0, fmt.Errorf("unsupported version go%d.%d.%d (min supported: go%d.%d.%d)", version.Major, version.Minor, version.Rev, 1, 13, 0)
	default:
		return 0, fmt.Errorf("unsupported architecture %q", goarch)
	}
}
