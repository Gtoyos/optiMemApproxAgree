//2 process Optimal memory approximate agreement protocol implementation.
use crate::snapshot::SnapshotAtomic;
use tokio::time::{sleep, Duration};
use std::sync::Arc;
use rand::Rng;

//Define the comm messages
#[derive(Debug,Clone,Copy,PartialEq,Eq)]
enum Msg {
    Bot = 0,
    A = 1,
    B = 2,
}

impl From<u8> for Msg{
    fn from(value: u8) -> Self {
        match value {
            0 => Msg::Bot,
            1 => Msg::A,
            2 => Msg::B,
            _ => {panic!("Out of range")},
        }
    }
}

pub async fn agreement_protocol(
    input : usize,
    rounds : usize,
    snapshots: Vec<Arc<SnapshotAtomic<u8>>>,
) -> f64 {
    let pid: i128 = input as i128;
    let mut state: i128 = 0;
    for i in 0..rounds {
        //write an encoding of the state
        snapshots[i].write( ((state%2) + 1) as u8, pid as usize).await;

        //get the view of the other process.
        let v = Msg::from(snapshots[i].snap().await[1-pid as usize]);
        println!("process {}, round {}, current state {}",pid,i,state);
        //compute next_state function
        match v { 
            Msg::B   => state = 3*state + pid + (1-2*(state%2))*(2*pid-1), // (-1)^s=1-2(s%2)
            Msg::Bot => state = 3*state + pid,
            Msg::A   => state = 3*state + pid + (1-2*(state%2))*(1-2*pid),
        }
        println!("process {}, round {}, read {:?}, new state is {}",pid,i,v,state);
    }
    
    //simulate random delay
    let random = rand::random::<u64>();
    tokio::time::sleep(Duration::from_millis(random % 5000)).await;

    let result = (2*state+pid) as f64 / (3f64.powi(rounds as i32));
    return result;
}