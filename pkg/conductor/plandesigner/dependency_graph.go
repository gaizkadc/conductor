/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/nalej/conductor/internal/entities"
    "github.com/yourbasic/graph"
)

// DependencyGraph represents the connectivity between different services inside a Nalej application.
type DependencyGraph struct {
    // The internal graph object
    graph *graph.Mutable
    // A reference map to translate between services id and vertex ids
    reference map[string]int
}

// NewDeopendencyGraph build a dependency graph from a list of service entities. This graphs represents
// the relationship between the entities in the temporal space. This is,
func NewDependencyGraph(services []entities.Service) *DependencyGraph {
    // Create the graph, one vertex per service
    // use the indexes in the array as ids in the graph
    g := graph.New(len(services))
    // Build a map to translate serviceid->array position
    reference := make(map[string]int,0)
    for i, serv := range services {
        reference[serv.ServiceId] = i
    }
    for _, serv := range services {
        if serv.DeployAfter != nil && len(serv.DeployAfter) >0 {
            sourceVertex := reference[serv.ServiceId]
            for _, afterId := range serv.DeployAfter {
                targetVertex := reference[afterId]
                //g.Add(sourceVertex, targetVertex)
                g.Add(targetVertex, sourceVertex)
            }
        }
    }
    return &DependencyGraph{graph: g, reference: reference}
}

func (dg *DependencyGraph) NumServices() int {
    return dg.graph.Order()
}

func (dg *DependencyGraph) NumDependencies() int {
    sum := 0
    for i:=0; i< dg.graph.Order(); i++ {
        sum = sum + dg.graph.Degree(i)
    }
    return sum
}

func (dg *DependencyGraph) GetDependencyOrderByGroups() [][]int {
    // Follow https://stackoverflow.com/questions/4073119/topological-sort-with-grouping
    //order, _ := graph.TopSort(dg.graph)
    //return order

    // This must be an acyclic graph
    isAcyclic := graph.Acyclic(dg.graph)
    if !isAcyclic {
        return nil
    }

    groups := make([]int, dg.NumServices())
    // mark nodes with degree 0 as root
    for i := 0; i < dg.NumServices(); i++ {
        groups[i] = 0
    }


    // store the max groups id to generate the corresponding data structure to be returned
    // log.Debug().Msgf("initial list of groups %v",groups)
    maxGroupId := 0
    changes := true
    for changes {
        changes = false
        for i := 0; i < dg.NumServices(); i++ {
            // Use breadth-first-search to iterate the graph and change the groups
            graph.BFS(dg.graph, i, func(v int, w int, _ int64){
                newGroupId := groups[v] + 1
                //log.Debug().Msgf("from %d to %v",v,w)
                if groups[w] < newGroupId {
                    //log.Debug().Msgf("--> %d changes group to %d",w, newGroupId)
                    groups[w] = newGroupId
                    if newGroupId > maxGroupId {
                        maxGroupId = newGroupId
                    }
                    changes = true
                }
            })
        }

    }

    //log.Debug().Msgf("list of groups %v",groups)

    toReturn := make([][]int,maxGroupId+1)
    // fill the list of groups
    for index, group := range groups {
        // log.Debug().Msgf("node %d goes to group %d", index, group)
        if toReturn[group] == nil {
            toReturn[group] = make([]int,0)
        }
        toReturn[group] = append(toReturn[group],index)
    }

    return toReturn
}

func (dg *DependencyGraph) String() string {
    return dg.graph.String()
}