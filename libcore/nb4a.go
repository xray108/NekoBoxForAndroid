package libcore

import (
	"libcore/device"
	"os"
	"path/filepath"
	"runtime"
	_ "unsafe"
	"log"

	"github.com/matsuridayo/libneko/neko_common"
	"github.com/matsuridayo/libneko/neko_log"
	"github.com/matsuridayo/sing-box-extra/boxmain"
	"github.com/sagernet/sing-box/nekoutils"
)

//go:linkname resourcePaths github.com/sagernet/sing-box/constant.resourcePaths
var resourcePaths []string

func NekoLogPrintln(s string) {
	log.Println(s)
}

func NekoLogClear() {
	neko_log.LogWriter.Truncate()
}

func ForceGc() {
	go runtime.GC()
}

func SetLocalResolver(lr LocalResolver) {
	localResolver = lr
}

func InitCore(process, cachePath, internalAssets, externalAssets string,
	maxLogSizeKb int32, logEnable bool,
	iif NB4AInterface,
) {
	defer device.DeferPanicToError("InitCore", func(err error) { log.Println(err) })

	neko_common.RunMode = neko_common.RunMode_NekoBoxForAndroid
	intfNB4A = iif

	// Working dir
	tmp := filepath.Join(cachePath, "../no_backup")
	os.MkdirAll(tmp, 0755)
	os.Chdir(tmp)

	// sing-box fs
	resourcePaths = append(resourcePaths, externalAssets)

	// Set up log
	if maxLogSizeKb < 50 {
		maxLogSizeKb = 50
	}
	neko_log.LogWriterDisable = !logEnable
	// neko_log.NB4AGuiLogWriter = iif.(io.Writer)
	neko_log.SetupLog(int(maxLogSizeKb)*1024, filepath.Join(cachePath, "neko.log"))
	boxmain.DisableColor()

	// nekoutils
	nekoutils.Selector_OnProxySelected = iif.Selector_OnProxySelected

}
