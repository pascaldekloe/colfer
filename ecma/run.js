var testrunner = require("qunit");

testrunner.setup({
	log: {
		errors: true,
		tests: true,
		coverage: true
	},
	coverage: {
		dir: "build/coverage"
	},
	maxBlockDuration: 2000
});

testrunner.run({
	code: "gen/Colfer.js",
	tests: "./test.js"
});
