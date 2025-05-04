/*
 *
 * Copyright 2021 gRPC authors.
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

package resolver

import (
	"encoding/base64"
	"sort"
	"strings"
)

<<<<<<< HEAD
type addressMapEntry[T any] struct {
	addr  Address
	value T
}

// AddressMap is an AddressMapV2[any].  It will be deleted in an upcoming
// release of grpc-go.
//
// Deprecated: use the generic AddressMapV2 type instead.
type AddressMap = AddressMapV2[any]

// AddressMapV2 is a map of addresses to arbitrary values taking into account
// Attributes.  BalancerAttributes are ignored, as are Metadata and Type.
// Multiple accesses may not be performed concurrently.  Must be created via
// NewAddressMap; do not construct directly.
type AddressMapV2[T any] struct {
=======
type addressMapEntry struct {
	addr  Address
	value any
}

// AddressMap is a map of addresses to arbitrary values taking into account
// Attributes.  BalancerAttributes are ignored, as are Metadata and Type.
// Multiple accesses may not be performed concurrently.  Must be created via
// NewAddressMap; do not construct directly.
type AddressMap struct {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	// The underlying map is keyed by an Address with fields that we don't care
	// about being set to their zero values. The only fields that we care about
	// are `Addr`, `ServerName` and `Attributes`. Since we need to be able to
	// distinguish between addresses with same `Addr` and `ServerName`, but
	// different `Attributes`, we cannot store the `Attributes` in the map key.
	//
	// The comparison operation for structs work as follows:
	//  Struct values are comparable if all their fields are comparable. Two
	//  struct values are equal if their corresponding non-blank fields are equal.
	//
	// The value type of the map contains a slice of addresses which match the key
	// in their `Addr` and `ServerName` fields and contain the corresponding value
	// associated with them.
<<<<<<< HEAD
	m map[Address]addressMapEntryList[T]
=======
	m map[Address]addressMapEntryList
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
}

func toMapKey(addr *Address) Address {
	return Address{Addr: addr.Addr, ServerName: addr.ServerName}
}

<<<<<<< HEAD
type addressMapEntryList[T any] []*addressMapEntry[T]

// NewAddressMap creates a new AddressMapV2[any].
//
// Deprecated: use the generic NewAddressMapV2 constructor instead.
func NewAddressMap() *AddressMap {
	return NewAddressMapV2[any]()
}

// NewAddressMapV2 creates a new AddressMapV2.
func NewAddressMapV2[T any]() *AddressMapV2[T] {
	return &AddressMapV2[T]{m: make(map[Address]addressMapEntryList[T])}
=======
type addressMapEntryList []*addressMapEntry

// NewAddressMap creates a new AddressMap.
func NewAddressMap() *AddressMap {
	return &AddressMap{m: make(map[Address]addressMapEntryList)}
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
}

// find returns the index of addr in the addressMapEntry slice, or -1 if not
// present.
<<<<<<< HEAD
func (l addressMapEntryList[T]) find(addr Address) int {
=======
func (l addressMapEntryList) find(addr Address) int {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	for i, entry := range l {
		// Attributes are the only thing to match on here, since `Addr` and
		// `ServerName` are already equal.
		if entry.addr.Attributes.Equal(addr.Attributes) {
			return i
		}
	}
	return -1
}

// Get returns the value for the address in the map, if present.
<<<<<<< HEAD
func (a *AddressMapV2[T]) Get(addr Address) (value T, ok bool) {
=======
func (a *AddressMap) Get(addr Address) (value any, ok bool) {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 {
		return entryList[entry].value, true
	}
<<<<<<< HEAD
	return value, false
}

// Set updates or adds the value to the address in the map.
func (a *AddressMapV2[T]) Set(addr Address, value T) {
=======
	return nil, false
}

// Set updates or adds the value to the address in the map.
func (a *AddressMap) Set(addr Address, value any) {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 {
		entryList[entry].value = value
		return
	}
<<<<<<< HEAD
	a.m[addrKey] = append(entryList, &addressMapEntry[T]{addr: addr, value: value})
}

// Delete removes addr from the map.
func (a *AddressMapV2[T]) Delete(addr Address) {
=======
	a.m[addrKey] = append(entryList, &addressMapEntry{addr: addr, value: value})
}

// Delete removes addr from the map.
func (a *AddressMap) Delete(addr Address) {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	entry := entryList.find(addr)
	if entry == -1 {
		return
	}
	if len(entryList) == 1 {
		entryList = nil
	} else {
		copy(entryList[entry:], entryList[entry+1:])
		entryList = entryList[:len(entryList)-1]
	}
	a.m[addrKey] = entryList
}

// Len returns the number of entries in the map.
<<<<<<< HEAD
func (a *AddressMapV2[T]) Len() int {
=======
func (a *AddressMap) Len() int {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	ret := 0
	for _, entryList := range a.m {
		ret += len(entryList)
	}
	return ret
}

// Keys returns a slice of all current map keys.
<<<<<<< HEAD
func (a *AddressMapV2[T]) Keys() []Address {
=======
func (a *AddressMap) Keys() []Address {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	ret := make([]Address, 0, a.Len())
	for _, entryList := range a.m {
		for _, entry := range entryList {
			ret = append(ret, entry.addr)
		}
	}
	return ret
}

// Values returns a slice of all current map values.
<<<<<<< HEAD
func (a *AddressMapV2[T]) Values() []T {
	ret := make([]T, 0, a.Len())
=======
func (a *AddressMap) Values() []any {
	ret := make([]any, 0, a.Len())
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	for _, entryList := range a.m {
		for _, entry := range entryList {
			ret = append(ret, entry.value)
		}
	}
	return ret
}

type endpointMapKey string

// EndpointMap is a map of endpoints to arbitrary values keyed on only the
// unordered set of address strings within an endpoint. This map is not thread
// safe, thus it is unsafe to access concurrently. Must be created via
// NewEndpointMap; do not construct directly.
<<<<<<< HEAD
type EndpointMap[T any] struct {
	endpoints map[endpointMapKey]endpointData[T]
}

type endpointData[T any] struct {
	// decodedKey stores the original key to avoid decoding when iterating on
	// EndpointMap keys.
	decodedKey Endpoint
	value      T
}

// NewEndpointMap creates a new EndpointMap.
func NewEndpointMap[T any]() *EndpointMap[T] {
	return &EndpointMap[T]{
		endpoints: make(map[endpointMapKey]endpointData[T]),
=======
type EndpointMap struct {
	endpoints map[endpointMapKey]endpointData
}

type endpointData struct {
	// decodedKey stores the original key to avoid decoding when iterating on
	// EndpointMap keys.
	decodedKey Endpoint
	value      any
}

// NewEndpointMap creates a new EndpointMap.
func NewEndpointMap() *EndpointMap {
	return &EndpointMap{
		endpoints: make(map[endpointMapKey]endpointData),
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	}
}

// encodeEndpoint returns a string that uniquely identifies the unordered set of
// addresses within an endpoint.
func encodeEndpoint(e Endpoint) endpointMapKey {
	addrs := make([]string, 0, len(e.Addresses))
	// base64 encoding the address strings restricts the characters present
	// within the strings. This allows us to use a delimiter without the need of
	// escape characters.
	for _, addr := range e.Addresses {
		addrs = append(addrs, base64.StdEncoding.EncodeToString([]byte(addr.Addr)))
	}
	sort.Strings(addrs)
	// " " should not appear in base64 encoded strings.
	return endpointMapKey(strings.Join(addrs, " "))
}

// Get returns the value for the address in the map, if present.
<<<<<<< HEAD
func (em *EndpointMap[T]) Get(e Endpoint) (value T, ok bool) {
=======
func (em *EndpointMap) Get(e Endpoint) (value any, ok bool) {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	val, found := em.endpoints[encodeEndpoint(e)]
	if found {
		return val.value, true
	}
<<<<<<< HEAD
	return value, false
}

// Set updates or adds the value to the address in the map.
func (em *EndpointMap[T]) Set(e Endpoint, value T) {
	en := encodeEndpoint(e)
	em.endpoints[en] = endpointData[T]{
=======
	return nil, false
}

// Set updates or adds the value to the address in the map.
func (em *EndpointMap) Set(e Endpoint, value any) {
	en := encodeEndpoint(e)
	em.endpoints[en] = endpointData{
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
		decodedKey: Endpoint{Addresses: e.Addresses},
		value:      value,
	}
}

// Len returns the number of entries in the map.
<<<<<<< HEAD
func (em *EndpointMap[T]) Len() int {
=======
func (em *EndpointMap) Len() int {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	return len(em.endpoints)
}

// Keys returns a slice of all current map keys, as endpoints specifying the
// addresses present in the endpoint keys, in which uniqueness is determined by
// the unordered set of addresses. Thus, endpoint information returned is not
// the full endpoint data (drops duplicated addresses and attributes) but can be
// used for EndpointMap accesses.
<<<<<<< HEAD
func (em *EndpointMap[T]) Keys() []Endpoint {
=======
func (em *EndpointMap) Keys() []Endpoint {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	ret := make([]Endpoint, 0, len(em.endpoints))
	for _, en := range em.endpoints {
		ret = append(ret, en.decodedKey)
	}
	return ret
}

// Values returns a slice of all current map values.
<<<<<<< HEAD
func (em *EndpointMap[T]) Values() []T {
	ret := make([]T, 0, len(em.endpoints))
=======
func (em *EndpointMap) Values() []any {
	ret := make([]any, 0, len(em.endpoints))
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	for _, val := range em.endpoints {
		ret = append(ret, val.value)
	}
	return ret
}

// Delete removes the specified endpoint from the map.
<<<<<<< HEAD
func (em *EndpointMap[T]) Delete(e Endpoint) {
=======
func (em *EndpointMap) Delete(e Endpoint) {
>>>>>>> e302735 ([backend] generate vendor folders for backend services)
	en := encodeEndpoint(e)
	delete(em.endpoints, en)
}
