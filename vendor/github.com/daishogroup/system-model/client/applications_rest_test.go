//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Application rest tests

package client

import (
    "github.com/stretchr/testify/suite"
    "testing"
    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/dhttp"
)

type ApplicationRestTestSuite struct {
    suite.Suite
    client Applications
    rest * dhttp.ClientMockup
}

func (suite * ApplicationRestTestSuite) SetupSuite() {
    rest := dhttp.NewClientMockup()
    suite.rest = rest
    client := ApplicationsRest{rest}
    suite.client = &client
}

func (suite * ApplicationRestTestSuite) SetupTest() {
    suite.rest.Reset()
}

func TestApplicationRest(t *testing.T){
    suite.Run(t, new(ApplicationRestTestSuite))
}

func (suite * ApplicationRestTestSuite) getDescriptor(networkID string, descriptorID string) * entities.AppDescriptor{
    return entities.NewAppDescriptorWithID(networkID,
        descriptorID, TestName, TestDescription, TestServiceName, TestServiceVersion, TestLabel, TestPort, []string{TestImage})
}

func (suite * ApplicationRestTestSuite) getInstance(networkID string, deployedID string) * entities.AppInstance {
    return entities.NewAppInstanceWithID(networkID,
        deployedID, TestDescriptorID, "",
        TestName, TestDescription, TestLabel, TestArguments, entities.AppInstReady, TestPersistenceSize, TestStorageType,
        make([]entities.ApplicationPort, 0), TestPort, "")
}

func (suite * ApplicationRestTestSuite) TestAddDescriptor(){
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := suite.getDescriptor("n1","randomDescId")
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestAddDescriptor(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestListDescriptors(){
    suite.rest.AddGet(func(path string) dhttp.Response{
        desc := suite.getDescriptor("n1",TestDescriptorID)
        result := [] entities.AppDescriptor { *desc }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListDescriptors(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestGetDescriptor(){
    suite.rest.AddGet(func(path string) dhttp.Response{
        result := suite.getDescriptor("n1",TestDescriptorID)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestGetDescriptor(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestAddInstance(){
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := suite.getInstance("n1","randomDescId")
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestAddInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestListInstances(){
    suite.rest.AddGet(func(path string) dhttp.Response{
        inst1 := suite.getInstance("n1",TestInstanceID)
        inst2 := suite.getInstance("n1",TestInstanceIDUpdate)
        result := [] entities.AppInstance { *inst1, *inst2 }
        statusCode := 200
        return dhttp.NewResponse(&result, &statusCode, nil)
    })
    TestListInstances(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestGetInstance(){
    suite.rest.AddGet(func(path string) dhttp.Response{
        result := suite.getInstance("n1",TestInstanceID)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestGetInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestUpdateInstance(){
    suite.rest.AddGet(func(path string) dhttp.Response{
        result := suite.getInstance("n1",TestInstanceIDUpdate)
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response{
        result := suite.getInstance("n1",TestInstanceIDUpdate)
        result.Description = "new description"
        statusCode := 200
        return dhttp.NewResponse(result, &statusCode, nil)
    })
    TestUpdateInstance(&suite.Suite, suite.client)
}

func (suite * ApplicationRestTestSuite) TestDeleteInstance(){
    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        result:=entities.NewSuccessfulOperation("DeleteInstance")
        return dhttp.NewResponse(result,&statusCode, nil)
    })
    TestDeleteInstance(&suite.Suite, suite.client)
}



