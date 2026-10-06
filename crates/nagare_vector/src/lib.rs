use std::ffi::c_char;
use std::path::Path;
use turbovec::IdMapIndex;

use crate::error::{LAST_ERROR, set_error};
use crate::vector_index::VectorIndex;

#[unsafe(no_mangle)]
pub extern "C" fn vector_index_last_error() -> *const c_char {
    match LAST_ERROR.lock() {
        Ok(error) => match error.as_ref() {
            Some(error) => error.as_ptr(),
            None => std::ptr::null(),
        },
        Err(_) => std::ptr::null(),
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn vector_index_create(dim: usize, bits: usize) -> *mut VectorIndex {
    match IdMapIndex::new(dim, bits) {
        Ok(index) => Box::into_raw(Box::new(VectorIndex { index })),

        Err(err) => {
            set_error(format!("failed to create IdMapIndex: {err}"));
            std::ptr::null_mut()
        }
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn vector_index_load(path: *const std::ffi::c_char) -> *mut VectorIndex {
    if path.is_null() {
        return std::ptr::null_mut();
    }

    let path = match unsafe { std::ffi::CStr::from_ptr(path) }.to_str() {
        Ok(path) => std::path::Path::new(path),
        Err(err) => {
            set_error(format!("failed to convert path to string: {err}"));
            return std::ptr::null_mut();
        }
    };

    match IdMapIndex::load(path) {
        Ok(index) => Box::into_raw(Box::new(VectorIndex { index })),
        Err(err) => {
            set_error(format!("failed to load IdMapIndex: {err}"));
            std::ptr::null_mut()
        }
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn vector_index_sync(
    ptr: *mut VectorIndex,
    path: *const std::ffi::c_char,
) -> i32 {
    if ptr.is_null() || path.is_null() {
        return -1;
    }

    let index = unsafe { &mut *ptr };

    let path = match unsafe { std::ffi::CStr::from_ptr(path) }.to_str() {
        Ok(path) => Path::new(path),
        Err(err) => {
            set_error(format!("failed to convert path to string: {err}"));
            return -2;
        }
    };

    match index.index.sync(path) {
        Ok(_) => 0,
        Err(err) => {
            set_error(format!("failed to sync vector index: {err}"));
            return -3;
        }
    }
}

mod error;
mod vector_index;
