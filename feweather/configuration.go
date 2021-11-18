package fe-measure

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

type deviceURLs []string

func (u *deviceURLs) String() string { return "" }
func (u *deviceURLs) Set(value string) error {
	*u = append(*u, value)
	return nil
}

type DeviceUrlConfiguration struct {
	DeviceURLs []string `json:"deviceurl"`
}

type FrontendConfiguration struct {
	StaticDir    string `json:"staticdir"`
	StaticURL    string `json:"staticurl"`
	TemplateFile string `json:"templatefile"`
	TemplateURL  string `json:"tempalteurl"`
}

type IFrontendConfiguration interface {
	GetStaticDir() string
	GetStaticURL() string
	GetTemplateFile() string
	GetTemplateURL() string
}

// Configuration ...
type Configuration struct {
	microservice.Configuration
	FrontendConfiguration
	DeviceUrlConfiguration
}

// IConfiguration ...
type IConfiguration interface {
	microservice.IConfiguration
	IFrontendConfiguration
}

// GetUpstream ...
func (cfg FrontendConfiguration) GetStaticDir() string    { return cfg.StaticDir }
func (cfg FrontendConfiguration) GetStaticURL() string    { return cfg.StaticURL }
func (cfg FrontendConfiguration) GetTemplateFile() string { return cfg.TemplateFile }
func (cfg FrontendConfiguration) GetTemplateURL() string  { return cfg.TemplateURL }

// ---------------------------------------------------------------------------

// InitConfigurationFromArgs ...
func InitConfigurationFromArgs(cfg *Configuration, args []string, flagset *flag.FlagSet) {
	var du deviceURLs

	if flagset == nil {
		flagset = flag.NewFlagSet("fe-measure", flag.PanicOnError)
	}

	flagset.Var(&du, "device", "Device URls.")
	pstaticdir := flagset.String("staticdir", "", "Path for static assets.")
	pstaticurl := flagset.String("staticurl", "", "The URL to the static assests in the frontend.")
	ptemplatefile := flagset.String("templatefile", "", "Path to the template file.")
	ptempalteurl := flagset.String("tempalteurl", "", "The URL of the template in the frontend.")

	microservice.InitConfigurationFromArgs(&cfg.Configuration, args, flagset)

	if len(*pstaticdir) > 0 {
		cfg.StaticDir = *pstaticdir
	}

	if len(*pstaticurl) > 0 {
		cfg.StaticURL = *pstaticurl
	}

	if len(*ptemplatefile) > 0 {
		cfg.TemplateFile = *ptemplatefile
	}

	if len(*ptempalteurl) > 0 {
		cfg.TemplateURL = *ptempalteurl
	}

	if len(du) > 0 {
		cfg.DeviceURLs = du
	}

	if len(cfg.GetConfigurationFile()) > 0 {
		file, err := os.Open(cfg.GetConfigurationFile())

		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf("Error: Failed to open onfiguration file '%s'. Error was %s!\n", cfg.GetConfigurationFile(), err.Error()))
		}

		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			flagset.Usage()
			panic(fmt.Sprintf("Error: Failed to parse onfiguration file '%s'. Error was %s!\n", cfg.GetConfigurationFile(), err.Error()))
		}

		if len(cfg.StaticDir) == 0 {
			cfg.StaticDir = configuration.StaticDir
		}

		if len(cfg.StaticURL) == 0 {
			cfg.StaticURL = configuration.StaticURL
		}

		if len(cfg.TemplateFile) == 0 {
			cfg.TemplateFile = configuration.TemplateFile
		}

		if len(cfg.TemplateURL) == 0 {
			cfg.TemplateURL = configuration.TemplateURL
		}
	}

	if len(cfg.StaticDir) == 0 {
		cfg.StaticDir = os.Getenv("FE_STATICDIR")
	}

	if len(cfg.StaticURL) == 0 {
		cfg.StaticURL = os.Getenv("FE_STATICURL")
	}

	if len(cfg.TemplateFile) == 0 {
		cfg.TemplateFile = os.Getenv("FE_TEMPLATEFILE")
	}

	if len(cfg.TemplateFile) == 0 {
		cfg.TemplateFile = os.Getenv("FE_TEMPLATEURL")
	}

	if len(cfg.StaticURL) == 0 {
		cfg.StaticURL = "/static/"
	}

	if len(cfg.TemplateFile) == 0 {
		cfg.TemplateFile = "index.html"
	}

	if len(cfg.TemplateURL) == 0 {
		cfg.TemplateURL = "/index.html"
	}

	cfg.TemplateFile = path.Clean(cfg.TemplateFile)

	if len(cfg.StaticDir) > 0 {
		cfg.StaticDir = path.Clean(cfg.StaticDir)
	}

	if len(cfg.StaticURL) > 0 {
		cfg.StaticURL = path.Clean(cfg.StaticURL) + "/"
	}
}
