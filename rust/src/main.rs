use optimemapproxagree::protocol::agreement_protocol;
use optimemapproxagree::snapshot::SnapshotAtomic;

use std::sync::Arc;
use std::{env,f64};

#[tokio::main]
async fn main() {
    //Parse command-line args
    let args: Vec<String> = env::args().collect();
    let input1: f64 = args.get(1).and_then(|s| s.parse().ok()).unwrap_or(0.0);
    let input2: f64 = args.get(2).and_then(|s| s.parse().ok()).unwrap_or(1.0);
    let agreement_lvl: f64 = args.get(3).and_then(|s| s.parse().ok()).unwrap_or(0.000001);
 
    //Compute required rounds based on agreement lvl and input.
    let r = ((1.0 / f64::ln(3.0)) * (f64::ln((input1-input2).abs()) - f64::ln(agreement_lvl))).ceil() as usize;

    println!("--- Approximate agreement 2-process task ---");
    println!("Process 0 input: {}", input1);
    println!("Process 1 input: {}", input2);
    println!("Agreement level: {}", agreement_lvl);
    println!("--------------------------------------------");
    println!("Agreement level requires {} communication rounds.", r);

    //Initialize IIS shared memory datastructure
    let snapshots: Vec<Arc<SnapshotAtomic<u8>>> = (0..r)
        .map(|_| Arc::new(SnapshotAtomic::new(2)))
        .collect();

    //Run the agreement protocol as two separate tasks
    let task1 = tokio::spawn(agreement_protocol(0, r, snapshots.clone()));
    let task2 = tokio::spawn(agreement_protocol(1, r, snapshots.clone()));
    
    // Ensure both tasks complete
    let r0 = task1.await.unwrap();
    let r1 = task2.await.unwrap();
    
    //Compute final values
    let decision1 = input1.min(input2) + (input1-input2).abs() * r0;
    let decision2 = input1.min(input2) + (input1-input2).abs() * r1;

    println!("Process 0 decision output: {}", decision1);
    println!("Process 1 decision output: {}", decision2);
    println!("Agreement task completed!");
}
