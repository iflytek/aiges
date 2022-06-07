package profiler

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pyroscope-io/dotnetdiag/nettrace"
)

type symbols struct {
	// Instruction pointer -> formatted string that includes
	// module name, namespace, method name and signature.
	resolved map[uint64]string
	// Slice of method addresses for sort and search.
	methodAddresses []uint64
	sorted          bool
	// MethodStartAddress -> method.
	methods map[uint64]*method
	// ModuleID -> module.
	modules map[int64]*module
}

type addresses []uint64

func (x addresses) Len() int { return len(x) }

func (x addresses) Less(i, j int) bool { return x[i] < x[j] }

func (x addresses) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// method describes MethodLoadUnloadTraceData event payload for
// Microsoft-Windows-DotNETRuntimeRundown Event ID 144.
type method struct {
	MethodID           int64
	ModuleID           int64
	MethodStartAddress uint64
	MethodSize         int32
	MethodToken        int32
	MethodFlags        int32
	MethodNamespace    string
	MethodName         string
	MethodSignature    string
}

func (d method) String() string {
	p := strings.Index(d.MethodSignature, "(")
	if p < 0 {
		p = 0
	}
	return fmt.Sprintf("%s.%s%s", d.MethodNamespace, d.MethodName, d.MethodSignature[p:])
}

// module describes ModuleLoadUnloadTraceData event payload for
// Microsoft-Windows-DotNETRuntimeRundown Event ID 152.
type module struct {
	ModuleID     int64
	AssemblyID   int64
	ModuleFlags  int32
	ModuleILPath string
}

func (d module) String() string {
	return strings.TrimSuffix(filepath.Base(d.ModuleILPath), filepath.Ext(d.ModuleILPath))
}

func newSymbols() *symbols {
	return &symbols{
		resolved: make(map[uint64]string),
		methods:  make(map[uint64]*method),
		modules:  make(map[int64]*module),
	}
}

const unresolvedSymbol = "?!?"

func (s *symbols) resolve(addr uint64) string {
	if n, ok := s.resolved[addr]; ok {
		return n
	}
	if !s.sorted {
		sort.Sort(addresses(s.methodAddresses))
		s.sorted = true
	}
	methodIdx := sort.Search(len(s.methodAddresses), func(i int) bool {
		return s.methodAddresses[i] > addr
	})
	// Method address not found.
	if methodIdx-1 == len(s.methods) || methodIdx == 0 {
		return unresolvedSymbol
	}
	methodAddress := s.methodAddresses[methodIdx-1]
	met, ok := s.methods[methodAddress]
	if !ok {
		return unresolvedSymbol
	}
	// Ensure the instruction pointer is within the method address space.
	if (met.MethodStartAddress + uint64(met.MethodSize)) <= addr {
		return unresolvedSymbol
	}
	var name string
	mod, ok := s.modules[met.ModuleID]
	if !ok {
		name = fmt.Sprintf("?!%s", met)
	} else {
		name = fmt.Sprintf("%s!%s", mod, met)
	}
	s.resolved[addr] = name
	return name
}

func (s *symbols) addModule(e *nettrace.Blob) error {
	m, err := parseModule(e.Payload)
	if err != nil {
		return err
	}
	s.modules[m.ModuleID] = &m
	return nil
}

func (s *symbols) addMethod(e *nettrace.Blob) error {
	m, err := parseMethod(e.Payload)
	if err != nil {
		return err
	}
	s.methods[m.MethodStartAddress] = &m
	s.methodAddresses = append(s.methodAddresses, m.MethodStartAddress)
	s.sorted = false
	return nil
}

func parseMethod(b *bytes.Buffer) (method, error) {
	p := &nettrace.Parser{Buffer: b}
	var d method
	p.Read(&d.MethodID)
	p.Read(&d.ModuleID)
	p.Read(&d.MethodStartAddress)
	p.Read(&d.MethodSize)
	p.Read(&d.MethodToken)
	p.Read(&d.MethodFlags)
	d.MethodNamespace = p.UTF16NTS()
	d.MethodName = p.UTF16NTS()
	d.MethodSignature = p.UTF16NTS()
	return d, p.Err()
}

func parseModule(b *bytes.Buffer) (module, error) {
	p := &nettrace.Parser{Buffer: b}
	var d module
	p.Read(&d.ModuleID)
	p.Read(&d.AssemblyID)
	p.Read(&d.ModuleFlags)
	p.Skip(12)
	d.ModuleILPath = p.UTF16NTS()
	return d, p.Err()
}
