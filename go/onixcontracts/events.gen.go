// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    rawEvent, err := UnmarshalRawEvent(bytes)
//    bytes, err = rawEvent.Marshal()
//
//    normEvent, err := UnmarshalNormEvent(bytes)
//    bytes, err = normEvent.Marshal()
//
//    cleanEvent, err := UnmarshalCleanEvent(bytes)
//    bytes, err = cleanEvent.Marshal()
//
//    wsMessage, err := UnmarshalWsMessage(bytes)
//    bytes, err = wsMessage.Marshal()

package onixcontracts

import "time"

import "encoding/json"

func UnmarshalRawEvent(data []byte) (RawEvent, error) {
	var r RawEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RawEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalNormEvent(data []byte) (NormEvent, error) {
	var r NormEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NormEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalCleanEvent(data []byte) (CleanEvent, error) {
	var r CleanEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CleanEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalWsMessage(data []byte) (WsMessage, error) {
	var r WsMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *WsMessage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Evento crudo de telemetría. Lo publica onix-hook en onix.raw.<session> y lo consume
// onix-ingestor. No contiene análisis ni redacción de credenciales todavía.
type RawEvent struct {
	// Rol del agente que generó el evento (deriva de ONIX_AGENT_ROLE o de la skill activa).                          
	AgentRole                                                                                  AgentRole              `json:"agent_role"`
	// Costo estimado en USD de la acción.                                                                            
	CostUsd                                                                                    *float64               `json:"cost_usd,omitempty"`
	// Directorio de trabajo del agente cuando ocurrió el evento.                                                     
	Cwd                                                                                        *string                `json:"cwd,omitempty"`
	// Hook de Claude Code que disparó el evento.                                                                     
	Hook                                                                                       Hook                   `json:"hook"`
	// Parámetros crudos de la tool. PUEDE contener credenciales: onix-guard los redacta más                          
	// adelante. Nunca se persiste tal cual.                                                                          
	Params                                                                                     map[string]interface{} `json:"params,omitempty"`
	// Identificador del proyecto monitoreado (ONIX_PROJECT).                                                         
	Project                                                                                    string                 `json:"project"`
	Result                                                                                     *RawEventResult        `json:"result,omitempty"`
	// claude_session_id: identifica la sesión-agente que originó el evento.                                          
	Session                                                                                    string                 `json:"session"`
	// Número de etapa (1..12) en curso, si aplica.                                                                   
	Stage                                                                                      *int64                 `json:"stage,omitempty"`
	// Tokens consumidos por la acción, si el hook los reporta.                                                       
	Tokens                                                                                     *int64                 `json:"tokens,omitempty"`
	// Nombre de la herramienta invocada (Bash, Read, Edit, WebSearch...). Ausente en hooks que                       
	// no son de tool.                                                                                                
	Tool                                                                                       *string                `json:"tool,omitempty"`
	// Momento del evento en el lado del agente (ISO-8601 UTC).                                                       
	Ts                                                                                         time.Time              `json:"ts"`
	// Versión del contrato de evento. Permite evolucionar el formato sin romper consumidores.                        
	V                                                                                          int64                  `json:"v"`
}

// Resultado de la ejecución de la tool.
type RawEventResult struct {
	// Duración de la acción en milisegundos.                          
	DurationMS                                                  *int64 `json:"duration_ms,omitempty"`
	// Código de salida (0 = éxito). exit_code != 0 marca error.       
	ExitCode                                                    *int64 `json:"exit_code,omitempty"`
}

// Evento validado y normalizado. Lo publica onix-ingestor en onix.norm.<session> y lo
// consume onix-guard. Misma forma que RawEvent pero con campos garantizados/normalizados
// (ts en UTC, params presentes). Aún NO está redactado ni analizado.
type NormEvent struct {
	// Rol del agente.                                                                                         
	AgentRole                                                                           AgentRole              `json:"agent_role"`
	// Costo en USD.                                                                                           
	CostUsd                                                                             *float64               `json:"cost_usd,omitempty"`
	// Directorio de trabajo.                                                                                  
	Cwd                                                                                 *string                `json:"cwd,omitempty"`
	// Hook que disparó el evento.                                                                             
	Hook                                                                                Hook                   `json:"hook"`
	// Parámetros de la tool (pueden contener credenciales; se redactan en onix-guard).                        
	Params                                                                              map[string]interface{} `json:"params,omitempty"`
	// Proyecto monitoreado.                                                                                   
	Project                                                                             string                 `json:"project"`
	// Momento en que onix-ingestor recibió el evento (UTC). Lo agrega la normalización.                       
	ReceivedAt                                                                          time.Time              `json:"received_at"`
	Result                                                                              *NormEventResult       `json:"result,omitempty"`
	// claude_session_id.                                                                                      
	Session                                                                             string                 `json:"session"`
	// Etapa en curso.                                                                                         
	Stage                                                                               *int64                 `json:"stage,omitempty"`
	// Tokens consumidos.                                                                                      
	Tokens                                                                              *int64                 `json:"tokens,omitempty"`
	// Nombre de la tool, si aplica.                                                                           
	Tool                                                                                *string                `json:"tool,omitempty"`
	// Timestamp del agente (UTC).                                                                             
	Ts                                                                                  time.Time              `json:"ts"`
	// Versión del contrato.                                                                                   
	V                                                                                   int64                  `json:"v"`
}

type NormEventResult struct {
	// Duración en ms.         
	DurationMS          *int64 `json:"duration_ms,omitempty"`
	// Código de salida.       
	ExitCode            *int64 `json:"exit_code,omitempty"`
}

// Evento limpio, enriquecido y SIN credenciales. Lo publica onix-guard en
// onix.clean.<session>; lo consumen onix-recorder (persiste) y onix-gateway (empuja al
// vivo). El valor real de cualquier credencial ya fue eliminado: solo queda su hash SHA-256.
type CleanEvent struct {
	// Rol del agente.                                                                                                 
	AgentRole                                                                                   AgentRole              `json:"agent_role"`
	// Costo en USD.                                                                                                   
	CostUsd                                                                                     *float64               `json:"cost_usd,omitempty"`
	// Credenciales DETECTADAS y redactadas en este evento. Solo metadatos + hash; NUNCA el                            
	// valor real.                                                                                                     
	Credentials                                                                                 []EventCleanSchema     `json:"credentials"`
	// Directorio de trabajo.                                                                                          
	Cwd                                                                                         *string                `json:"cwd,omitempty"`
	// Hook que disparó el evento.                                                                                     
	Hook                                                                                        Hook                   `json:"hook"`
	// true si la tool falló o exit_code != 0. Lo marca onix-guard.                                                    
	IsError                                                                                     bool                   `json:"is_error"`
	// true si onix-guard detectó una repetición según la regla configurada.                                           
	IsRepetition                                                                                bool                   `json:"is_repetition"`
	// SHA-256 de params_normalized. Base de la regla de repetición (mismo tool + params_hash >=                       
	// 3 veces).                                                                                                       
	ParamsHash                                                                                  string                 `json:"params_hash"`
	// Parámetros YA redactados (sin valores de credenciales) y normalizados para comparar                             
	// repeticiones.                                                                                                   
	ParamsNormalized                                                                            map[string]interface{} `json:"params_normalized,omitempty"`
	// Proyecto monitoreado.                                                                                           
	Project                                                                                     string                 `json:"project"`
	// Cuando lo recibió el ingestor (UTC).                                                                            
	ReceivedAt                                                                                  time.Time              `json:"received_at"`
	Result                                                                                      *CleanEventResult      `json:"result,omitempty"`
	// claude_session_id.                                                                                              
	Session                                                                                     string                 `json:"session"`
	// Etapa en curso.                                                                                                 
	Stage                                                                                       *int64                 `json:"stage,omitempty"`
	// Tokens consumidos.                                                                                              
	Tokens                                                                                      *int64                 `json:"tokens,omitempty"`
	// Nombre de la tool, si aplica.                                                                                   
	Tool                                                                                        *string                `json:"tool,omitempty"`
	// Timestamp del agente (UTC).                                                                                     
	Ts                                                                                          time.Time              `json:"ts"`
	// Versión del contrato.                                                                                           
	V                                                                                           int64                  `json:"v"`
}

// Credencial detectada. El contrato prohíbe transportar el valor real.
type EventCleanSchema struct {
	// Tipo de credencial detectada.                                                            
	Kind                                                                                 Kind   `json:"kind"`
	// Etiqueta legible que reemplaza al valor, p. ej. '[credencial · sha256:9f2c…e41a]'.       
	Label                                                                                string `json:"label"`
	// SHA-256 del valor real (lo único que se conserva).                                       
	Sha256                                                                               string `json:"sha256"`
}

type CleanEventResult struct {
	// Duración en ms.         
	DurationMS          *int64 `json:"duration_ms,omitempty"`
	// Código de salida.       
	ExitCode            *int64 `json:"exit_code,omitempty"`
}

// Mensaje que onix-gateway envía al frontend por WebSocket (/ws). Sobre estos type el
// cliente actualiza la lista de actividad, tarjetas de agente, oficina pixel y contadores.
type WsMessage struct {
	// Carga útil; su forma depende de type (CleanEvent para 'event', métricas para 'metrics',                       
	// etc.).                                                                                                        
	Data                                                                                      map[string]interface{} `json:"data"`
	// Categoría del mensaje en vivo.                                                                                
	Type                                                                                      Type                   `json:"type"`
}

// Rol del agente que generó el evento (deriva de ONIX_AGENT_ROLE o de la skill activa).
//
// Rol del agente.
type AgentRole string

const (
	Designer    AgentRole = "designer"
	Fullstack   AgentRole = "fullstack"
	Gm          AgentRole = "gm"
	Growth      AgentRole = "growth"
	ProjectLead AgentRole = "project_lead"
	Sales       AgentRole = "sales"
)

// Hook de Claude Code que disparó el evento.
//
// Hook que disparó el evento.
type Hook string

const (
	Notification Hook = "Notification"
	PostToolUse  Hook = "PostToolUse"
	PreToolUse   Hook = "PreToolUse"
	SessionStart Hook = "SessionStart"
	Stop         Hook = "Stop"
	SubagentStop Hook = "SubagentStop"
)

// Tipo de credencial detectada.
type Kind string

const (
	APIKey     Kind = "api_key"
	ConnString Kind = "conn_string"
	Env        Kind = "env"
	Password   Kind = "password"
	Token      Kind = "token"
)

// Categoría del mensaje en vivo.
type Type string

const (
	AgentStatus Type = "agent_status"
	Alert       Type = "alert"
	Event       Type = "event"
	Metrics     Type = "metrics"
	Report      Type = "report"
	Stage       Type = "stage"
)
