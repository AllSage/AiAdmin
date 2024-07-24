package envconfig

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type AiAdminHost struct {
	Scheme string
	Host   string
	Port   string
}

func (o AiAdminHost) String() string {
	return fmt.Sprintf("%s://%s:%s", o.Scheme, o.Host, o.Port)
}

var ErrInvalidHostPort = errors.New("invalid port specified in AiAdmin_HOST")

var (
	// Set via AiAdmin_ORIGINS in the environment
	AllowOrigins []string
	// Set via AiAdmin_DEBUG in the environment
	Debug bool
	// Experimental flash attention
	FlashAttention bool
	// Set via AiAdmin_HOST in the environment
	Host *AiAdminHost
	// Set via AiAdmin_KEEP_ALIVE in the environment
	KeepAlive time.Duration
	// Set via AiAdmin_LLM_LIBRARY in the environment
	LLMLibrary string
	// Set via AiAdmin_MAX_LOADED_MODELS in the environment
	MaxRunners int
	// Set via AiAdmin_MAX_QUEUE in the environment
	MaxQueuedRequests int
	// Set via AiAdmin_MODELS in the environment
	ModelsDir string
	// Set via AiAdmin_NOHISTORY in the environment
	NoHistory bool
	// Set via AiAdmin_NOPRUNE in the environment
	NoPrune bool
	// Set via AiAdmin_NUM_PARALLEL in the environment
	NumParallel int
	// Set via AiAdmin_RUNNERS_DIR in the environment
	RunnersDir string
	// Set via AiAdmin_SCHED_SPREAD in the environment
	SchedSpread bool
	// Set via AiAdmin_TMPDIR in the environment
	TmpDir string
	// Set via AiAdmin_INTEL_GPU in the environment
	IntelGpu bool

	// Set via CUDA_VISIBLE_DEVICES in the environment
	CudaVisibleDevices string
	// Set via HIP_VISIBLE_DEVICES in the environment
	HipVisibleDevices string
	// Set via ROCR_VISIBLE_DEVICES in the environment
	RocrVisibleDevices string
	// Set via GPU_DEVICE_ORDINAL in the environment
	GpuDeviceOrdinal string
	// Set via HSA_OVERRIDE_GFX_VERSION in the environment
	HsaOverrideGfxVersion string
)

type EnvVar struct {
	Name        string
	Value       any
	Description string
}

func AsMap() map[string]EnvVar {
	ret := map[string]EnvVar{
		"AiAdmin_DEBUG":             {"AiAdmin_DEBUG", Debug, "Show additional debug information (e.g. AiAdmin_DEBUG=1)"},
		"AiAdmin_FLASH_ATTENTION":   {"AiAdmin_FLASH_ATTENTION", FlashAttention, "Enabled flash attention"},
		"AiAdmin_HOST":              {"AiAdmin_HOST", Host, "IP Address for the AiAdmin server (default 127.0.0.1:11434)"},
		"AiAdmin_KEEP_ALIVE":        {"AiAdmin_KEEP_ALIVE", KeepAlive, "The duration that models stay loaded in memory (default \"5m\")"},
		"AiAdmin_LLM_LIBRARY":       {"AiAdmin_LLM_LIBRARY", LLMLibrary, "Set LLM library to bypass autodetection"},
		"AiAdmin_MAX_LOADED_MODELS": {"AiAdmin_MAX_LOADED_MODELS", MaxRunners, "Maximum number of loaded models per GPU"},
		"AiAdmin_MAX_QUEUE":         {"AiAdmin_MAX_QUEUE", MaxQueuedRequests, "Maximum number of queued requests"},
		"AiAdmin_MODELS":            {"AiAdmin_MODELS", ModelsDir, "The path to the models directory"},
		"AiAdmin_NOHISTORY":         {"AiAdmin_NOHISTORY", NoHistory, "Do not preserve readline history"},
		"AiAdmin_NOPRUNE":           {"AiAdmin_NOPRUNE", NoPrune, "Do not prune model blobs on startup"},
		"AiAdmin_NUM_PARALLEL":      {"AiAdmin_NUM_PARALLEL", NumParallel, "Maximum number of parallel requests"},
		"AiAdmin_ORIGINS":           {"AiAdmin_ORIGINS", AllowOrigins, "A comma separated list of allowed origins"},
		"AiAdmin_RUNNERS_DIR":       {"AiAdmin_RUNNERS_DIR", RunnersDir, "Location for runners"},
		"AiAdmin_SCHED_SPREAD":      {"AiAdmin_SCHED_SPREAD", SchedSpread, "Always schedule model across all GPUs"},
		"AiAdmin_TMPDIR":            {"AiAdmin_TMPDIR", TmpDir, "Location for temporary files"},
	}
	if runtime.GOOS != "darwin" {
		ret["CUDA_VISIBLE_DEVICES"] = EnvVar{"CUDA_VISIBLE_DEVICES", CudaVisibleDevices, "Set which NVIDIA devices are visible"}
		ret["HIP_VISIBLE_DEVICES"] = EnvVar{"HIP_VISIBLE_DEVICES", HipVisibleDevices, "Set which AMD devices are visible"}
		ret["ROCR_VISIBLE_DEVICES"] = EnvVar{"ROCR_VISIBLE_DEVICES", RocrVisibleDevices, "Set which AMD devices are visible"}
		ret["GPU_DEVICE_ORDINAL"] = EnvVar{"GPU_DEVICE_ORDINAL", GpuDeviceOrdinal, "Set which AMD devices are visible"}
		ret["HSA_OVERRIDE_GFX_VERSION"] = EnvVar{"HSA_OVERRIDE_GFX_VERSION", HsaOverrideGfxVersion, "Override the gfx used for all detected AMD GPUs"}
		ret["AiAdmin_INTEL_GPU"] = EnvVar{"AiAdmin_INTEL_GPU", IntelGpu, "Enable experimental Intel GPU detection"}
	}
	return ret
}

func Values() map[string]string {
	vals := make(map[string]string)
	for k, v := range AsMap() {
		vals[k] = fmt.Sprintf("%v", v.Value)
	}
	return vals
}

var defaultAllowOrigins = []string{
	"localhost",
	"127.0.0.1",
	"0.0.0.0",
}

// Clean quotes and spaces from the value
func clean(key string) string {
	return strings.Trim(os.Getenv(key), "\"' ")
}

func init() {
	// default values
	NumParallel = 0 // Autoselect
	MaxRunners = 0  // Autoselect
	MaxQueuedRequests = 512
	KeepAlive = 5 * time.Minute

	LoadConfig()
}

func LoadConfig() {
	if debug := clean("AiAdmin_DEBUG"); debug != "" {
		d, err := strconv.ParseBool(debug)
		if err == nil {
			Debug = d
		} else {
			Debug = true
		}
	}

	if fa := clean("AiAdmin_FLASH_ATTENTION"); fa != "" {
		d, err := strconv.ParseBool(fa)
		if err == nil {
			FlashAttention = d
		}
	}

	RunnersDir = clean("AiAdmin_RUNNERS_DIR")
	if runtime.GOOS == "windows" && RunnersDir == "" {
		// On Windows we do not carry the payloads inside the main executable
		appExe, err := os.Executable()
		if err != nil {
			slog.Error("failed to lookup executable path", "error", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			slog.Error("failed to lookup working directory", "error", err)
		}

		var paths []string
		for _, root := range []string{filepath.Dir(appExe), cwd} {
			paths = append(paths,
				root,
				filepath.Join(root, "windows-"+runtime.GOARCH),
				filepath.Join(root, "dist", "windows-"+runtime.GOARCH),
			)
		}

		// Try a few variations to improve developer experience when building from source in the local tree
		for _, p := range paths {
			candidate := filepath.Join(p, "AiAdmin_runners")
			_, err := os.Stat(candidate)
			if err == nil {
				RunnersDir = candidate
				break
			}
		}
		if RunnersDir == "" {
			slog.Error("unable to locate llm runner directory.  Set AiAdmin_RUNNERS_DIR to the location of 'AiAdmin_runners'")
		}
	}

	TmpDir = clean("AiAdmin_TMPDIR")

	LLMLibrary = clean("AiAdmin_LLM_LIBRARY")

	if onp := clean("AiAdmin_NUM_PARALLEL"); onp != "" {
		val, err := strconv.Atoi(onp)
		if err != nil {
			slog.Error("invalid setting, ignoring", "AiAdmin_NUM_PARALLEL", onp, "error", err)
		} else {
			NumParallel = val
		}
	}

	if nohistory := clean("AiAdmin_NOHISTORY"); nohistory != "" {
		NoHistory = true
	}

	if spread := clean("AiAdmin_SCHED_SPREAD"); spread != "" {
		s, err := strconv.ParseBool(spread)
		if err == nil {
			SchedSpread = s
		} else {
			SchedSpread = true
		}
	}

	if noprune := clean("AiAdmin_NOPRUNE"); noprune != "" {
		NoPrune = true
	}

	if origins := clean("AiAdmin_ORIGINS"); origins != "" {
		AllowOrigins = strings.Split(origins, ",")
	}
	for _, allowOrigin := range defaultAllowOrigins {
		AllowOrigins = append(AllowOrigins,
			fmt.Sprintf("http://%s", allowOrigin),
			fmt.Sprintf("https://%s", allowOrigin),
			fmt.Sprintf("http://%s", net.JoinHostPort(allowOrigin, "*")),
			fmt.Sprintf("https://%s", net.JoinHostPort(allowOrigin, "*")),
		)
	}

	AllowOrigins = append(AllowOrigins,
		"app://*",
		"file://*",
		"tauri://*",
	)

	maxRunners := clean("AiAdmin_MAX_LOADED_MODELS")
	if maxRunners != "" {
		m, err := strconv.Atoi(maxRunners)
		if err != nil {
			slog.Error("invalid setting, ignoring", "AiAdmin_MAX_LOADED_MODELS", maxRunners, "error", err)
		} else {
			MaxRunners = m
		}
	}

	if onp := os.Getenv("AiAdmin_MAX_QUEUE"); onp != "" {
		p, err := strconv.Atoi(onp)
		if err != nil || p <= 0 {
			slog.Error("invalid setting, ignoring", "AiAdmin_MAX_QUEUE", onp, "error", err)
		} else {
			MaxQueuedRequests = p
		}
	}

	ka := clean("AiAdmin_KEEP_ALIVE")
	if ka != "" {
		loadKeepAlive(ka)
	}

	var err error
	ModelsDir, err = getModelsDir()
	if err != nil {
		slog.Error("invalid setting", "AiAdmin_MODELS", ModelsDir, "error", err)
	}

	Host, err = getAiAdminHost()
	if err != nil {
		slog.Error("invalid setting", "AiAdmin_HOST", Host, "error", err, "using default port", Host.Port)
	}

	if set, err := strconv.ParseBool(clean("AiAdmin_INTEL_GPU")); err == nil {
		IntelGpu = set
	}

	CudaVisibleDevices = clean("CUDA_VISIBLE_DEVICES")
	HipVisibleDevices = clean("HIP_VISIBLE_DEVICES")
	RocrVisibleDevices = clean("ROCR_VISIBLE_DEVICES")
	GpuDeviceOrdinal = clean("GPU_DEVICE_ORDINAL")
	HsaOverrideGfxVersion = clean("HSA_OVERRIDE_GFX_VERSION")
}

func getModelsDir() (string, error) {
	if models, exists := os.LookupEnv("AiAdmin_MODELS"); exists {
		return models, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".AiAdmin", "models"), nil
}

func getAiAdminHost() (*AiAdminHost, error) {
	defaultPort := "11434"

	hostVar := os.Getenv("AiAdmin_HOST")
	hostVar = strings.TrimSpace(strings.Trim(strings.TrimSpace(hostVar), "\"'"))

	scheme, hostport, ok := strings.Cut(hostVar, "://")
	switch {
	case !ok:
		scheme, hostport = "http", hostVar
	case scheme == "http":
		defaultPort = "80"
	case scheme == "https":
		defaultPort = "443"
	}

	// trim trailing slashes
	hostport = strings.TrimRight(hostport, "/")

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host, port = "127.0.0.1", defaultPort
		if ip := net.ParseIP(strings.Trim(hostport, "[]")); ip != nil {
			host = ip.String()
		} else if hostport != "" {
			host = hostport
		}
	}

	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 0 {
		return &AiAdminHost{
			Scheme: scheme,
			Host:   host,
			Port:   defaultPort,
		}, ErrInvalidHostPort
	}

	return &AiAdminHost{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}, nil
}

func loadKeepAlive(ka string) {
	v, err := strconv.Atoi(ka)
	if err != nil {
		d, err := time.ParseDuration(ka)
		if err == nil {
			if d < 0 {
				KeepAlive = time.Duration(math.MaxInt64)
			} else {
				KeepAlive = d
			}
		}
	} else {
		d := time.Duration(v) * time.Second
		if d < 0 {
			KeepAlive = time.Duration(math.MaxInt64)
		} else {
			KeepAlive = d
		}
	}
}
