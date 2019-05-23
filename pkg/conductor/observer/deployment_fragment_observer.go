/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package observer

import (
    "context"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/internal/persistence/app_cluster"
    "github.com/nalej/derrors"
    "github.com/rs/zerolog/log"
    "time"
)

const CheckSleepTime = time.Second

// This observer can be used to wait for certain deployment fragments and take actions.


// Auxiliary structure
type ObservableDeploymentFragment struct {
    ClusterId string
    FragmentId string
    AppInstanceId string
}

type DeploymentFragmentsObserver struct {
    // Map with the clusterId -> appInstanceId
    Ids []ObservableDeploymentFragment
    // Monitored apps
    AppClusterDB *app_cluster.AppClusterDB
    // Remaining changes to occur
    RemainingChanges int
}

func NewDeploymentFragmentsObserver(ids []ObservableDeploymentFragment, appClusterDB *app_cluster.AppClusterDB) DeploymentFragmentsObserver {
    return DeploymentFragmentsObserver{Ids: ids, AppClusterDB: appClusterDB, RemainingChanges: len(ids)}
}

// Observe changes in the list of observed deployment fragments and run the indicated function if the deployment fragment
// changes into the given status type. The observer will stop when all the deployment fragments have been observed to
// change or when a timeout is reached.
// params:
//  timeout duration for the timeout of this context
//  status to be detected
//  f function to be called when the defined status is found
// return:
//  error if any
func (df * DeploymentFragmentsObserver) Observe(timeout time.Duration, status entities.DeploymentFragmentStatus, f func(*entities.DeploymentFragment) derrors.Error) {
    log.Debug().Interface("observableItems",df.Ids).Msgf("started deployments fragment observer with %d pending observations",df.RemainingChanges)
    sleep := time.Tick(CheckSleepTime)
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // store here processed items
    processed := make(map[string]bool,len(df.Ids))
    for _, id := range df.Ids {
        processed[id.FragmentId] = false
    }

    for {
        select {
        case <-sleep:
            for _, observed := range df.Ids {
                clusterId := observed.ClusterId
                fragmentId := observed.FragmentId

                // skip if we have already processed this entry
                done, _ := processed[fragmentId]
                if done {
                    continue
                }

                fragment, err := df.AppClusterDB.GetDeploymentFragment(clusterId, fragmentId)

                if fragment == nil {
                    log.Debug().Msgf("no deployment fragment stored for cluster %s with id %s", clusterId, fragmentId)
                    continue
                }

                if err != nil {
                    log.Error().Err(err).Str("clusterId",clusterId).Str("fragmentId", fragmentId).
                        Msg("error when collecting fragment data")
                    continue
                }

                if fragment.Status == status {
                    if e:= f(fragment); e!= nil {
                        log.Error().Err(err).Msg("error when executing callback function after observing change")
                    }
                    // set this entry as processed
                    processed[fragmentId] = true

                    // one observed reduce the counter
                    df.RemainingChanges = df.RemainingChanges - 1
                    log.Debug().Msgf("remaining %d deployment fragments to observe",df.RemainingChanges)

                    if df.RemainingChanges == 0 {
                        log.Debug().Msg("deployment fragments observer stops after all the elements were processed")
                        return
                    }
                }
            }

        case <- ctx.Done():
            log.Debug().Interface("observableItems",df.Ids).Interface("processed",processed).
                Msg("timeout reached for deployment fragments observer")
            return
        }
    }
}
