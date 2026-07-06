# Changelog

All notable changes to this project will be documented in this file. See [commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version) for commit guidelines.

## [0.1.6](https://github.com/engineervix/kwelea/compare/v0.1.5...v0.1.6) (2026-07-06)


### 🚀 Features

* add --output and --source CLI flags to override build dirs ([59b58bf](https://github.com/engineervix/kwelea/commit/59b58bfb178e79c5cef4b4e307dc3c06337456f7)), closes [#15](https://github.com/engineervix/kwelea/issues/15)
* add Open Graph and Twitter Card meta tags ([ac962fe](https://github.com/engineervix/kwelea/commit/ac962fe03e2aac6c47f4377bc0aa7173915e3454)), closes [#11](https://github.com/engineervix/kwelea/issues/11)
* generate sitemap.xml at build time ([ce2db31](https://github.com/engineervix/kwelea/commit/ce2db31765efb661b4b47b06114ba08edf59cb66)), closes [#4](https://github.com/engineervix/kwelea/issues/4)
* **parser:** add code block title and line highlighting ([9561183](https://github.com/engineervix/kwelea/commit/9561183fe95191ed19fec9f7c61b10aa7f6fe0f8))


### 🐛 Bug Fixes

* **deps:** update module github.com/alecthomas/chroma/v2 to v2.27.0 ([#32](https://github.com/engineervix/kwelea/issues/32)) ([e930256](https://github.com/engineervix/kwelea/commit/e930256bce55ff0ecad49c05ec65fe383473c0fc))
* **parser:** cap highlight ranges, drop dead code paths in code-attrs extension ([10722b6](https://github.com/engineervix/kwelea/commit/10722b624ab67847b880f6ee0ad5aa6b84ef4bce)), closes [#41](https://github.com/engineervix/kwelea/issues/41)
* **parser:** fix code-block highlight rendering and site-wide dark-mode colours ([4f1a492](https://github.com/engineervix/kwelea/commit/4f1a49276e56545b7e1f03ab888618d689741e72))
* reject empty values for --base-url, --output, --source ([d5cad54](https://github.com/engineervix/kwelea/commit/d5cad549a3f631f4a790279244f45dc423e1ab4f))
* **server:** reload kwelea.toml on change during serve ([cc85047](https://github.com/engineervix/kwelea/commit/cc850473fb2a355212bb52e4118b4f7a13e1d1df)), closes [#38](https://github.com/engineervix/kwelea/issues/38)
* tighten hasScheme to require non-empty host, update comments ([cd9f347](https://github.com/engineervix/kwelea/commit/cd9f34743b1e0ff8f79c5393ddbb41908acfed06)), closes [#35](https://github.com/engineervix/kwelea/issues/35)
* validate base_url has a scheme, add sitemap tests ([5defcdb](https://github.com/engineervix/kwelea/commit/5defcdbd59c05461f414e4ea40dd0a9cdb4a6df1)), closes [#35](https://github.com/engineervix/kwelea/issues/35)
* validate base_url has scheme, resolve relative og:image paths ([6574e39](https://github.com/engineervix/kwelea/commit/6574e394cc0daec0ac7f38560c1dd74fab61d9c5)), closes [#36](https://github.com/engineervix/kwelea/issues/36)


### 📝 Docs

* clarify which kwelea.toml values can be overridden via CLI ([f0bbd4b](https://github.com/engineervix/kwelea/commit/f0bbd4b072c01e5fa06e7119cf0d9776892ea059))
* document macOS Gatekeeper workaround for downloaded binaries ([e9488a2](https://github.com/engineervix/kwelea/commit/e9488a2c6153c2463801d8dee531a04110b0c9d5))


### ♻️ Code Refactoring

* extract applyFlagOverrides helper, drop package-level flag vars ([aedabc2](https://github.com/engineervix/kwelea/commit/aedabc296b4d6e91be8f8e2f6e1a2a240e67e907))


### ⚙️ Build System

* **deps:** update module github.com/alecthomas/chroma/v2 to v2.23.1 ([#22](https://github.com/engineervix/kwelea/issues/22)) ([d359305](https://github.com/engineervix/kwelea/commit/d3593055b360e33c464fd2f795c0d7f50c90908c))
* **deps:** update module github.com/fsnotify/fsnotify to v1.10.1 ([#33](https://github.com/engineervix/kwelea/issues/33)) ([5dee633](https://github.com/engineervix/kwelea/commit/5dee63308c0b35054f18c590e5596f83e0b1434e))
* **deps:** update module github.com/yuin/goldmark to v1.8.2 ([#23](https://github.com/engineervix/kwelea/issues/23)) ([2e7e406](https://github.com/engineervix/kwelea/commit/2e7e40635909dea1228f49058c41fd67ea9b5299))


### 👷 CI/CD

* **deps:** update actions/checkout action to v7 ([#34](https://github.com/engineervix/kwelea/issues/34)) ([d63eda7](https://github.com/engineervix/kwelea/commit/d63eda764a950911d7c26471d97e96a96915f05f))

## [0.1.5](https://github.com/engineervix/kwelea/compare/v0.1.4...v0.1.5) (2026-03-01)


### 🚀 Features

* add --base-url flag to kwelea build ([c701b4c](https://github.com/engineervix/kwelea/commit/c701b4c9e1856c530ec1673437caa3a56a26412c))


### 🐛 Bug Fixes

* ensure correct version is passed to the docs when we release ([6177b74](https://github.com/engineervix/kwelea/commit/6177b74dff4c12a1ed13926094472cd25ed85fec))


### 📝 Docs

* add contribution guide ([2c87ddc](https://github.com/engineervix/kwelea/commit/2c87ddc61fdc8175fc719fba93f16b92298e11b2))
* **readme:** show commits since last release ([9d43a34](https://github.com/engineervix/kwelea/commit/9d43a34b8c037b8461830ce1cb583c4e42521a6d))


### 👷 CI/CD

* add renovate config ([eb32e3d](https://github.com/engineervix/kwelea/commit/eb32e3de196fd81c36125027c35c9a54f4ba4bcc))

## [0.1.4](https://github.com/engineervix/kwelea/compare/v0.1.3...v0.1.4) (2026-03-01)


### 🚀 Features

* add site footer with extra_head/extra_footer injection hooks ([213c7c1](https://github.com/engineervix/kwelea/commit/213c7c12038644f014ff5050aa54b9e14e61b614)), closes [#10](https://github.com/engineervix/kwelea/issues/10)


### 🐛 Bug Fixes

* correct relative link to Installation page in Getting Started ([63962f6](https://github.com/engineervix/kwelea/commit/63962f6ab4434d8df9be7172918d3446fdaeaa9b))

## [0.1.3](https://github.com/engineervix/kwelea/compare/v0.1.2...v0.1.3) (2026-03-01)


### 🐛 Bug Fixes

* use relative font paths in theme.css; prefix SourceURL with BasePath ([f9473be](https://github.com/engineervix/kwelea/commit/f9473be059a0e6d69b1029e2722543cf15ea0877))

## [0.1.2](https://github.com/engineervix/kwelea/compare/v0.1.1...v0.1.2) (2026-03-01)


### 🐛 Bug Fixes

* build kwelea from source in CI; zero BasePath in dev server ([a5cb28f](https://github.com/engineervix/kwelea/commit/a5cb28fee532a7bc91410e2e90a69339dba78902))

## [0.1.1](https://github.com/engineervix/kwelea/compare/v0.1.0...v0.1.1) (2026-03-01)


### 🐛 Bug Fixes

* prefix all asset and page links with base path ([cef6898](https://github.com/engineervix/kwelea/commit/cef689816ab24a8537a80da3021fbe7410ed9243))


## [0.1.0](https://github.com/engineervix/kwelea/releases/tag/v0.1.0) (2026-02-28)


### 🚀 Features

* core parsing work ([aa3d68b](https://github.com/engineervix/kwelea/commit/aa3d68bc75e70b658857e1937d6e758cb8d3abf6))
* dev server ([dc0c5c7](https://github.com/engineervix/kwelea/commit/dc0c5c71a8b6392bd30416efa093361e6682b594))
* favicon, simple logo ([425b238](https://github.com/engineervix/kwelea/commit/425b2387ccc827118fa066b43310929a06efdb47))
* it ain't just github that's out there ([af8f9e9](https://github.com/engineervix/kwelea/commit/af8f9e9f1fc1a4c2a3d220db016908f3fb154d3b))
* nav tree ([c021725](https://github.com/engineervix/kwelea/commit/c021725840fe9979daa467bbf161607fc36d0e87))
* polish, distribution, and dogfooded docs ([f44b501](https://github.com/engineervix/kwelea/commit/f44b501719e781e670e4be7dd7c985d14d9d3fd8))
* search index ([2bc688f](https://github.com/engineervix/kwelea/commit/2bc688fd1603ce89a00cb890a6c8408f236ab79e))
* so we can also see the markdown source ([7b452cd](https://github.com/engineervix/kwelea/commit/7b452cdf266b54366aab6cc3843d3316a5727ded))
* templates & renderer ([cb566d1](https://github.com/engineervix/kwelea/commit/cb566d15b8d6fdf90e5e3c9017091377e242ac4a))
* this is kwelea - a fast, weaving documentation generator for Go ([cf6edea](https://github.com/engineervix/kwelea/commit/cf6edea2a491b198ba692c919105efc0c9f63e74))
* version info ([702129a](https://github.com/engineervix/kwelea/commit/702129ad43f8d6747534f4e215bc141814f858a4))


### 🐛 Bug Fixes

* colour contrast ([ce32815](https://github.com/engineervix/kwelea/commit/ce32815056081dfdd7384a6796d9346b6e66edbe))
* duplicate h1 issue ([3e3d290](https://github.com/engineervix/kwelea/commit/3e3d290f6a0f6f468fc29175d49d81d3548dd1bd))
* logging noise from D2 ([5eabcc4](https://github.com/engineervix/kwelea/commit/5eabcc497a402be75442aed7c5d13984621b4ec2))


### 📝 Docs

* installations options ([d6fe169](https://github.com/engineervix/kwelea/commit/d6fe1690661fc2908abebe146b49d15916b15e9a))
* organize the docs, cleanup, ensure consistency ([7f708df](https://github.com/engineervix/kwelea/commit/7f708dfe49740d1a92b8e192bcff6ead205dfcbe))


### ♻️ Code Refactoring

* **ui:** font size adjustments and more colour contrast improvements ([954a65a](https://github.com/engineervix/kwelea/commit/954a65a1eb68fcdbd82639b3e4867da5c64a9089))


### 👷 CI/CD

* setup github actions ([7db5274](https://github.com/engineervix/kwelea/commit/7db52748e881f99cfa4f4ba7459991ee2f1b8d00))
