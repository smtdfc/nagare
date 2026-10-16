use std::collections::HashMap;
use std::sync::{LazyLock, Mutex};

pub static ERRORS: LazyLock<Mutex<HashMap<usize, String>>> =
    LazyLock::new(|| Mutex::new(HashMap::new()));

pub fn set_error(id: usize, message: impl Into<String>) {
    let message = message.into();

    match ERRORS.lock() {
        Ok(mut errors) => {
            errors.insert(id, message);
        }
        Err(err) => {
            eprintln!("Failed to lock ERRORS: {err}");
        }
    }
}

pub fn get_error(id: usize) -> Option<String> {
    match ERRORS.lock() {
        Ok(errors) => errors.get(&id).cloned(),
        Err(err) => {
            eprintln!("Failed to lock ERRORS: {err}");
            None
        }
    }
}
