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
    #[prost(string, tag="6")]
    pub avg_ubt_price: ::prost::alloc::string::String,
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ClaimType {
    Unspecified = 0,
    ClaimUbtDeposited = 1,
}
/// Params defines the parameters for the module.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
}
/// GenesisState defines the baseledgerbridge module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    /// this line is used by starport scaffolding # genesis/proto/state
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
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
pub struct QueryDelegateKeysByEthAddressRequest {
    #[prost(string, tag="1")]
    pub eth_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryDelegateKeysByEthAddressResponse {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryDelegateKeysByOrchestratorAddressRequest {
    #[prost(string, tag="1")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryDelegateKeysByOrchestratorAddressResponse {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub eth_address: ::prost::alloc::string::String,
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
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs , clippy :: let_unit_value ,)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service."] # [derive (Debug , Clone)] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as Body > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor < F > (inner : T , interceptor : F) -> QueryClient < InterceptedService < T , F >> where F : tonic :: service :: Interceptor , T : tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody > , Response = http :: Response << T as tonic :: client :: GrpcService < tonic :: body :: BoxBody >> :: ResponseBody > > , < T as tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody >> > :: Error : Into < StdError > + Send + Sync , { QueryClient :: new (InterceptedService :: new (inner , interceptor)) } # [doc = r" Compress requests with `gzip`."] # [doc = r""] # [doc = r" This requires the server to support it otherwise it might respond with an"] # [doc = r" error."] pub fn send_gzip (mut self) -> Self { self . inner = self . inner . send_gzip () ; self } # [doc = r" Enable decompressing responses with `gzip`."] pub fn accept_gzip (mut self) -> Self { self . inner = self . inner . accept_gzip () ; self } # [doc = " Parameters queries the parameters of the module."] pub async fn params (& mut self , request : impl tonic :: IntoRequest < super :: QueryParamsRequest > ,) -> Result < tonic :: Response < super :: QueryParamsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Query/Params") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of DelegateKeysByEthAddress items."] pub async fn delegate_keys_by_eth_address (& mut self , request : impl tonic :: IntoRequest < super :: QueryDelegateKeysByEthAddressRequest > ,) -> Result < tonic :: Response < super :: QueryDelegateKeysByEthAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Query/DelegateKeysByEthAddress") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of DelegateKeysByOrchestratorAddress items."] pub async fn delegate_keys_by_orchestrator_address (& mut self , request : impl tonic :: IntoRequest < super :: QueryDelegateKeysByOrchestratorAddressRequest > ,) -> Result < tonic :: Response < super :: QueryDelegateKeysByOrchestratorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Query/DelegateKeysByOrchestratorAddress") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries a list of LastEventNonceByAddress items."] pub async fn last_event_nonce_by_address (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastEventNonceByAddressRequest > ,) -> Result < tonic :: Response < super :: QueryLastEventNonceByAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Query/LastEventNonceByAddress") ; self . inner . unary (request . into_request () , path , codec) . await } } }#[derive(Clone, PartialEq, ::prost::Message)]
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
pub struct MsgSetOrchestratorAddress {
    #[prost(string, tag="1")]
    pub validator: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub eth_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSetOrchestratorAddressResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs , clippy :: let_unit_value ,)] use tonic :: codegen :: * ; # [doc = " Msg defines the Msg service."] # [derive (Debug , Clone)] pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as Body > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor < F > (inner : T , interceptor : F) -> MsgClient < InterceptedService < T , F >> where F : tonic :: service :: Interceptor , T : tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody > , Response = http :: Response << T as tonic :: client :: GrpcService < tonic :: body :: BoxBody >> :: ResponseBody > > , < T as tonic :: codegen :: Service < http :: Request < tonic :: body :: BoxBody >> > :: Error : Into < StdError > + Send + Sync , { MsgClient :: new (InterceptedService :: new (inner , interceptor)) } # [doc = r" Compress requests with `gzip`."] # [doc = r""] # [doc = r" This requires the server to support it otherwise it might respond with an"] # [doc = r" error."] pub fn send_gzip (mut self) -> Self { self . inner = self . inner . send_gzip () ; self } # [doc = r" Enable decompressing responses with `gzip`."] pub fn accept_gzip (mut self) -> Self { self . inner = self . inner . accept_gzip () ; self } pub async fn ubt_deposited_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgUbtDepositedClaim > ,) -> Result < tonic :: Response < super :: MsgUbtDepositedClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Msg/UbtDepositedClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn set_orchestrator_address (& mut self , request : impl tonic :: IntoRequest < super :: MsgSetOrchestratorAddress > ,) -> Result < tonic :: Response < super :: MsgSetOrchestratorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/Baseledger.baseledgerbridge.baseledgerbridge.Msg/SetOrchestratorAddress") ; self . inner . unary (request . into_request () , path , codec) . await } } }