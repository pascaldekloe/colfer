var testrunner = require("qunit");

testrunner.setup({
	log: {
		errors: true,
		tests: true,
		coverage: true
	},
	coverage: true,
	maxBlockDuration: 2000
});

testrunner.run({
	code: "./Colfer.js",
	tests: "./test.js"
});
