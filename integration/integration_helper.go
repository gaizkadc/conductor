//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Main entry command line to run the integration for conductor.

package integration

// This is a collection of routines used during the integration test.

import (
    "os/exec"
    "context"
    "time"
    "os"
    "strconv"
)

// Timeout to wait for colony to finish in seconds
const (
    ColonyTimeout = 35*60
    ASMTimeout = 3*60
)

// Run an ASM command.
// params:
//    arguments additional arguments for the ASM
// return:
//    byte array coming from the output
//    error if any
func(i* Integration) runASMCommand(arguments ...string ) ([] byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), ASMTimeout*time.Second)
    defer cancel()

    tokens := [] string {"--ip="+AppMgrIP, "--port="+strconv.Itoa(AppMgrPort)}
    tokens = append(tokens, arguments...)
    cmd := exec.CommandContext(ctx, "./asmcli", tokens...)
    cmd.Dir = i.configuration.AsmPath
    return cmd.CombinedOutput()
}

// Run a colony command and connect the standard output.
// params:
//    arguments additional arguments to colony
// return:
//    error if any
func(i* Integration) runColonyCommandStdout(arguments ...string) error {
    ctx, cancel := context.WithTimeout(context.Background(), ColonyTimeout*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "./colonize", arguments...)
    // set dir to run in the colony folder
    cmd.Dir = i.configuration.ColonyPath
    // redirect to the standard output
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err:= cmd.Run(); err != nil {
        logger.Fatal(err)
        return err
    }

    return nil
}


// Run a colony command and return the obtained output
// params:
//    arguments additional arguments to colony
// return:
//    output from the command line in a byte array
//    error if any
func(i* Integration) runColonyCommand(arguments ...string) ([] byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), ColonyTimeout*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "./colonize", arguments...)
    // set dir to run in the colony folder
    cmd.Dir = i.configuration.ColonyPath

    return cmd.CombinedOutput()
}


