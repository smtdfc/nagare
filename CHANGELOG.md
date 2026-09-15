# [1.4.0](https://github.com/smtdfc/nagare/compare/v1.3.0...v1.4.0) (2026-09-15)


### Bug Fixes

* **core/plugin:** Potential fix for pull request finding 'CodeQL / Arbitrary file access during archive extraction ("Zip Slip")' ([3303792](https://github.com/smtdfc/nagare/commit/3303792f08b6fae3741d3b09e8386fc28f74093f))
* correct formatting of GitHub plugin configuration in release.config.mjs ([6796738](https://github.com/smtdfc/nagare/commit/679673898be44ead3499aa6a48f0ce95e34e2295))
* ensure dix is installed before building for platforms ([98adec8](https://github.com/smtdfc/nagare/commit/98adec86b557329bd107b89a57dd282221a0fcee))
* **plugin/client:** improve handshake error handling and remove unused functions ([68707f9](https://github.com/smtdfc/nagare/commit/68707f9b2ba543a32b6a39d60ad6cffca3fb808b))
* update dix installation command to use the correct repository ([da991b0](https://github.com/smtdfc/nagare/commit/da991b07836840c1f7f4cc8ca10b4e3b1854e92b))
* update Go version in release workflow to 1.27 ([154a3b6](https://github.com/smtdfc/nagare/commit/154a3b6a32c4078ac828c21f40b0e08715636e44))
* update release branch configuration to remove prerelease setting ([f0a6cc7](https://github.com/smtdfc/nagare/commit/f0a6cc7c3e61bb6efcb048caab8dfa0c49b077ae))
* update release workflow to use pnpm for running Semantic Release ([a2e0fac](https://github.com/smtdfc/nagare/commit/a2e0faca8e43292f2f404e0edc9da7884436b3e0))


### Features

* add chat event payloads and websocket events for plugin communication ([77d68bf](https://github.com/smtdfc/nagare/commit/77d68bf2ced905d5170b45442053ca46064004cb))
* add configuration files for chat and llm_provider modules; include environment variables and websocket settings ([d8a7f0e](https://github.com/smtdfc/nagare/commit/d8a7f0e60460d642a9c1bd7dbf2423e9efee5421))
* add go.sum file and update dependencies for plugin client ([2e9520d](https://github.com/smtdfc/nagare/commit/2e9520de2b9ee155c1b65de8b10af95fc694185e))
* add GROG_API_KEY to Development.yml and update llm_provider configurations for improved API integration ([2c0b34d](https://github.com/smtdfc/nagare/commit/2c0b34da087e933df58608ab28f1aac6ac0bcfdd))
* add initial DTOs and REST API structures for chat and plugin management ([8363f90](https://github.com/smtdfc/nagare/commit/8363f90d8e5c97894ff23100a6c6091321f9bf8c))
* add plugin builder and client structure with metadata handling ([74e0196](https://github.com/smtdfc/nagare/commit/74e0196c167b8d2c8b555584f3676c39b3209b60))
* add Type field to message structs for improved message type handling ([0d13159](https://github.com/smtdfc/nagare/commit/0d13159f68d1fb37ea0e22fc6328b896ad18b209))
* add websocket event payloads for chat message handling ([552c25e](https://github.com/smtdfc/nagare/commit/552c25e8ab913aaaa2544b93ef39bcff21c4a402))
* add YAML configuration files for plugin management including folder, install-local, and list ([c469c78](https://github.com/smtdfc/nagare/commit/c469c78db09486ee645b3ffb136a09b23e51a2b8))
* **core/plugin:**  add hostPort to PluginManager and update StartPlugin to include host port in environment variables ([3babcda](https://github.com/smtdfc/nagare/commit/3babcdadbf135a837f9e69376e629f9bb880bb97))
* **core/setup:** update Setup method to accept port and set plugin host port ([cc18b17](https://github.com/smtdfc/nagare/commit/cc18b1724bdc17138958fa1828b7a4dbca2f3d9c))
* enhance build and installation scripts for multi-platform support ([4c7c724](https://github.com/smtdfc/nagare/commit/4c7c72427867f27590f166e94f7ec09aae386ed2))
* enhance plugin connection handling and add weather tool implementation ([9f8c8d0](https://github.com/smtdfc/nagare/commit/9f8c8d021555be62354e3171278470536c6da208))
* **gateway:** enhance chat and plugin functionality with new auth middleware and Telegram plugin integration ([242185c](https://github.com/smtdfc/nagare/commit/242185c36a76b2eeb5e5d5ce64194f54b9f802ee))
* **gateway:** enhance chat and plugin functionality with RequestID and graceful shutdown support ([5bf234d](https://github.com/smtdfc/nagare/commit/5bf234d336c17aee9afea8b10af238079c008a2b))
* implement chat history saving functionality and add batch creation for messages ([d6058c9](https://github.com/smtdfc/nagare/commit/d6058c9d4f18d6eda088f4257f83dc056760bca0))
* **plugin/chat:** add chat message event payloads and update event constants ([f333ff7](https://github.com/smtdfc/nagare/commit/f333ff7dcc2b1da904fa8ce9a7ebe4603112fb1c))
* **plugin/chat:** implement chat session preparation and message sending functionality ([3bad8d2](https://github.com/smtdfc/nagare/commit/3bad8d261e298b48e7e07953c76e5538120c5891))
* **plugin/client:** implement PluginClient with connection handling and handshake support ([1141c91](https://github.com/smtdfc/nagare/commit/1141c91c037607b4422433017e40b9b62c995713))
* **plugin/example:** enhance main function with metadata loading and handshake support ([f40ae39](https://github.com/smtdfc/nagare/commit/f40ae39015d65161c247725a1a3d28bacf14bd30))
* **plugin/manager:** enhance unpackPlugin to validate destination directory and file paths ([7cc6e27](https://github.com/smtdfc/nagare/commit/7cc6e2772ee6d16924ae7d5aa377cb80d84f2436))
* **plugin/telegram:** add reset command ([122ef02](https://github.com/smtdfc/nagare/commit/122ef020547cae95ca46358638ecb2fa57333087))
* **plugin/telegram:** add session management for chat interactions ([4edec9c](https://github.com/smtdfc/nagare/commit/4edec9cf07aa496b4085c4fa24f2e2c6920ad2f5))
* **plugin/telegram:** add Telegram plugin with metadata, go.mod, and go.sum files ([2f62df9](https://github.com/smtdfc/nagare/commit/2f62df922d08c47d81b983f52788e3201ebcb730))
* **plugin/telegram:** implement message buffering and processing for chat sessions ([b5bbef0](https://github.com/smtdfc/nagare/commit/b5bbef0537c6afb31bd471176d14b764ae94e3c2))
* replace hello message with plugin client initialization and start ([4cc9260](https://github.com/smtdfc/nagare/commit/4cc9260c4597a9a7ebfd5522e2a8845feecf48f7))
* **session:** enhance session management with plugin support and error handling improvements ([f901bb5](https://github.com/smtdfc/nagare/commit/f901bb59f4ed57e9d931947b8768aa926daf1a1c))
* update .gitignore and add configuration files for JSON schema and Prettier ([82b5362](https://github.com/smtdfc/nagare/commit/82b536214a5a1682c1075f4902c668253fe7eb28))

# [1.3.0](https://github.com/smtdfc/nagare/compare/v1.2.0...v1.3.0) (2026-08-20)


### Bug Fixes

* **core/plugin:** fix plugin stdout/stderr streams not displaying ([4fe1755](https://github.com/smtdfc/nagare/commit/4fe1755c19013a7576742437f07070af16410b39))
* **core/session:** fix incorrect return type ([ba51804](https://github.com/smtdfc/nagare/commit/ba51804a9bde4b02cb80db4a9e70815e7b7d004f))
* **core:** fix missing condition in database query ([03895f2](https://github.com/smtdfc/nagare/commit/03895f23d72df2f9084ac8b0ec62227c87583691))
* **core:** remove unique constraint for soft-deleted plugins ([d16769c](https://github.com/smtdfc/nagare/commit/d16769c41372a18159c5dc07ce396846118db6a5))
* **server/app:** resolve missing dependency injection provider ([6ce1ae5](https://github.com/smtdfc/nagare/commit/6ce1ae58c18c664087cb800198fc626dba46f2df))
* **ui/main:** improve error toast display behavior ([1589cd3](https://github.com/smtdfc/nagare/commit/1589cd3541745e6c6eb6a0c5b6c2e7614daabeec))


### Features

* **core/logger:** add multi-writer support to logger ([06477f0](https://github.com/smtdfc/nagare/commit/06477f08d69c8de1cb3ebb630c3c89231158cf0d))
* **core/session:** implement GetOrCreateSessionByUserID method ([0675a82](https://github.com/smtdfc/nagare/commit/0675a8276f99b217b83f9d0c6a794b19d0b47ea9))
* enable plugin removal functionality ([5252b64](https://github.com/smtdfc/nagare/commit/5252b647feb443868d8f8a1e16f38e3b59ebab38))
* **plugin/client:** add context support to control and cancel plugin ([29e572a](https://github.com/smtdfc/nagare/commit/29e572ae7b90bbb23dae4c5496b86fb88c7cf770))
* **plugin/client:** add websocket event listener and unlistener helpers ([5b9b674](https://github.com/smtdfc/nagare/commit/5b9b6740e5ba8c657027017d9d49d4b870fb67bd))
* **plugin:** implement agent invocation over websocket and telegram ([4f3c3a6](https://github.com/smtdfc/nagare/commit/4f3c3a65a25020cfc45e8c80897c84ecdc6d0eab))
* **plugin:** introduce Nagare Telegram plugin ([c13484f](https://github.com/smtdfc/nagare/commit/c13484f3435c28b1ced938e60c3deb3dccdcc14a))
* **server:** add plugin event bus for internal messaging ([3243bc5](https://github.com/smtdfc/nagare/commit/3243bc5431c6536a347e5ffb8b0d46a121ed4cfd))
* **server:** add startup banner ([3f72770](https://github.com/smtdfc/nagare/commit/3f727702807d3dc363cfa5bbdd9c409915b9a1fd))
* **server:** introduce experimental ReduceMemoryUsage configuration ([d6bfb6a](https://github.com/smtdfc/nagare/commit/d6bfb6af3db0e796d169161dd3f0d76b12a5f120))
* **shared/bus:** implement core event bus ([0881a58](https://github.com/smtdfc/nagare/commit/0881a58cf9936563e0a829c8f7b85ae9db6efc46))
* **shared/paths:** add plugin configuration directory path ([61caefc](https://github.com/smtdfc/nagare/commit/61caefc5f18b65ee023d36592d0370b0a5494349))

# [1.2.0](https://github.com/smtdfc/nagare/compare/v1.1.0...v1.2.0) (2026-08-18)


### Features

* **core:** implement plugin connection code management(Triển khai ([e76ec56](https://github.com/smtdfc/nagare/commit/e76ec56016ce241dbdbba9cabd0210c970bf8375))
* **plugin:** establish bidirectional websocket communication with ([87cb078](https://github.com/smtdfc/nagare/commit/87cb078ff7248c21b9a5231dfe8481b0620afa65))
* **server:** add plugin handler for websocket connections ([b4f2f92](https://github.com/smtdfc/nagare/commit/b4f2f92a4de4088bd8419fec77b6dc5a99574f5f))
* **server:** implement routing for plugin connections ([ce6ac7d](https://github.com/smtdfc/nagare/commit/ce6ac7d96b1500418d8fe4c48b5c7270cfdd6d96))
* **server:** integrate new websocket handler ([91119d2](https://github.com/smtdfc/nagare/commit/91119d2770139198613487271b77783c166c0ba4))
* **shared:** add DTOs for plugin registration events ([05ac1f5](https://github.com/smtdfc/nagare/commit/05ac1f5c66ee658f02be3d51ef7b15206e800887))
* **shared:** add GetPluginWebsocketConnect helper ([87d462d](https://github.com/smtdfc/nagare/commit/87d462d2b9ecc72b12ddddabc44d5649671c79f8))
* **shared:** add websocket helpers ([99333da](https://github.com/smtdfc/nagare/commit/99333dac1f1aad43a304859c9a6b72fc908e61f0))
* **shared:** extend AuthPayload with target property ([222c92b](https://github.com/smtdfc/nagare/commit/222c92bb4e14c00f056d05c20a5b7d6da2a9e6bd))

# [1.1.0](https://github.com/smtdfc/nagare/compare/v1.0.2...v1.1.0) (2026-08-17)


### Features

* **core:** add DeletePluginByID method in PluginRepository ([f2082db](https://github.com/smtdfc/nagare/commit/f2082db6e834a9f1ab0ac18a92e75830724ed449))
* **core:** add log when plugin start ([4a66b40](https://github.com/smtdfc/nagare/commit/4a66b40cf10ca93469860cd920ad0d6edc23a2a2))
* **core:** add plugin logs ([ecf7f2e](https://github.com/smtdfc/nagare/commit/ecf7f2ef001d39652419ce860107233cdcb6b170))
* **core:** allow stop plugin process ([1870828](https://github.com/smtdfc/nagare/commit/187082884c535165c80c60a106839226e018acbf))
* **core:** save pid when start plugin ([5711666](https://github.com/smtdfc/nagare/commit/5711666cd8b9f01c8f309f8ccdc7ff4ad5f11bed))
* **ui:** add breadcrumb navigation for pages ([913d224](https://github.com/smtdfc/nagare/commit/913d2249ce57b12cd5ab51d8d5598b075c58e8a7))
* **ui:** add plugin installation menu item to sidebar ([5092041](https://github.com/smtdfc/nagare/commit/50920417e0f1b7ca7f2d6506f3e884ce083c2de6))
* **ui:** add sidebar trigger button ([ba95e2d](https://github.com/smtdfc/nagare/commit/ba95e2d7e799fcc542501561c9d2cf76fbf23da6))
