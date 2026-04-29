package specrt

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/vaastav/raglan/iridescent_rt/pass"
)

type SpecializationRuntime struct {
	Filename       string
	OriginalModule *plugin.Plugin
	Pts            []*CompileTimeSpecPoint[any]
	PtsMap         map[string]*CompileTimeSpecPoint[any]
	GlobalFns      map[string]bool
	PluginFile     string
	OrigPluginFile string
	Original       string
	Trampoline     string
	Counter        int
	CallbackFns    []func() error
	Passes         []pass.SpecPass
}

func NewSpecializationRuntime(ctx context.Context, filename string) (*SpecializationRuntime, error) {
	spec_rt := &SpecializationRuntime{}
	spec_rt.Filename = filename
	spec_rt.PtsMap = make(map[string]*CompileTimeSpecPoint[any])
	// Parse file to find specialization points!
	points, err := parseOriginalModule(filename)
	if err != nil {
		return nil, err
	}
	global_fns := make(map[string]bool)
	for _, pt := range points {
		global_fns[pt.ParentFn] = true
	}
	spec_rt.GlobalFns = global_fns
	log.Println("Found the following specialization points")
	for _, pt := range points {
		log.Println(pt.String())
		spec_rt.PtsMap[pt.Name] = pt
	}
	spec_rt.Pts = points
	outf, err := setupOriginalModule(filename, global_fns)
	if err != nil {
		log.Println("Failed to setup original module")
		return nil, err
	}
	spec_rt.Original = outf
	trampoline, err := setupTrampolineModule(filename, global_fns)
	if err != nil {
		log.Println("Failed to setup trampoline module")
		return nil, err
	}
	spec_rt.Trampoline = trampoline
	specialized, err := spec_rt.setupSpecializedModule(filename)
	if err != nil {
		log.Println("Failed to setup specialized module")
		return nil, err
	}
	plugin_file := filepath.Dir(outf) + "/module.so"
	spec_rt.PluginFile = plugin_file
	spec_rt.OrigPluginFile = plugin_file
	err = spec_rt.buildModule(plugin_file, outf, trampoline, specialized)
	if err != nil {
		log.Println("Failed to build module with error: ", err)
		return nil, err
	}
	log.Println("Loading plugin:", plugin_file)
	p, err := plugin.Open(plugin_file)
	if err != nil {
		log.Println("Failed to open the built plugin")
		return nil, err
	}
	spec_rt.OriginalModule = p
	return spec_rt, nil
}

func (srt *SpecializationRuntime) buildModule(plugin_file string, orig_filename string, trampoline_filename string, spec_filename string) error {
	version_filename := filepath.Dir(plugin_file) + fmt.Sprintf("version%d.go", srt.Counter)
	f, err := os.Create(version_filename)
	defer f.Close()
	f.WriteString(fmt.Sprintf("package main\n\nvar version string=\"v%d\"\n", srt.Counter))
	cmd := exec.Command("go", "build", "-o", plugin_file, "-buildmode=plugin", orig_filename, trampoline_filename, spec_filename)
	cmd.Dir = filepath.Dir(plugin_file)
	out, err := cmd.CombinedOutput()
	res := string(out)
	if res != "" {
		log.Println(res)
	}
	return err
}

func (srt *SpecializationRuntime) UpdatePlugin() error {
	specialized, err := srt.setupSpecializedModule(srt.Filename)
	if err != nil {
		return err
	}
	//Clean up previous file
	err = os.Remove(srt.PluginFile)
	if err != nil {
		return err
	}
	srt.Counter += 1
	plugin_file := strings.ReplaceAll(srt.OrigPluginFile, "module.so", fmt.Sprintf("module_%d.so", srt.Counter))
	srt.PluginFile = plugin_file
	err = srt.buildModule(plugin_file, srt.Original, srt.Trampoline, specialized)
	if err != nil {
		return err
	}
	log.Println("Loading plugin:", plugin_file)
	p, err := plugin.Open(plugin_file)
	if err != nil {
		return err
	}
	srt.OriginalModule = p
	for _, fn := range srt.CallbackFns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (srt *SpecializationRuntime) AddCallbackFn(fn func() error) {
	srt.CallbackFns = append(srt.CallbackFns, fn)
}

func (srt *SpecializationRuntime) Lookup(fn_name string) (plugin.Symbol, error) {
	orig_sym, err := srt.OriginalModule.Lookup(fn_name + "_Trampoline")
	if err != nil {
		return nil, err
	}
	return orig_sym, nil
}

func (srt *SpecializationRuntime) Specialize(name string, index int) {
	pt := srt.PtsMap[name]
	pt.Specialize(index)
}

func (srt *SpecializationRuntime) Instrument(name string, key int) {
	if pt, ok := srt.PtsMap[name]; ok {
		pt.Incr(key)
	}
}

func (srt *SpecializationRuntime) AddSpecializationPass(p pass.SpecPass) {
	srt.Passes = append(srt.Passes, p)
}
