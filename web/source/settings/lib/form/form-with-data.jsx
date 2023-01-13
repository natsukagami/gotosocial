"use strict";

const React = require("react");

const Loading = require("../../components/loading");

// Wrap Form component inside component that fires the RTK Query call,
// so Form will only be rendered when data is available to generate form-fields for
module.exports = function FormWithData({ dataQuery, DataForm, queryArg, ...formProps }) {
	const { data, isLoading } = dataQuery(queryArg);

	if (isLoading) {
		return (
			<div>
				<Loading />
			</div>
		);
	} else {
		return <DataForm data={data} {...formProps} />;
	}
};