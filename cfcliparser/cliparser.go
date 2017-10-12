package cfcliparser

import (
	"flag"
	"errors"
	"strings"
)

const ValidateMode string = "validate"
const ConvertMode string = "convert"
const JSON string = "json"
const YAML string = "yaml"

type CliArguments struct {
	Mode *string
	FilePath *string
	OutputFilePath *string
	OutputFileFormat *string
	Region *string
}

func ParseCliArguments() (cliArguments CliArguments, err error) {
	cliArguments.Mode = flag.String("mode", "", "Mode: " + ValidateMode + "|" + ConvertMode)
	cliArguments.FilePath = flag.String("file", "", "A path to the template")
	cliArguments.OutputFilePath = flag.String("output", "", "A path, where converted file will be saved")
	cliArguments.OutputFileFormat = flag.String("format", "", "Output format: " + strings.ToUpper(JSON) +
		"|" + strings.ToUpper(YAML))
	cliArguments.Region = flag.String("region", "", "Region (e.g. eu-central-1)")

	flag.Parse()

	if *cliArguments.Mode == "" {
		err = errors.New("You should specify what you want to do with -mode flag")
		return
	}

	if *cliArguments.Mode != ValidateMode && *cliArguments.Mode != ConvertMode {
		err = errors.New("Invalid mode. Use validate or convert")
		return
	}

	if *cliArguments.FilePath == "" {
		err = errors.New("You should specify a source of the template file with -file flag")
		return
	}

	if *cliArguments.Mode == ValidateMode {
		if *cliArguments.Region == "" {
			err = errors.New("You should specify a region with -region flag")
			return
		}
	}

	if *cliArguments.Mode == ConvertMode {
		if *cliArguments.OutputFilePath == "" {
			err = errors.New("You should specify a output file path with -output flag")
			return
		}

		if *cliArguments.OutputFileFormat == "" {
			err = errors.New("You should specify a output file format with -format flag")
			return
		}

		*cliArguments.OutputFileFormat = strings.ToLower(*cliArguments.OutputFileFormat)
		if *cliArguments.OutputFileFormat != JSON && *cliArguments.OutputFileFormat != YAML {
			err = errors.New("Invalid mode. Use validate or convert")
			return
		}
	}

	return
}