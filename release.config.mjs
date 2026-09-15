/**
 * @type {import('semantic-release').GlobalConfig}
 */
export default {
  
  branches: [
    "release",
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
   [
      "@semantic-release/github",
      {
        assets: [
          // Linux
          { path: "dist/linux-amd64/nagare", name: "nagare-linux-amd64-${nextRelease.version}" },
          { path: "dist/linux-amd64/nagare-gateway", name: "nagare-gateway-linux-amd64-${nextRelease.version}" },
          { path: "dist/linux-arm64/nagare", name: "nagare-linux-arm64-${nextRelease.version}" },
          { path: "dist/linux-arm64/nagare-gateway", name: "nagare-gateway-linux-arm64-${nextRelease.version}" },

          // Windows
          { path: "dist/windows-amd64/nagare", name: "nagare-windows-amd64-${nextRelease.version}.exe" },
          { path: "dist/windows-amd64/nagare-gateway", name: "nagare-gateway-windows-amd64-${nextRelease.version}.exe" },
          { path: "dist/windows-arm64/nagare", name: "nagare-windows-arm64-${nextRelease.version}.exe" },
          { path: "dist/windows-arm64/nagare-gateway", name: "nagare-gateway-windows-arm64-${nextRelease.version}.exe" },
          
          // macOS (Darwin)
          { path: "dist/darwin-amd64/nagare", name: "nagare-darwin-amd64-${nextRelease.version}" },
          { path: "dist/darwin-amd64/nagare-gateway", name: "nagare-gateway-darwin-amd64-${nextRelease.version}" },
          { path: "dist/darwin-arm64/nagare", name: "nagare-darwin-arm64-${nextRelease.version}" },
          { path: "dist/darwin-arm64/nagare-gateway", name: "nagare-gateway-darwin-arm64-${nextRelease.version}" }
        ]
      }
    ],
    "@semantic-release/git",
  ],
};