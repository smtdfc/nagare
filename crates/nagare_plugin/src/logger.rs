use crate::PluginError;
use std::sync::OnceLock;
use std::time::SystemTime;

static LOG_INIT: OnceLock<()> = OnceLock::new();

pub fn init_logging_once(path: &str) {
    LOG_INIT.get_or_init(|| {
        match setup_logger(path) {
            Ok(_) => {}
            Err(e) => {
                eprintln!("Error: {e}");
            }
        };
    });
}

pub fn setup_logger(path: &str) -> Result<(), PluginError> {
    fern::Dispatch::new()
        .format(|out, message, record| {
            out.finish(format_args!(
                "[{} {} {}] {}",
                humantime::format_rfc3339_seconds(SystemTime::now()),
                record.level(),
                record.target(),
                message
            ))
        })
        .level(log::LevelFilter::Debug)
        .chain(std::io::stdout())
        .chain(fern::log_file(path)?)
        .apply()?;
    Ok(())
}
