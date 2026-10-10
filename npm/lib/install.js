'use strict';

const crypto = require('node:crypto');
const fs = require('node:fs');
const path = require('node:path');
const zlib = require('node:zlib');

const REPO = 'linkbreakers-com/linkbreakers-cli';
const VENDOR_DIR = path.join(__dirname, '..', 'vendor');

function target(platform = process.platform, arch = process.arch) {
  const os = { darwin: 'darwin', linux: 'linux', win32: 'windows' }[platform];
  const cpu = { x64: 'amd64', arm64: 'arm64' }[arch];
  if (!os || !cpu) {
    throw new Error(`unsupported platform ${platform}/${arch}; download a binary from https://github.com/${REPO}/releases`);
  }
  return { os, cpu, ext: os === 'windows' ? 'zip' : 'tar.gz', exe: os === 'windows' ? 'linkbreakers.exe' : 'linkbreakers' };
}

function binaryPath() {
  return path.join(VENDOR_DIR, target().exe);
}

function packageVersion() {
  return require('../package.json').version;
}

async function download(url) {
  const res = await fetch(url, { headers: { 'User-Agent': 'linkbreakers-cli-npm' } });
  if (!res.ok) throw new Error(`GET ${url} failed: ${res.status} ${res.statusText}`);
  return Buffer.from(await res.arrayBuffer());
}

function expectedChecksum(checksums, archiveName) {
  for (const line of checksums.split('\n')) {
    const [sum, name] = line.trim().split(/\s+/);
    if (name === archiveName) return sum;
  }
  throw new Error(`${archiveName} is missing from checksums.txt`);
}

function extractFromTarGz(archive, fileName) {
  const tar = zlib.gunzipSync(archive);
  let offset = 0;
  while (offset + 512 <= tar.length) {
    const header = tar.subarray(offset, offset + 512);
    if (header.every((b) => b === 0)) break;
    const name = header.toString('utf8', 0, 100).replace(/\0.*$/s, '');
    const prefix = header.toString('utf8', 345, 500).replace(/\0.*$/s, '');
    const size = parseInt(header.toString('utf8', 124, 136).replace(/\0.*$/s, '').trim() || '0', 8);
    const type = String.fromCharCode(header[156]);
    const fullName = prefix ? `${prefix}/${name}` : name;
    offset += 512;
    if ((type === '0' || type === '\0') && path.posix.basename(fullName) === fileName) {
      return tar.subarray(offset, offset + size);
    }
    offset += Math.ceil(size / 512) * 512;
  }
  throw new Error(`${fileName} not found in archive`);
}

function extractFromZip(archive, fileName) {
  const eocd = archive.lastIndexOf(Buffer.from([0x50, 0x4b, 0x05, 0x06]));
  if (eocd < 0) throw new Error('invalid zip archive');
  const entries = archive.readUInt16LE(eocd + 10);
  let offset = archive.readUInt32LE(eocd + 16);
  for (let i = 0; i < entries; i++) {
    if (archive.readUInt32LE(offset) !== 0x02014b50) throw new Error('invalid zip central directory');
    const method = archive.readUInt16LE(offset + 10);
    const compressedSize = archive.readUInt32LE(offset + 20);
    const nameLength = archive.readUInt16LE(offset + 28);
    const extraLength = archive.readUInt16LE(offset + 30);
    const commentLength = archive.readUInt16LE(offset + 32);
    const localOffset = archive.readUInt32LE(offset + 42);
    const name = archive.toString('utf8', offset + 46, offset + 46 + nameLength);
    offset += 46 + nameLength + extraLength + commentLength;
    if (path.posix.basename(name) !== fileName) continue;

    const dataStart = localOffset + 30 + archive.readUInt16LE(localOffset + 26) + archive.readUInt16LE(localOffset + 28);
    const data = archive.subarray(dataStart, dataStart + compressedSize);
    if (method === 0) return data;
    if (method === 8) return zlib.inflateRawSync(data);
    throw new Error(`unsupported zip compression method ${method}`);
  }
  throw new Error(`${fileName} not found in archive`);
}

async function install() {
  const version = packageVersion();
  if (!version || version === '0.0.0') {
    throw new Error('this package has no release version; install linkbreakers-cli from npm');
  }
  const t = target();
  const archiveName = `linkbreakers-cli_${version}_${t.os}_${t.cpu}.${t.ext}`;
  const base = `https://github.com/${REPO}/releases/download/v${version}`;

  const [archive, checksums] = await Promise.all([
    download(`${base}/${archiveName}`),
    download(`${base}/checksums.txt`).then((b) => b.toString('utf8')),
  ]);

  const actual = crypto.createHash('sha256').update(archive).digest('hex');
  const expected = expectedChecksum(checksums, archiveName);
  if (actual !== expected) {
    throw new Error(`checksum mismatch for ${archiveName}: expected ${expected}, got ${actual}`);
  }

  const binary = t.ext === 'zip' ? extractFromZip(archive, t.exe) : extractFromTarGz(archive, t.exe);
  fs.mkdirSync(VENDOR_DIR, { recursive: true });
  const dest = path.join(VENDOR_DIR, t.exe);
  const tmp = `${dest}.${process.pid}.tmp`;
  fs.writeFileSync(tmp, binary, { mode: 0o755 });
  fs.renameSync(tmp, dest);
  return dest;
}

async function ensureBinary() {
  const bin = binaryPath();
  if (fs.existsSync(bin)) return bin;
  return install();
}

module.exports = { ensureBinary, install, target, expectedChecksum, extractFromTarGz, extractFromZip };
