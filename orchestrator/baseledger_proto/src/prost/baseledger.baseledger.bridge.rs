/// Params defines the parameters for the module.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attestation {
    #[prost(bool, tag="1")]
    pub observed: bool,
    #[prost(string, repeated, tag="2")]
    pub votes: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(uint64, tag="3")]
    pub height: u64,
    #[prost(message, optional, tag="4")]
    pub claim: ::core::option::Option<::prost_types::Any>,
    #[prost(string, repeated, tag="5")]
    pub ubt_prices: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ClaimType {
    Unspecified = 0,
    ClaimUbtDeposited = 1,
    ClaimValidatorPowerChanged = 2,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OrchestratorValidatorAddress {
    #[prost(string, tag="1")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub validator_address: ::prost::alloc::string::String,
}
/// QueryParamsRequest is request type for the Query/Params RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {
}
/// QueryParamsResponse is response type for the Query/Params RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    /// params holds all the parameters of this module.
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastEventNonceByAddressRequest {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastEventNonceByAddressResponse {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAttestationsRequest {
    #[prost(uint64, tag="1")]
    pub limit: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAttestationsResponse {
    #[prost(message, repeated, tag="1")]
    pub attestations: ::prost::alloc::vec::Vec<Attestation>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetOrchestratorValidatorAddressRequest {
    #[prost(string, tag="1")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetOrchestratorValidatorAddressResponse {
    #[prost(message, optional, tag="1")]
    pub orchestrator_validator_address: ::core::option::Option<OrchestratorValidatorAddress>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllOrchestratorValidatorAddressRequest {
    #[prost(message, optional, tag="1")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllOrchestratorValidatorAddressResponse {
    #[prost(message, repeated, tag="1")]
    pub orchestrator_validator_address: ::prost::alloc::vec::Vec<OrchestratorValidatorAddress>,
    #[prost(message, optional, tag="2")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs , clippy :: let_unit_value ,)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service."] # [derive (Debug , Clone)] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as Body > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor < F > (inner : T , interceptor : F) -> QueryClient < InterceptedService < T , F >> where F : tonic :: service :: Interceptor , T : tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody > , Response = http :: Response << T as tonic :: client :: GrpcService < tonic :: body :: BoxBody >> :: ResponseBody > > , < T as tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody >> > :: Error : Into < StdError > + Send + Sync , { QueryClient :: new (InterceptedService :: new (inner , interceptor)) } # [doc = r" Compress requests with `gzip`."] # [doc = r""] # [doc = r" This requires the server to support it otherwise it might respond with an"] # [doc = r" error."] pub fn send_gzip (mut self) -> Self { self . inner = self . inner . send_gzip () ; self } # [doc = r" Enable decompressing responses with `gzip`."] pub fn accept_gzip (mut self) -> Self { self . inner = self . inner . accept_gzip () ; self } # [doc = " Parameters queries the parameters of the module."] pub async fn params (& mut self , request : impl tonic :: IntoRequest < super :: QueryParamsRequest > ,) -> Result < tonic :: Response < super :: QueryParamsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Query/Params") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of LastEventNonceByAddress items."] pub async fn last_event_nonce_by_address (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastEventNonceByAddressRequest > ,) -> Result < tonic :: Response < super :: QueryLastEventNonceByAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Query/LastEventNonceByAddress") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of Attestations items."] pub async fn attestations (& mut self , request : impl tonic :: IntoRequest < super :: QueryAttestationsRequest > ,) -> Result < tonic :: Response < super :: QueryAttestationsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Query/Attestations") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a OrchestratorValidatorAddress by index."] pub async fn orchestrator_validator_address (& mut self , request : impl tonic :: IntoRequest < super :: QueryGetOrchestratorValidatorAddressRequest > ,) -> Result < tonic :: Response < super :: QueryGetOrchestratorValidatorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Query/OrchestratorValidatorAddress") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of OrchestratorValidatorAddress items."] pub async fn orchestrator_validator_address_all (& mut self , request : impl tonic :: IntoRequest < super :: QueryAllOrchestratorValidatorAddressRequest > ,) -> Result < tonic :: Response < super :: QueryAllOrchestratorValidatorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Query/OrchestratorValidatorAddressAll") ; self . inner . unary (request . into_request () , path , codec) . await } } }#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUbtDepositedClaim {
    #[prost(string, tag="1")]
    pub creator: ::prost::alloc::string::String,
    #[prost(uint64, tag="2")]
    pub event_nonce: u64,
    #[prost(uint64, tag="3")]
    pub block_height: u64,
    #[prost(string, tag="4")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag="6")]
    pub ethereum_sender: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub cosmos_receiver: ::prost::alloc::string::String,
    #[prost(string, tag="8")]
    pub ubt_price: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUbtDepositedClaimResponse {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgValidatorPowerChangedClaim {
    #[prost(string, tag="1")]
    pub creator: ::prost::alloc::string::String,
    #[prost(uint64, tag="2")]
    pub event_nonce: u64,
    #[prost(uint64, tag="3")]
    pub block_height: u64,
    #[prost(string, tag="4")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag="6")]
    pub ethereum_sender: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub cosmos_receiver: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgValidatorPowerChangedClaimResponse {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCreateOrchestratorValidatorAddress {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCreateOrchestratorValidatorAddressResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs , clippy :: let_unit_value ,)] use tonic :: codegen :: * ; # [doc = " Msg defines the Msg service."] # [derive (Debug , Clone)] pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as Body > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor < F > (inner : T , interceptor : F) -> MsgClient < InterceptedService < T , F >> where F : tonic :: service :: Interceptor , T : tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody > , Response = http :: Response << T as tonic :: client :: GrpcService < tonic :: body :: BoxBody >> :: ResponseBody > > , < T as tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody >> > :: Error : Into < StdError > + Send + Sync , { MsgClient :: new (InterceptedService :: new (inner , interceptor)) } # [doc = r" Compress requests with `gzip`."] # [doc = r""] # [doc = r" This requires the server to support it otherwise it might respond with an"] # [doc = r" error."] pub fn send_gzip (mut self) -> Self { self . inner = self . inner . send_gzip () ; self } # [doc = r" Enable decompressing responses with `gzip`."] pub fn accept_gzip (mut self) -> Self { self . inner = self . inner . accept_gzip () ; self } pub async fn ubt_deposited_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgUbtDepositedClaim > ,) -> Result < tonic :: Response < super :: MsgUbtDepositedClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Msg/UbtDepositedClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn validator_power_changed_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgValidatorPowerChangedClaim > ,) -> Result < tonic :: Response < super :: MsgValidatorPowerChangedClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Msg/ValidatorPowerChangedClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn create_orchestrator_validator_address (& mut self , request : impl tonic :: IntoRequest < super :: MsgCreateOrchestratorValidatorAddress > ,) -> Result < tonic :: Response < super :: MsgCreateOrchestratorValidatorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledger.bridge.Msg/CreateOrchestratorValidatorAddress") ; self . inner . unary (request . into_request () , path , codec) . await } } }/// GenesisState defines the bridge module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, repeated, tag="5")]
    pub orchestrator_validator_address_list: ::prost::alloc::vec::Vec<OrchestratorValidatorAddress>,
    /// this line is used by starport scaffolding # genesis/proto/state
    #[prost(message, repeated, tag="2")]
    pub attestations: ::prost::alloc::vec::Vec<Attestation>,
    #[prost(uint64, tag="3")]
    pub last_observed_nonce: u64,
}
