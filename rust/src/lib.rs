//! Crate `onix-contracts`: tipos compartidos de OnixGuard para los servicios en Rust
//! (principalmente `onix-guard`). Los tipos de `events` son GENERADOS desde
//! `schemas/*.schema.json` con `codegen/generate.sh` — no editar a mano.

pub mod events {
    include!("events.gen.rs");
}

pub use events::*;
