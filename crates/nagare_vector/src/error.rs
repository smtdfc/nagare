use std::{ffi::CString, sync::Mutex};

pub static LAST_ERROR: Mutex<Option<CString>> = Mutex::new(None);

pub fn set_error(message: impl Into<String>) {
    let message = message.into();

    if let Ok(mut error) = LAST_ERROR.lock() {
        *error = CString::new(message).ok();
    }
}
