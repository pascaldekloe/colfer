global._ = require('./lodash.js')
const Benchmark = require('./benchmark.js');

const Colfer = require('./build/Colfer.js')

const testData = [
	new Colfer.bench.Colfer({key: 1234567890, host: "db003lz12", port: 389, size: 452, hash: 0x5c2428488918, ratio: 0.99, route: true}),
	new Colfer.bench.Colfer({key: 1234567891, host: "localhost", port: 22, size: 4096, hash: 0x48899c24c824, ratio: 0.20, route: false}),
	new Colfer.bench.Colfer({key: 1234567892, host: "kdc.local", port: 88, size: 1984, hash: 0x48891c24485c, ratio: 0.06, route: false}),
	new Colfer.bench.Colfer({key: 1234567893, host: "vhost8.dmz.example.com", port: 27017, size: 59741, hash: 0x08488b9c2489, ratio: 0.0, route: true}),
];

// Corresponding testData Colfer serials.
var testColfer = new Array(testData.length);
testData.forEach(function(o, i) {
	testColfer[i] = o.marshal();
});

// Corresponding testData JSON serials.
var testJSON = new Array(testData.length);
testData.forEach(function(o, i) {
	testJSON[i] = JSON.stringify(o);
});

var suite = new Benchmark.Suite;
suite.add('marshal Colfer', function() {
		testData[0].marshal();
		testData[1].marshal();
		testData[2].marshal();
		testData[3].marshal();
	})
	.add('unmarshal Colfer', function() {
		new Colfer.bench.Colfer().unmarshal(testColfer[0]);
		new Colfer.bench.Colfer().unmarshal(testColfer[1]);
		new Colfer.bench.Colfer().unmarshal(testColfer[2]);
		new Colfer.bench.Colfer().unmarshal(testColfer[3]);
	})
	.add('marshal JSON', function() {
		JSON.stringify(testData[0]);
		JSON.stringify(testData[1]);
		JSON.stringify(testData[2]);
		JSON.stringify(testData[3]);
	})
	.add('unmarshal JSON', function() {
		JSON.parse(testJSON[0]);
		JSON.parse(testJSON[1]);
		JSON.parse(testJSON[2]);
		JSON.parse(testJSON[3]);
	})
	.on('error', function(event) {
		console.log(String(event.target.error));
	})
	.on('cycle', function(event) {
		console.log(String(event.target))
	})
	.run({async: true});
