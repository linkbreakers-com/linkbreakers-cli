'use strict';

const assert = require('node:assert/strict');
const test = require('node:test');
const zlib = require('node:zlib');
const { target, expectedChecksum, extractFromTarGz, extractFromZip } = require('./install');

function tarGz(files) {
  const blocks = [];
  for (const [name, content] of Object.entries(files)) {
    const header = Buffer.alloc(512);
    header.write(name, 0);
    header.write(content.length.toString(8).padStart(11, '0'), 124);
    header.write('0', 156);
    blocks.push(header, content, Buffer.alloc((512 - (content.length % 512)) % 512));
  }
  blocks.push(Buffer.alloc(1024));
  return zlib.gzipSync(Buffer.concat(blocks));
}

function zip(files) {
  const locals = [];
  const centrals = [];
  let offset = 0;
  for (const [name, content] of Object.entries(files)) {
    const data = zlib.deflateRawSync(content);
    const nameBuf = Buffer.from(name);
    const local = Buffer.alloc(30);
    local.writeUInt32LE(0x04034b50, 0);
    local.writeUInt16LE(8, 8);
    local.writeUInt32LE(data.length, 18);
    local.writeUInt32LE(content.length, 22);
    local.writeUInt16LE(nameBuf.length, 26);
    const central = Buffer.alloc(46);
    central.writeUInt32LE(0x02014b50, 0);
    central.writeUInt16LE(8, 10);
    central.writeUInt32LE(data.length, 20);
    central.writeUInt32LE(content.length, 24);
    central.writeUInt16LE(nameBuf.length, 28);
    central.writeUInt32LE(offset, 42);
    locals.push(local, nameBuf, data);
    centrals.push(central, nameBuf);
    offset += local.length + nameBuf.length + data.length;
  }
  const cd = Buffer.concat(centrals);
  const eocd = Buffer.alloc(22);
  eocd.writeUInt32LE(0x06054b50, 0);
  eocd.writeUInt16LE(Object.keys(files).length, 10);
  eocd.writeUInt32LE(cd.length, 12);
  eocd.writeUInt32LE(offset, 16);
  return Buffer.concat([...locals, cd, eocd]);
}

test('maps node platforms to release archive names', () => {
  assert.deepEqual(target('darwin', 'arm64'), { os: 'darwin', cpu: 'arm64', ext: 'tar.gz', exe: 'linkbreakers' });
  assert.deepEqual(target('linux', 'x64'), { os: 'linux', cpu: 'amd64', ext: 'tar.gz', exe: 'linkbreakers' });
  assert.deepEqual(target('win32', 'x64'), { os: 'windows', cpu: 'amd64', ext: 'zip', exe: 'linkbreakers.exe' });
  assert.throws(() => target('freebsd', 'x64'), /unsupported platform/);
  assert.throws(() => target('linux', 'ia32'), /unsupported platform/);
});

test('reads the expected checksum for an archive', () => {
  const checksums = 'aaa  linkbreakers-cli_1.0.0_linux_amd64.tar.gz\nbbb  linkbreakers-cli_1.0.0_darwin_arm64.tar.gz\n';
  assert.equal(expectedChecksum(checksums, 'linkbreakers-cli_1.0.0_darwin_arm64.tar.gz'), 'bbb');
  assert.throws(() => expectedChecksum(checksums, 'missing.zip'), /missing from checksums/);
});

test('extracts the binary from a tar.gz archive', () => {
  const archive = tarGz({ 'README.md': Buffer.from('readme'), linkbreakers: Buffer.from('binary-bytes') });
  assert.equal(extractFromTarGz(archive, 'linkbreakers').toString(), 'binary-bytes');
  assert.throws(() => extractFromTarGz(archive, 'other'), /not found/);
});

test('extracts the binary from a zip archive', () => {
  const archive = zip({ 'LICENSE': Buffer.from('mit'), 'linkbreakers.exe': Buffer.from('windows-binary') });
  assert.equal(extractFromZip(archive, 'linkbreakers.exe').toString(), 'windows-binary');
  assert.throws(() => extractFromZip(archive, 'other.exe'), /not found/);
});
