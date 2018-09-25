/*
 * Copyright 2018 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conductor

// Set of common routines for conductor components. A pool of already opened client connections is maintained
// for the components below and implemented in a singleton instance accessible by all the elements in this package.
// When running tests, this pool uses listening buffers.

import (
    "github.com/nalej/conductor/tools"
    "google.golang.org/grpc"
    "github.com/rs/zerolog/log"
    "sync"
    //"flag"
)

var (
    // Singleton instance of connections with musician clients
    MusicianClients *tools.ConnectionsMap
    onceMusicians   sync.Once
    // Singleton instance of connections with deployment managers
    DMClients *tools.ConnectionsMap
    onceDM sync.Once
)


func GetMusicianClients() *tools.ConnectionsMap {
    /*
    onceMusicians.Do(func(){
        if flag.Lookup("test.v") != nil {
            log.Debug().Msg("using testing musician clients factory")
            MusicianClients = tools.NewConnectionsMap(conductorClientFactoryTest)
        } else {
            MusicianClients = tools.NewConnectionsMap(conductorClientFactory)
        }
    })
    */

    onceMusicians.Do(func(){
        MusicianClients = tools.NewConnectionsMap(conductorClientFactory)
    })

    return MusicianClients
}

func GetDMClients() *tools.ConnectionsMap {
    onceDM.Do(func() {
        DMClients = tools.NewConnectionsMap(dmClientFactory)
    })
    return DMClients
}

/*
// Factory in charge of generating new connections for Conductor->Musician communication in test environments.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func conductorClientFactoryTest(address string) (*grpc.ClientConn, error) {
    conn, err := tools.GetConn(*tools.GetDefaultListener())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}
*/

// Factory in charge of generating new connections for Conductor->Musician communication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func conductorClientFactory(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}

// Factory in charge of generating new connections for Conductor->DM communication.
//  params:
//   address the communication has to be done with
//  return:
//   client and error if any
func dmClientFactory(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}
