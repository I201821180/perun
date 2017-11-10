// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cfofflinevalidator privides tools for offline cloudformation template
// validation.
package cfofflinevalidator

import (
	"encoding/json"
	"github.com/Appliscale/cftool/cfspecification"
	"io/ioutil"
	"path"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/Appliscale/cftool/cfofflinevalidator/cftemplate"
	"github.com/Appliscale/cftool/cfofflinevalidator/cfvalidators"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfcontext"
)

var validators = map[string]interface{}{
	"AWS::EC2::VPC": cfvalidators.IsVpcValid,
}

// Validate cloudformation template.
func Validate(context *cfcontext.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	specification, err := cfspecification.GetSpecification(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	var template cftemplate.Template
	templateFileExtension := path.Ext(*context.CliArguments.TemplatePath)
	if templateFileExtension == ".json" {
		template, err = parseJSON(rawTemplate)
	} else if templateFileExtension == ".yaml" ||  templateFileExtension == ".yml" {
		template, err = parseYAML(rawTemplate)
	} else {
		err = errors.New("Invalid template file format.")
	}
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	valid = validateResources(template.Resources, &specification, context.Logger)
	return valid
}

func printResult(valid *bool, logger *cflogger.Logger) {
	logger.PrintValidationErrors()
	if !*valid {
		logger.Info("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}

func validateResources(resources map[string]cftemplate.Resource, specification *cfspecification.Specification, logger *cflogger.Logger) (bool) {
	valid := true
	for resourceName, resourceValue := range resources {
		if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
			if !areRequiredPropertiesPresent(resourceSpecification, resourceValue, resourceName, logger) {
				valid = false
			}
		} else {
			logger.ValidationError(resourceName, "Type needs to be specified")
			valid = false
		}
		if validator, ok := validators[resourceValue.Type]; ok {
			if !validator.(func(string, cftemplate.Resource, *cflogger.Logger)(bool))(resourceName, resourceValue, logger) {
				valid = false
			}
		}
	}

	return valid
}
func areRequiredPropertiesPresent(resourceSpecification cfspecification.Resource, resourceValue cftemplate.Resource, resourceName string, logger *cflogger.Logger) bool {
	valid := true
	for propertyName, propertyValue := range resourceSpecification.Properties {
		if propertyValue.Required {
			if _, ok := resourceValue.Properties[propertyName]; !ok {
				logger.ValidationError(resourceName, "Property " + propertyName + " is required")
				valid = false
			}
		}
	}
	return valid
}

func parseJSON(templateFile []byte) (template cftemplate.Template, err error) {
	err = json.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}

func parseYAML(templateFile []byte) (template cftemplate.Template, err error) {
	err = yaml.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}
