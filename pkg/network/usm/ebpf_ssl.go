// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package usm

import (
	"debug/elf"
	"fmt"
	"os"
	"regexp"
	"strings"

	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"

	ddebpf "github.com/DataDog/datadog-agent/pkg/ebpf"
	"github.com/DataDog/datadog-agent/pkg/network/config"
	"github.com/DataDog/datadog-agent/pkg/network/ebpf/probes"
	"github.com/DataDog/datadog-agent/pkg/network/go/bininspect"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/http"
	errtelemetry "github.com/DataDog/datadog-agent/pkg/network/telemetry"
	"github.com/DataDog/datadog-agent/pkg/util/common"
	"github.com/DataDog/datadog-agent/pkg/util/kernel"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

var openSSLProbes = []manager.ProbesSelector{
	&manager.BestEffort{
		Selectors: []manager.ProbesSelector{
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_read_ex",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_read_ex",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_write_ex",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_write_ex",
				},
			},
		},
	},
	&manager.AllOf{
		Selectors: []manager.ProbesSelector{
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_do_handshake",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_do_handshake",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_connect",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_connect",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_set_bio",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_set_fd",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_read",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_read",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_write",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__SSL_write",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__SSL_shutdown",
				},
			},
		},
	},
}

var cryptoProbes = []manager.ProbesSelector{
	&manager.AllOf{
		Selectors: []manager.ProbesSelector{
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__BIO_new_socket",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__BIO_new_socket",
				},
			},
		},
	},
}

var gnuTLSProbes = []manager.ProbesSelector{
	&manager.AllOf{
		Selectors: []manager.ProbesSelector{
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_handshake",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__gnutls_handshake",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_transport_set_int2",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_transport_set_ptr",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_transport_set_ptr2",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_record_recv",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__gnutls_record_recv",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_record_send",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uretprobe__gnutls_record_send",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_bye",
				},
			},
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "uprobe__gnutls_deinit",
				},
			},
		},
	},
}

const (
	sslSockByCtxMap        = "ssl_sock_by_ctx"
	sharedLibrariesPerfMap = "shared_libraries"

	// probe used for streaming shared library events
	openatSysCall  = "openat"
	openat2SysCall = "openat2"
)

var (
	traceTypes = []string{"enter", "exit"}
)

type sslProgram struct {
	cfg                     *config.Config
	sockFDMap               *ebpf.Map
	perfHandler             *ddebpf.PerfHandler
	perfMap                 *manager.PerfMap
	watcher                 *soWatcher
	manager                 *errtelemetry.Manager
	sysOpenHooksIdentifiers []manager.ProbeIdentificationPair
}

var _ subprogram = &sslProgram{}

func newSSLProgram(c *config.Config, sockFDMap *ebpf.Map) *sslProgram {
	if !c.EnableHTTPSMonitoring || !http.HTTPSSupported(c) {
		return nil
	}

	return &sslProgram{
		cfg:                     c,
		sockFDMap:               sockFDMap,
		perfHandler:             ddebpf.NewPerfHandler(100),
		sysOpenHooksIdentifiers: getSysOpenHooksIdentifiers(),
	}
}

func (o *sslProgram) Name() string {
	return "openssl"
}

func (o *sslProgram) IsBuildModeSupported(_ buildMode) bool {
	return true
}

func (o *sslProgram) ConfigureManager(m *errtelemetry.Manager) {
	o.manager = m

	o.perfMap = &manager.PerfMap{
		Map: manager.Map{Name: sharedLibrariesPerfMap},
		PerfMapOptions: manager.PerfMapOptions{
			PerfRingBufferSize: 8 * os.Getpagesize(),
			Watermark:          1,
			RecordHandler:      o.perfHandler.RecordHandler,
			LostHandler:        o.perfHandler.LostHandler,
			RecordGetter:       o.perfHandler.RecordGetter,
		},
	}

	m.PerfMaps = append(m.PerfMaps, o.perfMap)

	for _, identifier := range o.sysOpenHooksIdentifiers {
		m.Probes = append(m.Probes,
			&manager.Probe{
				ProbeIdentificationPair: identifier,
				KProbeMaxActive:         maxActive,
			},
		)
	}
}

func (o *sslProgram) ConfigureOptions(options *manager.Options) {
	options.MapSpecEditors[sslSockByCtxMap] = manager.MapSpecEditor{
		Type:       ebpf.Hash,
		MaxEntries: o.cfg.MaxTrackedConnections,
		EditorFlag: manager.EditMaxEntries,
	}

	for _, identifier := range o.sysOpenHooksIdentifiers {
		options.ActivatedProbes = append(options.ActivatedProbes,
			&manager.ProbeSelector{
				ProbeIdentificationPair: identifier,
			},
		)
	}

	if options.MapEditors == nil {
		options.MapEditors = make(map[string]*ebpf.Map)
	}

	options.MapEditors[probes.SockByPidFDMap] = o.sockFDMap
}

func (o *sslProgram) Start() {
	// Setup shared library watcher and configure the appropriate callbacks
	o.watcher = newSOWatcher(o.perfHandler,
		soRule{
			re:           regexp.MustCompile(`libssl.so`),
			registerCB:   addHooks(o.manager, openSSLProbes),
			unregisterCB: removeHooks(o.manager, openSSLProbes),
		},
		soRule{
			re:           regexp.MustCompile(`libcrypto.so`),
			registerCB:   addHooks(o.manager, cryptoProbes),
			unregisterCB: removeHooks(o.manager, cryptoProbes),
		},
		soRule{
			re:           regexp.MustCompile(`libgnutls.so`),
			registerCB:   addHooks(o.manager, gnuTLSProbes),
			unregisterCB: removeHooks(o.manager, gnuTLSProbes),
		},
	)

	o.watcher.Start()
}

func (o *sslProgram) Stop() {
	// Detaching the sys-open hooks, as they are feeding the perf map we're going to close next.
	for _, identifier := range o.sysOpenHooksIdentifiers {
		probe, found := o.manager.GetProbe(identifier)
		if !found {
			continue
		}
		if err := probe.Stop(); err != nil {
			log.Errorf("Failed to stop hook %q. Error: %s", identifier.EBPFFuncName, err)
		}
	}

	if o.perfMap != nil {
		if err := o.perfMap.Stop(manager.CleanAll); err != nil {
			log.Errorf("Failed to stop perf map. Error: %s", err)
		}
	}

	// We must stop the watcher first, as we can read from the perfHandler, before terminating the perfHandler, otherwise
	// we might try to send events over the perfHandler.
	o.watcher.Stop()
	o.perfHandler.Stop()
}

func addHooks(m *errtelemetry.Manager, probes []manager.ProbesSelector) func(pathIdentifier, string, string) error {
	return func(id pathIdentifier, root string, path string) error {
		uid := getUID(id)

		elfFile, err := elf.Open(root + path)
		if err != nil {
			return err
		}
		defer elfFile.Close()

		symbolsSet := make(common.StringSet, 0)
		symbolsSetBestEffort := make(common.StringSet, 0)
		for _, singleProbe := range probes {
			_, isBestEffort := singleProbe.(*manager.BestEffort)
			for _, selector := range singleProbe.GetProbesIdentificationPairList() {
				_, symbol, ok := strings.Cut(selector.EBPFFuncName, "__")
				if !ok {
					continue
				}
				if isBestEffort {
					symbolsSetBestEffort[symbol] = struct{}{}
				} else {
					symbolsSet[symbol] = struct{}{}
				}
			}
		}
		symbolMap, err := bininspect.GetAllSymbolsByName(elfFile, symbolsSet)
		if err != nil {
			return err
		}
		/* Best effort to resolve symbols, so we don't care about the error */
		symbolMapBestEffort, _ := bininspect.GetAllSymbolsByName(elfFile, symbolsSetBestEffort)

		for _, singleProbe := range probes {
			_, isBestEffort := singleProbe.(*manager.BestEffort)
			for _, selector := range singleProbe.GetProbesIdentificationPairList() {
				identifier := manager.ProbeIdentificationPair{
					EBPFFuncName: selector.EBPFFuncName,
					UID:          uid,
				}
				singleProbe.EditProbeIdentificationPair(selector, identifier)
				probe, found := m.GetProbe(identifier)
				if found {
					if !probe.IsRunning() {
						err := probe.Attach()
						if err != nil {
							return err
						}
					}

					continue
				}

				_, symbol, ok := strings.Cut(selector.EBPFFuncName, "__")
				if !ok {
					continue
				}

				sym := symbolMap[symbol]
				if isBestEffort {
					sym, found = symbolMapBestEffort[symbol]
					if !found {
						continue
					}
				}
				manager.SanitizeUprobeAddresses(elfFile, []elf.Symbol{sym})
				offset, err := bininspect.SymbolToOffset(elfFile, sym)
				if err != nil {
					return err
				}

				newProbe := &manager.Probe{
					ProbeIdentificationPair: identifier,
					BinaryPath:              root + path,
					UprobeOffset:            uint64(offset),
					HookFuncName:            symbol,
				}
				_ = m.AddHook("", newProbe)
			}
			if err := singleProbe.RunValidator(m.Manager); err != nil {
				return err
			}
		}

		return nil
	}
}

func removeHooks(m *errtelemetry.Manager, probes []manager.ProbesSelector) func(pathIdentifier) error {
	return func(lib pathIdentifier) error {
		uid := getUID(lib)
		for _, singleProbe := range probes {
			for _, selector := range singleProbe.GetProbesIdentificationPairList() {
				identifier := manager.ProbeIdentificationPair{
					EBPFFuncName: selector.EBPFFuncName,
					UID:          uid,
				}
				probe, found := m.GetProbe(identifier)
				if !found {
					continue
				}

				program := probe.Program()
				err := m.DetachHook(identifier)
				if err != nil {
					log.Debugf("detach hook %s/%s : %s", selector.EBPFFuncName, uid, err)
				}
				if program != nil {
					program.Close()
				}
			}
		}

		return nil
	}
}

// getUID() return a key of length 5 as the kernel uprobe registration path is limited to a length of 64
// ebpf-manager/utils.go:GenerateEventName() MaxEventNameLen = 64
// MAX_EVENT_NAME_LEN (linux/kernel/trace/trace.h)
//
// Length 5 is arbitrary value as the full string of the eventName format is
//
//	fmt.Sprintf("%s_%.*s_%s_%s", probeType, maxFuncNameLen, functionName, UID, attachPIDstr)
//
// functionName is variable but with a minimum guarantee of 10 chars
func getUID(lib pathIdentifier) string {
	return lib.Key()[:5]
}

func (*sslProgram) GetAllUndefinedProbes() []manager.ProbeIdentificationPair {
	var probeList []manager.ProbeIdentificationPair

	for _, sslProbeList := range [][]manager.ProbesSelector{openSSLProbes, cryptoProbes, gnuTLSProbes} {
		for _, singleProbe := range sslProbeList {
			for _, identifier := range singleProbe.GetProbesIdentificationPairList() {
				probeList = append(probeList, manager.ProbeIdentificationPair{
					EBPFFuncName: identifier.EBPFFuncName,
				})
			}
		}
	}

	for _, hook := range []string{openatSysCall, openat2SysCall} {
		for _, traceType := range traceTypes {
			probeList = append(probeList, manager.ProbeIdentificationPair{
				EBPFFuncName: fmt.Sprintf("tracepoint__syscalls__sys_%s_%s", traceType, hook),
			})
		}
	}

	return probeList
}

func sysOpenAt2Supported() bool {
	missing, err := ddebpf.VerifyKernelFuncs("do_sys_openat2")
	if err == nil && len(missing) == 0 {
		return true
	}
	kversion, err := kernel.HostVersion()
	if err != nil {
		log.Error("could not determine the current kernel version. fallback to do_sys_open")
		return false
	}

	return kversion >= kernel.VersionCode(5, 6, 0)
}

// getSysOpenHooksIdentifiers returns the enter and exit tracepoints for openat and openat2 (if supported).
func getSysOpenHooksIdentifiers() []manager.ProbeIdentificationPair {
	openatProbes := []string{openatSysCall}
	if sysOpenAt2Supported() {
		openatProbes = append(openatProbes, openat2SysCall)
	}

	res := make([]manager.ProbeIdentificationPair, 0, len(traceTypes)*len(openatProbes))
	for _, probe := range openatProbes {
		for _, traceType := range traceTypes {
			res = append(res, manager.ProbeIdentificationPair{
				EBPFFuncName: fmt.Sprintf("tracepoint__syscalls__sys_%s_%s", traceType, probe),
				UID:          probeUID,
			})
		}
	}

	return res
}
