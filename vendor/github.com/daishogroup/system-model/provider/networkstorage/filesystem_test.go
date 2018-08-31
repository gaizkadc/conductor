//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// File system tests.

package networkstorage

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/user"
    "path"
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
)

func TestFileSystemProvider_Networks(t *testing.T) {

    dir, ioError := ioutil.TempDir("", "testNetwork")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    toAdd1 := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    toAdd2 := entities.NewNetwork("n2", "desc", "admin2", "2", "admin2")
    err := provider.Add(* toAdd1)
    assert.Nil(t, err, "network1 should be added")
    err = provider.Add(* toAdd2)
    assert.Nil(t, err, "network2 should be added")

    r1, err := provider.RetrieveNetwork(toAdd1.ID)
    assert.Nil(t, err, "network1 should be retrieved")
    assert.EqualValues(t, toAdd1, r1, "network1 should match")
    r2, err := provider.RetrieveNetwork(toAdd2.ID)
    assert.Nil(t, err, "network1 should be retrieved")
    assert.EqualValues(t, toAdd2, r2, "network1 should match")

    networks, err := provider.ListNetworks()
    assert.Nil(t, err, "network list should be retrieved")
    assert.Equal(t, 2, len(networks))

}

func TestFileSystemProvider_Clusters(t *testing.T) {

    dir, ioError := ioutil.TempDir("", "testNetwork")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.Nil(t, err, "network1 should be added")

    // List clusters
    clusters, err := provider.ListClusters(net.ID)
    assert.Nil(t, err, "clusters should be retrieved")
    assert.Equal(t, 0, len(clusters))

    c1 := entities.NewCluster(net.ID, "c1", "d1", entities.EdgeType, "l",
        "email1", entities.ClusterCreated, false, false)
    c2 := entities.NewCluster(net.ID, "c1", "d1", entities.EdgeType, "l",
        "email1", entities.ClusterCreated, false, false)
    err = provider.AttachCluster(net.ID, c1.ID)
    assert.Nil(t, err, "cluster should be attached")
    err = provider.AttachCluster(net.ID, c2.ID)
    assert.Nil(t, err, "cluster should be attached")
    clusters, err = provider.ListClusters(net.ID)
    assert.Nil(t, err, "clusters should be retrieved")
    assert.Equal(t, 2, len(clusters))
    err = provider.DeleteCluster(net.ID, c1.ID)
    assert.Nil(t, err, "cluster should be detached")
    clusters, err = provider.ListClusters(net.ID)
    assert.Nil(t, err, "clusters should be retrieved")
    assert.Equal(t, 1, len(clusters))

}

func TestFileSystemProvider_Descriptors(t *testing.T) {

    dir, ioError := ioutil.TempDir("", "testNetwork")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.Nil(t, err, "network1 should be added")

    descriptors, err := provider.ListAppDesc(net.ID)
    assert.Nil(t, err, "descriptors should be retrieved")
    assert.Equal(t, 0, len(descriptors))

    d1 := entities.NewAppDescriptor(net.ID, "n1", "d1","s1","1.0","l",0,[]string{"repo1:tag1"})
    d2 := entities.NewAppDescriptor(net.ID, "n2", "d2","s2","1.0","l",0,[]string{"repo2:tag2"})

    err = provider.RegisterAppDesc(net.ID, d1.ID)
    assert.Nil(t, err, "descriptor should be attached")
    err = provider.RegisterAppDesc(net.ID, d2.ID)
    assert.Nil(t, err, "descriptor should be attached")
    descriptors, err = provider.ListAppDesc(net.ID)
    assert.Nil(t, err, "descriptors should be retrieved")
    assert.Equal(t, 2, len(descriptors))
}

func TestFileSystemProvider_Instances(t *testing.T) {

    dir, ioError := ioutil.TempDir("", "testNetwork")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.Nil(t, err, "network1 should be added")

    instances, err := provider.ListAppInst(net.ID)
    assert.Nil(t, err, "instances should be retrieved")
    assert.Equal(t, 0, len(instances))

    i1 := entities.NewAppInstance(net.ID, "d1", "c1", "n1", "d1", "l1",
        "arg","1Gb", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 0, "0.0.0.0")
    i2 := entities.NewAppInstance(net.ID, "d2", "c2", "n2", "d2", "l2",
        "arg","1Gb", entities.AppStorageDefault, make([]entities.ApplicationPort, 0), 0, "0.0.0.0")
    err = provider.RegisterAppInst(net.ID, i1.DeployedID)
    assert.Nil(t, err, "instance should be attached")
    err = provider.RegisterAppInst(net.ID, i2.DeployedID)
    assert.Nil(t, err, "instance should be attached")
    instances, err = provider.ListAppInst(net.ID)
    assert.Nil(t, err, "instances should be retrieved")
    assert.Equal(t, 2, len(instances))
    err = provider.DeleteAppInstance(net.ID, i1.DeployedID)
    assert.Nil(t, err, "instance should be detached")
    instances, err = provider.ListAppInst(net.ID)
    assert.Nil(t, err, "instances should be retrieved")
    assert.Equal(t, 1, len(instances))

}

func TestFileSystemProvider_Delete(t *testing.T) {
    dir, ioError := ioutil.TempDir("", "testNetwork")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.Nil(t, err, "network1 should be added")

    err = provider.DeleteNetwork(net.ID)
    assert.Nil(t,err, "network1 should have been deleted")

    // try to retrieve it
    retrieved, err := provider.RetrieveNetwork(net.ID)
    assert.NotNil(t, err, "network1 was already recovered after deletion")
    assert.Nil(t, retrieved, "network1 was returned")
}

func TestReadOnlyFileSystem(t *testing.T) {
    currentUser, err := user.Current()
    assert.Nil(t, err, "should be able to retrieve current user")
    if currentUser.Username == "root" {
        t.Skip("This test cannot be executed under the root user")
    }else {
        doTestReadOnlyFileSystem(t)
    }
}

func doTestReadOnlyFileSystem(t *testing.T) {
    dir, ioError := ioutil.TempDir("", "testFSPermissions")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)
    roPath := path.Join(dir, "readOnly")
    ioError = os.Mkdir(roPath, 000)
    if ioError != nil {
        fmt.Println("Cannot create read only directory", ioError.Error())
    }
    assert.Nil(t, ioError, "should be able to create directories.")
    fmt.Println("Temp file system", dir)
    provider := NewFileSystemProvider(roPath)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.NotNil(t, err, "network1 should not be added")
    fmt.Println(err.DebugReport())
    os.Chmod(roPath, 777)
}

func TestNewBasePath(t *testing.T) {
    dir, ioError := ioutil.TempDir("", "testFSPermissions")
    if ioError != nil {
        fmt.Println("Cannot create temp dir", ioError.Error())
    }
    assert.Nil(t, ioError, "directory should be created")
    defer os.RemoveAll(dir)
    targetPath := path.Join(dir, "doesNotExists/doesNotExists")
    fmt.Println("Temp file system", dir)
    provider := NewFileSystemProvider(targetPath)
    net := entities.NewNetwork("n1", "desc", "admin1", "1", "admin1")
    err := provider.Add(* net)
    assert.Nil(t, err, "network1 should be added")
}