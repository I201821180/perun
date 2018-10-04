// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
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

package stack

import (
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyStackPolicy(t *testing.T) {
	stackName := "StackName"
	policyPath := "./test_resources/test_stackpolicy.json"
	ctx := mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, policyPath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	template := mocks.ReadFile(t, policyPath)

	input := createStackPolicyInput(&template, &stackName)
	mockAWSPI.EXPECT().SetStackPolicy(&input).Return(nil, nil).Times(1)

	ApplyStackPolicy(ctx)
}

func TestCreateStackPolicyInput(t *testing.T) {
	stackName := "StackName"
	templateBody := "TestTemplate"
	returnedValue := createStackPolicyInput(&templateBody, &stackName)
	assert.Equal(t, *returnedValue.StackName, stackName)
	assert.Equal(t, *returnedValue.StackPolicyBody, templateBody)
}
