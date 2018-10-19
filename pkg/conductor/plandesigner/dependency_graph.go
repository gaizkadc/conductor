/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/nalej/conductor/internal/entities"
    "github.com/yourbasic/graph"
)

type DependencyGraph struct {
    // The internal graph object
    graph *graph.Mutable
    // A reference map to translate between services id and vertex ids
    reference map[string]int
}

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

func (dg *DependencyGraph) GetDependencyOrder() []int {
    // Follow https://stackoverflow.com/questions/4073119/topological-sort-with-grouping
    order, _ := graph.TopSort(dg.graph)
    return order
}

func (dg *DependencyGraph) String() string {
    return dg.graph.String()
}