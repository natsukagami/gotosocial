/*
	GoToSocial
	Copyright (C) GoToSocial Authors admin@gotosocial.org
	SPDX-License-Identifier: AGPL-3.0-or-later

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

"use strict";

const React = require("react");
const { Switch, Route } = require("wouter");

const UserDetail = require("./detail");

module.exports = function Users({ baseUrl }) {
	return (
		<div className="users">
			<Switch>
				<Route path={`${baseUrl}/:userId`}>
					<UserDetail />
				</Route>
				<UserOverview />
			</Switch>
		</div>
	);
};

function UserOverview({ }) {
	return (
		<>
			<h1>Users</h1>
			<div>
				Pending <a href="https://github.com/superseriousbusiness/gotosocial/issues/582">#582</a> and <a href="https://github.com/superseriousbusiness/gotosocial/issues/581">#581</a>,
				there is currently no way to list user accounts.<br />
				You can perform actions on reported users by clicking their name.
			</div>
		</>
	);
}
