use log::{error, info};
use nagare_plugin::PluginClient;

const PLUGIN_METADATA: &str = include_str!("../metadata.json");

#[tokio::main]
async fn main() {
    let plugin = PluginClient::new_with_logger();
    match plugin.load_metadata(PLUGIN_METADATA) {
        Ok(_) => {
            info!("Plugin metadata loaded")
        }
        Err(e) => {
            error!("Failed to load plugin metadata {e}");
        }
    }

    plugin
        .start(|| async {
            info!("starting plugin ....");
            match plugin.handshake().await {
                Ok(_) => {
                    info!("hanshake success")
                }
                Err(e) => {
                    error!("Failed to handshake: {e}");
                }
            }
        })
        .await
        .unwrap();
}
