// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"codeberg.org/gruf/go-mutexes"
	"github.com/superseriousbusiness/gotosocial/internal/cache"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/storage"
	"github.com/superseriousbusiness/gotosocial/internal/timeline"
	"github.com/superseriousbusiness/gotosocial/internal/workers"
)

// State provides a means of dependency injection and sharing of resources
// across different subpackages of the GoToSocial codebase. DO NOT assume
// that any particular field will be initialized if you are accessing this
// during initialization. A pointer to a State{} is often passed during
// subpackage initialization, while the returned subpackage type will later
// then be set and stored within the State{} itself.
type State struct {
	// Caches provides access to this state's collection of caches.
	Caches cache.Caches

	// Timelines provides access to this state's collection of timelines.
	Timelines timeline.Timelines

	// DB provides access to the database.
	DB db.DB

	// FedLocks provides access to this state's
	// mutex map of per URI federation locks.
	//
	// Used during account and status dereferencing,
	// message processing in the FromFediAPI worker
	// functions, and by the go-fed/activity library.
	FedLocks mutexes.MutexMap

	// AccountLocks provides access to this state's
	// mutex map of per URI locks, intended for use
	// when updating accounts, migrating, approving
	// or rejecting an account, changing stats,
	// pinned statuses, etc.
	AccountLocks mutexes.MutexMap

	// Storage provides access to the storage driver.
	Storage *storage.Driver

	// Workers provides access to this state's collection of worker pools.
	Workers workers.Workers

	// prevent pass-by-value.
	_ nocopy
}

// nocopy when embedded will signal linter to
// error on pass-by-value of parent struct.
type nocopy struct{}

func (*nocopy) Lock() {}

func (*nocopy) Unlock() {}
