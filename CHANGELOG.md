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
