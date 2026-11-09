mod handlers;
mod module;
mod types;

use std::process::exit;

use module::ComputerUseModule;

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
    match ComputerUseModule::new() {
        Ok(module) => {
            tokio::select! {
                _ = module.clone().start() => {},
                _ = wait_for_signal() => {},
            }
            module.shutdown(true);
        }
        Err(e) => {
            eprintln!("Failed to start module: {e}");
            exit(1)
        }
    };
}
