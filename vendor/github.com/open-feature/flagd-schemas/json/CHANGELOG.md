# Changelog

## [0.2.10](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.9...json/json-schema-v0.2.10) (2025-01-06)


### ✨ New Features

* rename named top-level prop from id to flagSetId ([#177](https://github.com/open-feature/flagd-schemas/issues/177)) ([bb4d9fd](https://github.com/open-feature/flagd-schemas/commit/bb4d9fd6dc4f0b8a1af77335404c949697b766e2))

## [0.2.9](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.8...json/json-schema-v0.2.9) (2024-12-05)


### ✨ New Features

* add flag set and flag metadata to schema ([#173](https://github.com/open-feature/flagd-schemas/issues/173)) ([9fe3677](https://github.com/open-feature/flagd-schemas/commit/9fe36777df7f2697d78830125f04f6972cf86e7a))


### 🧹 Chore

* remove comment about summing to 100 ([#172](https://github.com/open-feature/flagd-schemas/issues/172)) ([d03b772](https://github.com/open-feature/flagd-schemas/commit/d03b772f167aa21a635e70d22641d792deefcb7f))

## [0.2.8](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.7...json/json-schema-v0.2.8) (2024-07-08)


### ✨ New Features

* standalone targeting, fractional shorthand ([#168](https://github.com/open-feature/flagd-schemas/issues/168)) ([6211e3a](https://github.com/open-feature/flagd-schemas/commit/6211e3a5ea413809c2818fbe54ba1b6769ac0f54))

## [0.2.7](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.6...json/json-schema-v0.2.7) (2024-06-17)


### 🐛 Bug Fixes

* revert full schema versioning ([#166](https://github.com/open-feature/flagd-schemas/issues/166)) ([64b6d98](https://github.com/open-feature/flagd-schemas/commit/64b6d9831bc40c4b78df78f7f77f2d482632903f))

## [0.2.6](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.5...json/json-schema-v0.2.6) (2024-05-27)


### ✨ New Features

* fully version schema ([#163](https://github.com/open-feature/flagd-schemas/issues/163)) ([c0692e4](https://github.com/open-feature/flagd-schemas/commit/c0692e484cfa969ad19077f10b153752ce017676))

## [0.2.5](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.4...json/json-schema-v0.2.5) (2024-05-23)


### 🐛 Bug Fixes

* unary op shorthand unsupported ([#142](https://github.com/open-feature/flagd-schemas/issues/142)) ([946ffc4](https://github.com/open-feature/flagd-schemas/commit/946ffc4d2ed29e91d3beddd7b0650df350411c71))

## [0.2.4](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.3...json/json-schema-v0.2.4) (2024-05-17)


### 🐛 Bug Fixes

* change $defs to definitions per draft07 ([#140](https://github.com/open-feature/flagd-schemas/issues/140)) ([87ab28a](https://github.com/open-feature/flagd-schemas/commit/87ab28a58cac260b911c03b35451c8175ec1f148))

## [0.2.3](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.2...json/json-schema-v0.2.3) (2024-04-08)


### ✨ New Features

* allow custom seed in fractional op ([#136](https://github.com/open-feature/flagd-schemas/issues/136)) ([ea4f119](https://github.com/open-feature/flagd-schemas/commit/ea4f119d2bd716ec4aec05e554d51e7e79ba187b))
* support chainable if conditions ([#138](https://github.com/open-feature/flagd-schemas/issues/138)) ([cc7832a](https://github.com/open-feature/flagd-schemas/commit/cc7832ab20c9e0b8e438ffc4299f661974149454))

## [0.2.2](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.1...json/json-schema-v0.2.2) (2024-02-28)


### 🐛 Bug Fixes

* workaround bad $ref derefs ([#134](https://github.com/open-feature/flagd-schemas/issues/134)) ([40a0677](https://github.com/open-feature/flagd-schemas/commit/40a0677a6a97b9e15d4c8e9419d4b666bfa778b3))

## [0.2.1](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.2.0...json/json-schema-v0.2.1) (2024-02-14)


### 🐛 Bug Fixes

* support "if" 1-arg shorthand ([#131](https://github.com/open-feature/flagd-schemas/issues/131)) ([463d33a](https://github.com/open-feature/flagd-schemas/commit/463d33a3895f1cd1149c9d99cdc5fb7981abd296))

## [0.2.0](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.1.3...json/json-schema-v0.2.0) (2024-01-18)


### ⚠ BREAKING CHANGES

* improve naming, fix module url ([#128](https://github.com/open-feature/flagd-schemas/issues/128))

### 🐛 Bug Fixes

* "if" json schema description ([#126](https://github.com/open-feature/flagd-schemas/issues/126)) ([eee6a28](https://github.com/open-feature/flagd-schemas/commit/eee6a2810d9c45a360841a30f3fb92b534f5611d))


### ✨ New Features

* improve naming, fix module url ([#128](https://github.com/open-feature/flagd-schemas/issues/128)) ([0e6cbbd](https://github.com/open-feature/flagd-schemas/commit/0e6cbbd89d1591df728c9ab06d6cf4065f432dfe))

## [0.1.3](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.1.2...json/json-schema-v0.1.3) (2024-01-15)


### 🐛 Bug Fixes

* add "/" op and fix special cases in schema ([#122](https://github.com/open-feature/flagd-schemas/issues/122)) ([4a0be4f](https://github.com/open-feature/flagd-schemas/commit/4a0be4f48816ea0ac83d909ae58b8dcf5acda4b8))

## [0.1.2](https://github.com/open-feature/flagd-schemas/compare/json/json-schema-v0.1.1...json/json-schema-v0.1.2) (2024-01-08)


### ✨ New Features

* add targeting schema, improve def. schema ([#120](https://github.com/open-feature/flagd-schemas/issues/120)) ([6041fc7](https://github.com/open-feature/flagd-schemas/commit/6041fc7ef05fdd6ea9013718f253c869cb528b68))

## [0.1.1](https://github.com/open-feature/schemas/compare/json/json-schema-v0.1.0...json/json-schema-v0.1.1) (2023-01-26)


### Bug Fixes

* add ajv schema linting to pr workflow and validation fixes ([#71](https://github.com/open-feature/schemas/issues/71)) ([407a39c](https://github.com/open-feature/schemas/commit/407a39c2049e95ae1d80c28b68aa2658d597fbc5))
