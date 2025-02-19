// src/snapshot.rs
// Immediate snapshot implementation.

use std::sync::Arc;
use tokio::sync::RwLock;

pub struct SnapshotAtomic<T: Clone + Default>{
    raw_array : Arc<RwLock<Vec<T>>>,
}

impl<T: Clone + Default> SnapshotAtomic<T> {

    //Create a IS instance
    pub fn new(size: usize) -> Self {
        let arr = vec![T::default(); size]; //an array with the default init value
        return Self {
            raw_array: Arc::new(RwLock::new(arr)),
        }
    }

    //Takes a snapshot of the array
    pub async fn snap(&self) -> Vec<T> {
        let snap = self.raw_array.read().await;
        return snap.clone();
    }
    
    //Write value at snapshot
    pub async fn write(&self, value: T, index: usize){
        let mut write_ptr = self.raw_array.write().await;
        if index < write_ptr.len() {
            write_ptr[index] = value;
        }
    }
}
