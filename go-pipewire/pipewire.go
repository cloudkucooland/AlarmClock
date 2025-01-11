package main

// main for now, will be go-pipewire

// https://gitlab.freedesktop.org/pipewire/wireplumber/-/blob/master/src/tools/wpctl.c?ref_type=heads

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
)

type wp_core_t struct {
	parent              unsafe.Pointer
	g_main_contxt       unsafe.Pointer
	properties          unsafe.Pointer
	pw_context          unsafe.Pointer
	pw_core             unsafe.Pointer
	info                unsafe.Pointer
	core_listener       unsafe.Pointer
	proxy_core_listener unsafe.Pointer
	conf                *wp_core_conf_t
	register            unsafe.Pointer
	async_tasks         unsafe.Pointer
}

type wp_core_conf_t struct {
	parent        unsafe.Pointer
	name          unsafe.Pointer
	properties    unsafe.Pointer
	conf_sections unsafe.Pointer
	files         unsafe.Pointer
}

var wp_core *wp_core_t
var wp_init func(byte)
var wp_get_library_version func() string

type g_context_t unsafe.Pointer

var g_option_context_new func(string) g_context_t

// wp_core_new(GMainContext *context, WpConf *conf, WpProperties *properties)
var wp_core_new func(g_context_t, unsafe.Pointer, unsafe.Pointer) *wp_core_t

type wp_settings_t unsafe.Pointer

var wp_settings_new func(*wp_core_t, string) wp_settings_t
var wp_core_connect func(*wp_core_t) bool
var wp_core_disconnect func(*wp_core_t)
var wp_core_get_remote_name func(*wp_core_t) string
var wp_core_get_remote_version func(*wp_core_t) string
var wp_core_load_component func(*wp_core_t, string, string, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer)

type wp_objectmanager_t unsafe.Pointer

var om wp_objectmanager_t
var wp_object_manager_new func() wp_objectmanager_t
var wp_object_manager_add_interest func(wp_objectmanager_t, int) // takes varargs

type wp_plugin_t unsafe.Pointer

var wp_plugin_find func(*wp_core_t, string) wp_plugin_t

// mixer_api, "get-volume", id, &variant
var g_signal_emit_by_name func(wp_plugin_t, string, int, *unsafe.Pointer)
var g_variant_lookup func(unsafe.Pointer, string, string, *float32)

// obviously wrong, but just a hack for testing
var g_object_new func(int, string, *wp_core_t, string, string, unsafe.Pointer) unsafe.Pointer
var g_object_set func(wp_plugin_t, string, int, unsafe.Pointer)
var wp_core_install_object_manager func(*wp_core_t, wp_objectmanager_t)

func getWirePlumberLibrary() string {
	switch runtime.GOOS {
	case "linux":
		return "libwireplumber-0.4.so.0"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func setup_libwp() {
	wp, err := purego.Dlopen(getWirePlumberLibrary(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&wp_init, wp, "wp_init")
	purego.RegisterLibFunc(&wp_get_library_version, wp, "wp_get_library_version")
	purego.RegisterLibFunc(&g_option_context_new, wp, "g_option_context_new")
	purego.RegisterLibFunc(&wp_core_new, wp, "wp_core_new")
	// purego.RegisterLibFunc(&wp_settings_new, wp, "wp_settings_new")
	purego.RegisterLibFunc(&wp_core_connect, wp, "wp_core_connect")
	purego.RegisterLibFunc(&wp_core_get_remote_name, wp, "wp_core_get_remote_name")
	purego.RegisterLibFunc(&wp_core_get_remote_version, wp, "wp_core_get_remote_version")
	purego.RegisterLibFunc(&wp_core_disconnect, wp, "wp_core_disconnect")
	purego.RegisterLibFunc(&wp_core_load_component, wp, "wp_core_load_component")
	purego.RegisterLibFunc(&wp_object_manager_new, wp, "wp_object_manager_new")
	purego.RegisterLibFunc(&wp_object_manager_add_interest, wp, "wp_object_manager_add_interest")
	purego.RegisterLibFunc(&wp_plugin_find, wp, "wp_plugin_find")
	purego.RegisterLibFunc(&g_signal_emit_by_name, wp, "g_signal_emit_by_name")
	purego.RegisterLibFunc(&g_variant_lookup, wp, "g_variant_lookup")
	purego.RegisterLibFunc(&g_object_new, wp, "g_object_new")
	purego.RegisterLibFunc(&g_object_set, wp, "g_object_set")
	purego.RegisterLibFunc(&wp_core_install_object_manager, wp, "wp_core_install_object_manager")
}

func main() {
	setup_libwp()
	wp_init(0xf) // everything for now

	fmt.Printf("WirePlumber library version: %s\n", wp_get_library_version())

	// g_context := g_option_context_new("go-wireplumber")
	g_context := g_context_t(nil)
	wp_core = wp_core_new(g_context, nil, nil)
	om = wp_object_manager_new()

	// missing from older verrsions?
	// wp_settings := wp_settings_new(wp_core, "go-pipewire")
	// nasty kludge for testing
	wp_settings := g_object_new(80, "core", wp_core, "metadata-name", "go-pipewire", nil) // guessing on the 82...
	fmt.Printf("wp_settings: %+v\n", wp_settings)

	wp_core_load_component(wp_core, "libwireplumber-module-default-nodes-api", "module", nil, nil, nil, nil, nil)
	wp_core_load_component(wp_core, "libwireplumber-module-mixer-api", "module", nil, nil, nil, nil, nil)

	// why this list? because these are the ones that didn't err out
	types := []int{8, 9, 10, 11, 80, 81, 82, 83}
	for _, i := range types {
		wp_object_manager_add_interest(om, i)
	}

	mixer_api := wp_plugin_find(wp_core, "mixer-api")
	if mixer_api == nil {
		panic("unable to get mixer")
	}

	g_object_set(mixer_api, "scale", 1, nil)
	wp_core_install_object_manager(wp_core, om)

	if connected := wp_core_connect(wp_core); !connected {
		panic("unable to connect")
	}

	// check that we are really connected
	if name := wp_core_get_remote_name(wp_core); name != "" {
		fmt.Printf("remote name: %s\n", name)
	}
	if remver := wp_core_get_remote_version(wp_core); remver != "" {
		fmt.Printf("remote version: %s\n", remver)
	}
	// fmt.Printf("wp_core.info: %+v\n", wp_core.info)

	var variant unsafe.Pointer

	// try all of them to see...
	for id := 0; id < 100; id++ {
		g_signal_emit_by_name(mixer_api, "get-volume", id, &variant)
		if variant != nil {
			fmt.Printf("%d: %+v\n", id, variant)
			var curr_volume = float32(0.0)
			g_variant_lookup(variant, "volume", "d", &curr_volume)
			fmt.Printf("%f\n", curr_volume)
		}
	}

	// set := float32(0.45)
	/*

		b := G_VARIANT_BUILDER_INIT (G_VARIANT_TYPE_VARDICT);
		g_variant_builder_add (&b, "{sv}", "volume", g_variant_new_double(set));
		variant = g_variant_builder_end (&b);

		g_signal_emit_by_name(mixer_api, "set-volume", 37, &variant)

		if variant != nil {
			fmt.Printf("%+v\n", variant)
			g_variant_lookup(variant, "volume", "d", &set);
			fmt.Printf("%f", set)
		}
	*/

	// temporary run as a daemon so I can check 'wpctl status' and see if we are connected...
	sigch := make(chan os.Signal, 3)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	fmt.Printf("shutdown requested by signal: %s\n\n", sig)

	wp_core_disconnect(wp_core)
}
