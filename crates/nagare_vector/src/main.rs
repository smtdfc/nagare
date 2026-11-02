use std::{
    path::{Path, PathBuf},
    sync::Arc,
};

use log::{error, info};
use nagare_module::ModuleClient;

use crate::vector_index::VectorIndex;

mod event;
mod vector_index;

const INDEX_DIM: usize = 1024;
const INDEX_BITS: usize = 4;
const INDEX_FILE_NAME: &str = "vector_index.tvim";

struct VectorModule {
    client: Arc<ModuleClient>,
    vector_index: VectorIndex,
    data_path: PathBuf,
}

impl VectorModule {
    pub fn new() -> Result<Self, String> {
        let data_path = PathBuf::from(
            std::env::var("NAGARE_MODULE_DATA_DIR").unwrap_or_else(|_| "./data".to_string()),
        );

        std::fs::create_dir_all(&data_path)
            .map_err(|e| format!("Failed to create data directory: {e}"))?;

        let vector_index = Self::ensure_index_exists(&data_path)?;

        Ok(Self {
            data_path,
            client: ModuleClient::new_with_logger(),
            vector_index,
        })
    }

    fn ensure_index_exists(data_path: &Path) -> Result<VectorIndex, String> {
        let vector_file = data_path.join(INDEX_FILE_NAME);

        if vector_file.exists() {
            VectorIndex::load(&vector_file)
        } else {
            VectorIndex::new(INDEX_DIM, INDEX_BITS)
        }
    }

    pub async fn start(&self) {
        let result = self
            .client
            .start(|| async {
                info!("starting module ....");
                match self.client.handshake().await {
                    Ok(_) => {
                        info!("hanshake success")
                    }
                    Err(e) => {
                        error!("Failed to handshake: {e}");
                    }
                }

                self.client.on_message(move |_event, _payload| {});
            })
            .await;

        match result {
            Ok(_) => {}
            Err(err) => {
                eprintln!("Error: {}", err.to_string())
            }
        }
    }

    pub fn shutdown(&self) {
        let vector_file = self.data_path.join(INDEX_FILE_NAME);

        if let Err(e) = self.vector_index.index.write(&vector_file) {
            error!("Failed to save vector index: {e}");
        }
    }
}

#[cfg(unix)]
async fn wait_for_signal() {
    use tokio::signal::unix::{SignalKind, signal};

    let mut signal_terminate = signal(SignalKind::terminate()).unwrap();
    let mut signal_interrupt = signal(SignalKind::interrupt()).unwrap();

    tokio::select! {
        _ = signal_terminate.recv() => {},
        _ = signal_interrupt.recv() => {},
    };
}

#[cfg(windows)]
async fn wait_for_signal() {
    use tokio::signal::windows;

    let mut signal_c = windows::ctrl_c().unwrap();
    let mut signal_break = windows::ctrl_break().unwrap();
    let mut signal_close = windows::ctrl_close().unwrap();
    let mut signal_shutdown = windows::ctrl_shutdown().unwrap();

    tokio::select! {
        _ = signal_c.recv() => {},
        _ = signal_break.recv() => {},
        _ = signal_close.recv() => {},
        _ = signal_shutdown.recv() => {},
    };
}

#[tokio::main]
async fn main() {
    match VectorModule::new() {
        Ok(module) => {
            tokio::select! {
                _ = module.start() => {},
                _ = wait_for_signal() => {},
            }
            module.shutdown();
        }
        Err(e) => eprintln!("Failed to start module: {e}"),
    };
}
