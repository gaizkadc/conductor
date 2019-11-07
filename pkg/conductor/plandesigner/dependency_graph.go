/*
 * Copyright 2019 Nalej
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
 *
 */

package plandesigner

import (
	"errors"
	"github.com/nalej/conductor/internal/entities"
	"github.com/yourbasic/graph"
)

// DependencyGraph represents the connectivity between different services inside a Nalej application.
type DependencyGraph struct {
	// The internal graph object
	graph *graph.Mutable
	// A id2graph map to translate between services name and vertex ids
	id2graph map[string]int
	// Array where the i-th node id corresponds to the i-th service  in the array
	graph2id []entities.Service
}

// NewDeopendencyGraph build a dependency graph from a list of service entities. This graphs represents
// the relationship between the entities in the temporal space. This is,
func NewDependencyGraph(services []entities.Service) *DependencyGraph {
	// Create the graph, one vertex per service
	// use the indexes in the array as ids in the graph
	g := graph.New(len(services))
	// Build a map to translate serviceName->array position
	reference := make(map[string]int, 0)
	// Build the array to translate nodeid -> serviceName
	greference := make([]entities.Service, len(services))
	for i, serv := range services {
		reference[serv.Name] = i
		greference[i] = serv
	}
	for _, serv := range services {
		if serv.DeployAfter != nil && len(serv.DeployAfter) > 0 {
			sourceVertex := reference[serv.Name]
			for _, afterName := range serv.DeployAfter {
				targetVertex := reference[afterName]
				//g.Add(sourceVertex, targetVertex)
				// create graph in temporal order
				g.Add(targetVertex, sourceVertex)
			}
		}
	}
	return &DependencyGraph{graph: g, id2graph: reference, graph2id: greference}
}

func (dg *DependencyGraph) NumServices() int {
	return dg.graph.Order()
}

func (dg *DependencyGraph) NumDependencies() int {
	sum := 0
	for i := 0; i < dg.graph.Order(); i++ {
		sum = sum + dg.graph.Degree(i)
	}
	return sum
}

// GetDependencyOrderByGraph returns a array of arrays with the stages and the ids of the services that can be
// executed in parallel in at the same time when the previous stage finishes.
// return:
//  array of services per group. E.g.: [[service2,service3], [service0], [service1]]
func (dg *DependencyGraph) GetDependencyOrderByGroups() ([][]entities.Service, error) {

	// If there is only one node, simply return it
	if dg.NumServices() == 1 {
		return [][]entities.Service{dg.graph2id}, nil
	}

	// This must be an acyclic graph
	isAcyclic := graph.Acyclic(dg.graph)
	if !isAcyclic {
		error := errors.New("cyclic dependency graph")
		return nil, error
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
			graph.BFS(dg.graph, i, func(v int, w int, _ int64) {
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

	toReturn := make([][]entities.Service, maxGroupId+1)
	// fill the list of groups
	for index, group := range groups {
		// log.Debug().Msgf("node %d goes to group %d", index, group)
		if toReturn[group] == nil {
			toReturn[group] = make([]entities.Service, 0)
		}
		toReturn[group] = append(toReturn[group], dg.graph2id[index])
	}

	return toReturn, nil
}

func (dg *DependencyGraph) String() string {
	return dg.graph.String()
}
