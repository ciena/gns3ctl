/*
Copyright 2022 Ciena Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gns3

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	grab "github.com/cavaliergopher/grab/v3"
	units "github.com/docker/go-units"
	"github.com/spf13/viper"
	"zgo.at/termfo"
	"zgo.at/termfo/caps"
)

const (
	AppliancesPath = "v2/appliances"
)

//nolint:tagliatelle
type ApplianceImage struct {
	DirectDownloadUrl string `json:"direct_download_url,omitempty" yaml:"direct_download_url"`
	DownloadUrl       string `json:"download_url,omitempty" yaml:"download_url"`
	Filename          string `json:"filename,omitempty" yaml:"filename"`
	FileSize          int64  `json:"filesize,omitempty" yaml:"filesize"`
	Md5Sum            string `json:"md5sum,omitempty" yaml:"md5sum"`
	Version           string `json:"version,omitempty" yaml:"version"`
}

type ApplianceDocker struct {
	Image             string   `json:"image,omitempty" yaml:"image,omitempty"`
	Usage             string   `json:"usage,omitempty" yaml:"usage,omitempty"`
	Adapters          int      `json:"adapters,omitempty" yaml:"adapters,omitempty"`
	StartCommand      string   `json:"start_command,omitempty" yaml:"start_command,omitempty"`
	Environment       string   `json:"environment,omitempty" yaml:"environment,omitempty"`
	ConsoleType       string   `json:"console_type,omitempty" yaml:"console_type,omitempty"`
	ConsoleAutoStart  bool     `json:"console_auto_start,omitempty" yaml:"console_auto_start,omitempty"`
	ConsoleHttpPort   int      `json:"console_http_port,omitempty" yaml:"console_http_port,omitempty"`
	ConsoleHttpPath   string   `json:"console_http_path,omitempty" yaml:"console_http_path,omitempty"`
	ConsoleResolution string   `json:"console_resolution,omitempty" yaml:"console_resolution,omitempty"`
	ExtraHosts        string   `json:"extra_hosts,omitempty" yaml:"extra_hosts,omitempty"`
	ExtraVolumes      []string `json:"extra_volumes,omitempty" yaml:"extra_volumes,omitempty"`
	CustomAdapters    []*struct {
		Schema               string   `json:"$schema,omitempty" yaml:"$schema,omitempty"`
		Description          string   `json:"description,omitempty" yaml:"description,omitempty"`
		Type                 any      `json:"type,omitempty" yaml:"type,omitempty"`
		Properties           any      `json:"properties,omitempty" yaml:"properties,omitempty"`
		Required             []string `json:"required,omitempty" yaml:"required,omitempty"`
		AdditionalProperties bool     `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	} `json:"custom_adapters,omitempty" yaml:"custom_adapters,omitempty"`
}

type ApplianceIou struct{}
type ApplianceDynamips struct{}

//nolint:tagliatelle
type ApplianceQemu struct {
	AdapterType       string `json:"adapter_type,omitempty" yaml:"adapter_type"`
	Adapters          int    `json:"adapters,omitempty" yaml:"adapters"`
	Arch              string `json:"arch,omitempty" yaml:"arch"`
	BootPriority      string `json:"boot_priority,omitempty" yaml:"boot_priority"`
	ConsoleType       string `json:"console_type,omitempty" yaml:"console_type"`
	HdaDiskInterface  string `json:"hda_disk_interface,omitempty" yaml:"hda_disk_interface"`
	Kvm               string `json:"kvm,omitempty" yaml:"kvm"`
	KernelCommandLine string `json:"kernel_command_line,omitempty" yaml:"kernel_command_line"`
	Options           string `json:"options,omitempty" yaml:"options"`
	Ram               int64  `json:"ram,omitempty" yaml:"ram"`
}

//nolint:tagliatelle
type ApplianceVersion struct {
	Images struct {
		CdromImage   string `json:"cdrom_image,omitempty" yaml:"cdrom_image"`
		BiosImage    string `json:"bios_image,omitempty" yaml:"bios_image"`
		HdaDiskImage string `json:"hda_disk_image,omitempty" yaml:"hda_disk_image"`
	} `json:"images,omitempty" yaml:"images"`
	Name string `json:"name,omitempty" yaml:"name"`
}

//nolint:tagliatelle
type Appliance struct {
	Builtin          bool                `json:"builtin,omitempty" yaml:"builtin"`
	Category         string              `json:"category,omitempty" yaml:"category"`
	Description      string              `json:"description,omitempty" yaml:"description"`
	DocumentationUrl string              `json:"documentation_url,omitempty" yaml:"documentation_url"`
	Images           []*ApplianceImage   `json:"images,omitempty" yaml:"images"`
	Maintainer       string              `json:"maintainer,omitempty" yaml:"maintainer"`
	MaintaineEmail   string              `json:"maintainer_email,omitempty" yaml:"maintainer_email"`
	RegistryVersion  int                 `json:"registry_version,omitempty" yaml:"registry_version"`
	ProductName      string              `json:"product_name,omitempty" yaml:"product_name"`
	ProductUrl       string              `json:"product_url,omitempty" yaml:"product_url"`
	Name             string              `json:"name,omitempty" yaml:"name"`
	Symbol           string              `json:"symbol,omitempty" yaml:"symbol"`
	Status           string              `json:"status,omitempty" yaml:"status"`
	FirstPortName    string              `json:"first_port_name,omitempty" yaml:"first_port_name"`
	PortNameFormat   string              `json:"port_name_format,omitempty" yaml:"port_name_format"`
	Qemu             *ApplianceQemu      `json:"qemu,omitempty" yaml:"qemu"`
	Docker           *ApplianceDocker    `json:"docker,omitempty" yaml:"docker"`
	Iou              *ApplianceIou       `json:"iou,omitempty" yaml:"iou"`
	Dynamips         *ApplianceDynamips  `json:"dynamips,omitempty" yaml:"dynamips"`
	Usage            string              `json:"usage,omitempty" yaml:"usage"`
	VendorName       string              `json:"vendor_name,omitempty" yaml:"vendor_name"`
	VendorUrl        string              `json:"vendor_url,omitempty" yaml:"vendor_url"`
	Versions         []*ApplianceVersion `json:"versions,omitempty" yaml:"versions"`
}

type Appliances struct {
	gns3 *Gns3
}

var (
	ImageTypes = map[string]string{TemplateTypeQemu: "QEMU", TemplateTypeDynamips: "IOS", TemplateTypeIou: "IOU", TemplateTypeDocker: "DOCKER"}
)

func (g *Gns3) Appliances() *Appliances {
	return &Appliances{gns3: g}
}

func (a *Appliances) List() ([]Appliance, error) {
	list := []Appliance{}
	err := a.gns3.Get(AppliancesPath, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *Appliances) Get(id string) (*Appliance, error) {
	var appliance Appliance
	// Think this may be a UUID, so try to delete directly
	err := a.gns3.Get(fmt.Sprintf("%s/%s", AppliancesPath, id), &appliance)
	if err == nil {
		return &appliance, nil
	}

	return nil, err
}

const ImagePath = "%s/images/%s/%s"
const ImageMd5Path = "%s/images/%s/%s.md5sum"

func (a Appliances) downloadUrlToFile(url, outname, outmd5 string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(outname, url)
	if err != nil {
		return err
	}
	req.NoResume = true
	if val, err := units.FromHumanSize(viper.GetString("download-buffer-size")); err == nil {
		req.BufferSize = int(val)
	} else {
		fmt.Printf("ERROR: unable to parse download buffer size '%s', defaulting to '10M'\n",
			viper.GetString("download-buffer-size"))
		val, _ := units.FromHumanSize("10MB")
		req.BufferSize = int(val)
	}

	ti, _ := termfo.New("")
	cr := ti.Strings[caps.CarriageReturn]
	ceol := ti.Strings[caps.ClrEol]

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(time.Second)
	defer t.Stop()
	status := false

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("%s  transferred %.2f MB/%.2f MB (%.2f%%)%s",
				cr,
				float64(resp.BytesComplete()/1024/1024),
				float64(resp.Size()/1024/1024),
				100*resp.Progress(),
				ceol)
			status = true

		case <-resp.Done:
			// download is complete
			if status {
				fmt.Println()
			}
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	if outmd5 != "" {
		fmt.Printf("Generating MD5 sum for '%s'\n", outname)
		err = a.generateMd5Sum(outname, outmd5)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Appliances) generateMd5Sum(filename, md5name string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	data := fmt.Sprintf("%x", h.Sum(nil))
	err = os.WriteFile(md5name, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *Appliances) checkAndDownload(url, imgType, outname, outmd5 string) error {
	basedir := viper.GetString("base-directory")
	filename := fmt.Sprintf(ImagePath, basedir, imgType, outname)
	md5name := fmt.Sprintf(ImageMd5Path, basedir, imgType, outname)

	// If we already have a local md5 file, read it and compare to expected
	// md5 value
	var err error
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		err := a.downloadUrlToFile(url, filename, md5name)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(md5name)
	if os.IsNotExist(err) {
		fmt.Printf("Generating MD5 sum for file '%s'\n", filename)
		err = a.generateMd5Sum(filename, md5name)
		if err != nil {
			return err
		}
	}

	data, err := os.ReadFile(md5name)
	if err != nil {
		return err
	}

	if string(data) != outmd5 {
		fmt.Println("MD5 mismatch, initiating download")
		err := a.downloadUrlToFile(url, filename, md5name)
		if err != nil {
			return err
		}
	}

	data, err = os.ReadFile(md5name)
	if err != nil {
		return err
	}
	if string(data) != outmd5 {
		return ErrMd5Mismatch
	}
	return nil
}

func (a *Appliances) Import(file io.Reader, inputDirectory string) (*Appliance, *Template, error) {
	var app Appliance

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&app); err != nil {
		return nil, nil, err
	}

	templateName := app.Name
	if len(app.Versions) > 0 && app.Versions[0].Name != "" {
		templateName += " " + app.Versions[0].Name
	}

	tmpl := &Template{
		ComputeId:     "local",
		Name:          templateName,
		Usage:         app.Usage,
		FirstPortName: app.FirstPortName,
	}

	if app.Category == CategoryMultilayerSwitch {
		tmpl.Category = CategorySwitch
	} else {
		tmpl.Category = app.Category
	}
	if app.Symbol != "" {
		// If the appliance specified a symbol, then copy it to the GNS3 data
		// directory. First we will check of it is a URL of local file
		// reference.
		if u, err := url.Parse(app.Symbol); err == nil {
			dest := fmt.Sprintf("%s/symbols/%s", viper.GetString("base-directory"), path.Base(u.Path))
			switch strings.ToLower(u.Scheme) {
			case "http", "https":
				// Download and write file
				if err := a.downloadUrlToFile(app.Symbol, dest, ""); err != nil {
					return nil, nil, err
				}
			case "":
				// Local file copy
				// if the app.Symbol reference does not begin with a `/`, then assume
				// it is relative to the input file
				src := app.Symbol
				if !strings.HasPrefix("/", src) {
					src = fmt.Sprintf("%s/%s", inputDirectory, src)
				}

				input, err := os.ReadFile(src)
				if err != nil {
					return nil, nil, err
				}

				err = os.WriteFile(dest, input, 0644)
				if err != nil {
					return nil, nil, err
				}
			default:
				fmt.Printf("WARNING: unsupported file schema type, '%s', ignoring\n", u.Scheme)
			}
		} else {
			fmt.Printf("WARNING: unable to parse symbol referece, '%s', ignoring\n", app.Symbol)
		}
		tmpl.Symbol = app.Symbol
	} else {
		switch app.Category {
		case CategoryGuest:
			if app.Docker != nil {
				tmpl.Symbol = SymbolDockerGuest
			} else {
				tmpl.Symbol = SymbolQemuGuest
			}
		case CategoryRouter:
			tmpl.Symbol = SymbolRouter
		case CategorySwitch:
			tmpl.Symbol = SymbolEthernetSwitch
		case CategoryMultilayerSwitch:
			tmpl.Symbol = SymbolMultilayerSwitch
		case CategoryFirewall:
			tmpl.Symbol = SymbolFirewall
		}
	}

	if app.Qemu != nil {
		tmpl.TemplateType = TemplateTypeQemu
		err := a.addQemuConfigToTemplate(tmpl, &app)
		if err != nil {
			return nil, nil, err
		}
	} else if app.Iou != nil {
		tmpl.TemplateType = TemplateTypeIou
	} else if app.Dynamips != nil {
		tmpl.TemplateType = TemplateTypeDynamips
	} else if app.Docker != nil {
		tmpl.TemplateType = TemplateTypeDocker
		err := a.addDockerConfigToTemplate(tmpl, &app)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, fmt.Errorf("%s no configuration found for known emulators", tmpl.Name)
	}

	// Download images
	for _, img := range app.Images {
		if img.DirectDownloadUrl != "" {
			err := a.checkAndDownload(img.DirectDownloadUrl, ImageTypes[tmpl.TemplateType], img.Filename, img.Md5Sum)
			if err != nil {
				fmt.Printf("ERROR: check '%s': %v\n", img.DirectDownloadUrl, err)
				return nil, nil, err
			}
			fmt.Printf("INFO: '%s', downloaded and verified\n", img.Filename)
		} else if img.DownloadUrl != "" {
			err := a.checkAndDownload(img.DownloadUrl, ImageTypes[tmpl.TemplateType], img.Filename, img.Md5Sum)
			if err != nil {
				fmt.Printf("ERROR: check '%s': %v\n", img.DownloadUrl, err)
				return nil, nil, err
			}
			fmt.Printf("INFO: '%s', downloaded and verified\n", img.Filename)
		}
	}

	return &app, tmpl, nil
}

func (a *Appliances) addDockerConfigToTemplate(tmpl *Template, app *Appliance) error {
	tmpl.Docker = &TemplateDocker{
		Adapters: app.Docker.Adapters,
		Image:    app.Docker.Image,
	}
	return nil
}

func (a *Appliances) addQemuConfigToTemplate(tmpl *Template, app *Appliance) error {
	// copy over the qemu appliance config to the template
	tmpl.Qemu = &TemplateQemu{
		AdapterType:       app.Qemu.AdapterType,
		Adapters:          app.Qemu.Adapters,
		BootPriority:      app.Qemu.BootPriority,
		ConsoleType:       app.Qemu.ConsoleType,
		HdaDiskInterface:  app.Qemu.HdaDiskInterface,
		Ram:               app.Qemu.Ram,
		KernelCommandLine: app.Qemu.KernelCommandLine,
	}

	options := app.Qemu.Options
	if app.Qemu.Kvm == "disable" && !strings.Contains(options, "-machine accel=tcg") {
		options += " -machine accel=tcg"
	}

	tmpl.Qemu.Options = options
	tmpl.Qemu.PortNameFormat = app.PortNameFormat
	tmpl.Qemu.Path = "qemu-system-" + app.Qemu.Arch

	// load the images from the application version
	// we don't seem to have the image path or type in the appliance.
	// the version images has the type as the key followed by the image name that can be used for the qemu template

	if len(app.Versions) > 0 {
		tmpl.Qemu.BiosImage = app.Versions[0].Images.BiosImage
		tmpl.Qemu.HdaDiskImage = app.Versions[0].Images.HdaDiskImage
		tmpl.Qemu.CdromImage = app.Versions[0].Images.CdromImage
	}

	return nil
}
