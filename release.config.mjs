/**
 * @type {import('semantic-release').GlobalConfig}
 */
export default {
  branches: [
    { name: "release", prerelease: "rc" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    "@semantic-release/github",
    [
      "@semantic-release/exec",
      {
        prepareCmd: "./build.sh",
      },
    ],
    "@semantic-release/github",
    {
      assets: [
        { path: "dist/cli/nagare", name: "nagare-${nextRelease.version}" },
        { path: "dist/gateway/nagare-gateway", name: "nagare-gateway-${nextRelease.version}" }
      ]
    },
    "@semantic-release/git",
  ],
};