use std::path::Path;

use turbovec::IdMapIndex;

pub struct VectorIndex {
    pub index: IdMapIndex,
}

impl VectorIndex {
    pub fn new(dim: usize, bits: usize) -> Result<Self, String> {
        let index = IdMapIndex::new(dim, bits)
            .map_err(|err| format!("Failed to create IdMapIndex: {}", err))?;
        Ok(Self { index })
    }

    pub fn load(path: &Path) -> Result<Self, String> {
        let index =
            IdMapIndex::load(path).map_err(|err| format!("Failed to load IdMapIndex: {}", err))?;
        Ok(Self { index })
    }
}
