var addon = require('bindings')('hello');

console.log(addon.createInt32()); // 'world'