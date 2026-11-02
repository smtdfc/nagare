pub mod connector;
pub mod errors;
pub mod ipc;

pub use connector::{Connector, OnEvent};
pub use errors::IpcError;
pub use ipc::{MAX_MESSAGE_SIZE, encode_message};
