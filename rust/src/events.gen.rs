// Example code that deserializes and serializes the model.
// extern crate serde;
// #[macro_use]
// extern crate serde_derive;
// extern crate serde_json;
//
// use generated_module::RawEvent;
//
// fn main() {
//     let json = r#"{"answer": 42}"#;
//     let model: RawEvent = serde_json::from_str(&json).unwrap();
// }

use serde::{Serialize, Deserialize};
use std::collections::HashMap;

/// Evento crudo de telemetría. Lo publica onix-hook en onix.raw.<session> y lo consume
/// onix-ingestor. No contiene análisis ni redacción de credenciales todavía.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RawEvent {
    /// Rol del agente que generó el evento (deriva de ONIX_AGENT_ROLE o de la skill activa).
    pub agent_role: AgentRole,

    /// Costo estimado en USD de la acción.
    pub cost_usd: Option<f64>,

    /// Directorio de trabajo del agente cuando ocurrió el evento.
    pub cwd: Option<String>,

    /// Hook de Claude Code que disparó el evento.
    pub hook: Hook,

    /// Parámetros crudos de la tool. PUEDE contener credenciales: onix-guard los redacta más
    /// adelante. Nunca se persiste tal cual.
    pub params: Option<HashMap<String, Option<serde_json::Value>>>,

    /// Identificador del proyecto monitoreado (ONIX_PROJECT).
    pub project: String,

    pub result: Option<RawEventResult>,

    /// claude_session_id: identifica la sesión-agente que originó el evento.
    pub session: String,

    /// Número de etapa (1..12) en curso, si aplica.
    pub stage: Option<i32>,

    /// Tokens consumidos por la acción, si el hook los reporta.
    pub tokens: Option<i64>,

    /// Nombre de la herramienta invocada (Bash, Read, Edit, WebSearch...). Ausente en hooks que
    /// no son de tool.
    pub tool: Option<String>,

    /// Momento del evento en el lado del agente (ISO-8601 UTC).
    pub ts: String,

    /// Versión del contrato de evento. Permite evolucionar el formato sin romper consumidores.
    pub v: i64,
}

/// Rol del agente que generó el evento (deriva de ONIX_AGENT_ROLE o de la skill activa).
///
/// Rol del agente.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum AgentRole {
    Designer,

    Fullstack,

    Gm,

    Growth,

    #[serde(rename = "project_lead")]
    ProjectLead,

    Sales,
}

/// Hook de Claude Code que disparó el evento.
///
/// Hook que disparó el evento.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Hook {
    Notification,

    #[serde(rename = "PostToolUse")]
    PostToolUse,

    #[serde(rename = "PreToolUse")]
    PreToolUse,

    #[serde(rename = "SessionStart")]
    SessionStart,

    Stop,

    #[serde(rename = "SubagentStop")]
    SubagentStop,
}

/// Resultado de la ejecución de la tool.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RawEventResult {
    /// Duración de la acción en milisegundos.
    pub duration_ms: Option<i64>,

    /// Código de salida (0 = éxito). exit_code != 0 marca error.
    pub exit_code: Option<i64>,
}

/// Evento validado y normalizado. Lo publica onix-ingestor en onix.norm.<session> y lo
/// consume onix-guard. Misma forma que RawEvent pero con campos garantizados/normalizados
/// (ts en UTC, params presentes). Aún NO está redactado ni analizado.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NormEvent {
    /// Rol del agente.
    pub agent_role: AgentRole,

    /// Costo en USD.
    pub cost_usd: Option<f64>,

    /// Directorio de trabajo.
    pub cwd: Option<String>,

    /// Hook que disparó el evento.
    pub hook: Hook,

    /// Parámetros de la tool (pueden contener credenciales; se redactan en onix-guard).
    pub params: Option<HashMap<String, Option<serde_json::Value>>>,

    /// Proyecto monitoreado.
    pub project: String,

    /// Momento en que onix-ingestor recibió el evento (UTC). Lo agrega la normalización.
    pub received_at: String,

    pub result: Option<NormEventResult>,

    /// claude_session_id.
    pub session: String,

    /// Etapa en curso.
    pub stage: Option<i32>,

    /// Tokens consumidos.
    pub tokens: Option<i64>,

    /// Nombre de la tool, si aplica.
    pub tool: Option<String>,

    /// Timestamp del agente (UTC).
    pub ts: String,

    /// Versión del contrato.
    pub v: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NormEventResult {
    /// Duración en ms.
    pub duration_ms: Option<i64>,

    /// Código de salida.
    pub exit_code: Option<i64>,
}

/// Evento limpio, enriquecido y SIN credenciales. Lo publica onix-guard en
/// onix.clean.<session>; lo consumen onix-recorder (persiste) y onix-gateway (empuja al
/// vivo). El valor real de cualquier credencial ya fue eliminado: solo queda su hash SHA-256.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CleanEvent {
    /// Rol del agente.
    pub agent_role: AgentRole,

    /// Costo en USD.
    pub cost_usd: Option<f64>,

    /// Credenciales DETECTADAS y redactadas en este evento. Solo metadatos + hash; NUNCA el
    /// valor real.
    pub credentials: Vec<EventCleanSchema>,

    /// Directorio de trabajo.
    pub cwd: Option<String>,

    /// Hook que disparó el evento.
    pub hook: Hook,

    /// true si la tool falló o exit_code != 0. Lo marca onix-guard.
    pub is_error: bool,

    /// true si onix-guard detectó una repetición según la regla configurada.
    pub is_repetition: bool,

    /// SHA-256 de params_normalized. Base de la regla de repetición (mismo tool + params_hash >=
    /// 3 veces).
    pub params_hash: String,

    /// Parámetros YA redactados (sin valores de credenciales) y normalizados para comparar
    /// repeticiones.
    pub params_normalized: Option<HashMap<String, Option<serde_json::Value>>>,

    /// Proyecto monitoreado.
    pub project: String,

    /// Cuando lo recibió el ingestor (UTC).
    pub received_at: String,

    pub result: Option<CleanEventResult>,

    /// claude_session_id.
    pub session: String,

    /// Etapa en curso.
    pub stage: Option<i32>,

    /// Tokens consumidos.
    pub tokens: Option<i64>,

    /// Nombre de la tool, si aplica.
    pub tool: Option<String>,

    /// Timestamp del agente (UTC).
    pub ts: String,

    /// Versión del contrato.
    pub v: i64,
}

/// Credencial detectada. El contrato prohíbe transportar el valor real.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EventCleanSchema {
    /// Tipo de credencial detectada.
    pub kind: Kind,

    /// Etiqueta legible que reemplaza al valor, p. ej. '[credencial · sha256:9f2c…e41a]'.
    pub label: String,

    /// SHA-256 del valor real (lo único que se conserva).
    pub sha256: String,
}

/// Tipo de credencial detectada.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum Kind {
    #[serde(rename = "api_key")]
    ApiKey,

    #[serde(rename = "conn_string")]
    ConnString,

    Env,

    Password,

    Token,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CleanEventResult {
    /// Duración en ms.
    pub duration_ms: Option<i64>,

    /// Código de salida.
    pub exit_code: Option<i64>,
}

/// Mensaje que onix-gateway envía al frontend por WebSocket (/ws). Sobre estos type el
/// cliente actualiza la lista de actividad, tarjetas de agente, oficina pixel y contadores.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WsMessage {
    /// Carga útil; su forma depende de type (CleanEvent para 'event', métricas para 'metrics',
    /// etc.).
    pub data: HashMap<String, Option<serde_json::Value>>,

    /// Categoría del mensaje en vivo.
    #[serde(rename = "type")]
    pub ws_message_type: Type,
}

/// Categoría del mensaje en vivo.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum Type {
    #[serde(rename = "agent_status")]
    AgentStatus,

    Alert,

    Event,

    Metrics,

    Report,

    Stage,
}
