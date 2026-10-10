'use strict';

const { install } = require('./install');

install().catch((err) => {
  console.warn(`linkbreakers-cli: could not download the linkbreakers binary during install (${err.message}).`);
  console.warn('linkbreakers-cli: it will be downloaded on first run instead.');
});
