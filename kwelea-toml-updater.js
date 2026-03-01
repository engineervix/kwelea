// Custom commit-and-tag-version updater for kwelea.toml
// Handles the [site].version field and preserves the "v" prefix.

module.exports.readVersion = function (contents) {
  const match = contents.match(/^version\s*=\s*["']v?([^"']+)["']/m);
  return match ? match[1] : "0.0.0";
};

module.exports.writeVersion = function (contents, version) {
  return contents.replace(
    /^(version\s*=\s*["'])v?[^"']*["']/m,
    `$1v${version}"`
  );
};
