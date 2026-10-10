#!/usr/bin/env node
'use strict';

const { spawnSync } = require('node:child_process');
const { ensureBinary } = require('../lib/install');

ensureBinary()
  .then((bin) => {
    const result = spawnSync(bin, process.argv.slice(2), { stdio: 'inherit', windowsHide: true });
    if (result.error) throw result.error;
    if (result.signal) process.kill(process.pid, result.signal);
    process.exit(result.status ?? 1);
  })
  .catch((err) => {
    console.error(`linkbreakers: ${err.message}`);
    process.exit(1);
  });
