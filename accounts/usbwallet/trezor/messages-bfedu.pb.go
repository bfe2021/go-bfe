// This file originates from the SatoshiLabs Trezor `common` repository at:
//   https://github.com/trezor/trezor-common/blob/master/protob/messages-bfedu.proto
// dated 28.05.2019, commit 893fd219d4a01bcffa0cd9cfa631856371ec5aa9.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.3
// source: messages-bfedu.proto

package trezor

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//*
// Request: Ask device for public key corresponding to address_n path
// @start
// @next BfeduPublicKey
// @next Failure
type BfeduGetPublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressN    []uint32 `protobuf:"varint,1,rep,name=address_n,json=addressN" json:"address_n,omitempty"`          // BIP-32 path to derive the key from master node
	ShowDisplay *bool    `protobuf:"varint,2,opt,name=show_display,json=showDisplay" json:"show_display,omitempty"` // optionally show on display before sending the result
}

func (x *BfeduGetPublicKey) Reset() {
	*x = BfeduGetPublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduGetPublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduGetPublicKey) ProtoMessage() {}

func (x *BfeduGetPublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduGetPublicKey.ProtoReflect.Descriptor instead.
func (*BfeduGetPublicKey) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{0}
}

func (x *BfeduGetPublicKey) GetAddressN() []uint32 {
	if x != nil {
		return x.AddressN
	}
	return nil
}

func (x *BfeduGetPublicKey) GetShowDisplay() bool {
	if x != nil && x.ShowDisplay != nil {
		return *x.ShowDisplay
	}
	return false
}

//*
// Response: Contains public key derived from device private seed
// @end
type BfeduPublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *HDNodeType `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"` // BIP32 public node
	Xpub *string     `protobuf:"bytes,2,opt,name=xpub" json:"xpub,omitempty"` // serialized form of public node
}

func (x *BfeduPublicKey) Reset() {
	*x = BfeduPublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduPublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduPublicKey) ProtoMessage() {}

func (x *BfeduPublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduPublicKey.ProtoReflect.Descriptor instead.
func (*BfeduPublicKey) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{1}
}

func (x *BfeduPublicKey) GetNode() *HDNodeType {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *BfeduPublicKey) GetXpub() string {
	if x != nil && x.Xpub != nil {
		return *x.Xpub
	}
	return ""
}

//*
// Request: Ask device for Bfedu address corresponding to address_n path
// @start
// @next BfeduAddress
// @next Failure
type BfeduGetAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressN    []uint32 `protobuf:"varint,1,rep,name=address_n,json=addressN" json:"address_n,omitempty"`          // BIP-32 path to derive the key from master node
	ShowDisplay *bool    `protobuf:"varint,2,opt,name=show_display,json=showDisplay" json:"show_display,omitempty"` // optionally show on display before sending the result
}

func (x *BfeduGetAddress) Reset() {
	*x = BfeduGetAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduGetAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduGetAddress) ProtoMessage() {}

func (x *BfeduGetAddress) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduGetAddress.ProtoReflect.Descriptor instead.
func (*BfeduGetAddress) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{2}
}

func (x *BfeduGetAddress) GetAddressN() []uint32 {
	if x != nil {
		return x.AddressN
	}
	return nil
}

func (x *BfeduGetAddress) GetShowDisplay() bool {
	if x != nil && x.ShowDisplay != nil {
		return *x.ShowDisplay
	}
	return false
}

//*
// Response: Contains an Bfedu address derived from device private seed
// @end
type BfeduAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressBin []byte  `protobuf:"bytes,1,opt,name=addressBin" json:"addressBin,omitempty"` // Bfedu address as 20 bytes (legacy firmwares)
	AddressHex *string `protobuf:"bytes,2,opt,name=addressHex" json:"addressHex,omitempty"` // Bfedu address as hex string (newer firmwares)
}

func (x *BfeduAddress) Reset() {
	*x = BfeduAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduAddress) ProtoMessage() {}

func (x *BfeduAddress) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduAddress.ProtoReflect.Descriptor instead.
func (*BfeduAddress) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{3}
}

func (x *BfeduAddress) GetAddressBin() []byte {
	if x != nil {
		return x.AddressBin
	}
	return nil
}

func (x *BfeduAddress) GetAddressHex() string {
	if x != nil && x.AddressHex != nil {
		return *x.AddressHex
	}
	return ""
}

//*
// Request: Ask device to sign transaction
// All fields are optional from the protocol's point of view. Each field defaults to value `0` if missing.
// Note: the first at most 1024 bytes of data MUST be transmitted as part of this message.
// @start
// @next BfeduTxRequest
// @next Failure
type BfeduSignTx struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressN         []uint32 `protobuf:"varint,1,rep,name=address_n,json=addressN" json:"address_n,omitempty"`                          // BIP-32 path to derive the key from master node
	Nonce            []byte   `protobuf:"bytes,2,opt,name=nonce" json:"nonce,omitempty"`                                                 // <=256 bit unsigned big endian
	GasPrice         []byte   `protobuf:"bytes,3,opt,name=gas_price,json=gasPrice" json:"gas_price,omitempty"`                           // <=256 bit unsigned big endian (in wei)
	GasLimit         []byte   `protobuf:"bytes,4,opt,name=gas_limit,json=gasLimit" json:"gas_limit,omitempty"`                           // <=256 bit unsigned big endian
	ToBin            []byte   `protobuf:"bytes,5,opt,name=toBin" json:"toBin,omitempty"`                                                 // recipient address (20 bytes, legacy firmware)
	ToHex            *string  `protobuf:"bytes,11,opt,name=toHex" json:"toHex,omitempty"`                                                // recipient address (hex string, newer firmware)
	Value            []byte   `protobuf:"bytes,6,opt,name=value" json:"value,omitempty"`                                                 // <=256 bit unsigned big endian (in wei)
	DataInitialChunk []byte   `protobuf:"bytes,7,opt,name=data_initial_chunk,json=dataInitialChunk" json:"data_initial_chunk,omitempty"` // The initial data chunk (<= 1024 bytes)
	DataLength       *uint32  `protobuf:"varint,8,opt,name=data_length,json=dataLength" json:"data_length,omitempty"`                    // Length of transaction payload
	ChainId          *uint32  `protobuf:"varint,9,opt,name=chain_id,json=chainId" json:"chain_id,omitempty"`                             // Chain Id for EIP 155
	TxType           *uint32  `protobuf:"varint,10,opt,name=tx_type,json=txType" json:"tx_type,omitempty"`                               // (only for Wanchain)
}

func (x *BfeduSignTx) Reset() {
	*x = BfeduSignTx{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduSignTx) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduSignTx) ProtoMessage() {}

func (x *BfeduSignTx) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduSignTx.ProtoReflect.Descriptor instead.
func (*BfeduSignTx) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{4}
}

func (x *BfeduSignTx) GetAddressN() []uint32 {
	if x != nil {
		return x.AddressN
	}
	return nil
}

func (x *BfeduSignTx) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *BfeduSignTx) GetGasPrice() []byte {
	if x != nil {
		return x.GasPrice
	}
	return nil
}

func (x *BfeduSignTx) GetGasLimit() []byte {
	if x != nil {
		return x.GasLimit
	}
	return nil
}

func (x *BfeduSignTx) GetToBin() []byte {
	if x != nil {
		return x.ToBin
	}
	return nil
}

func (x *BfeduSignTx) GetToHex() string {
	if x != nil && x.ToHex != nil {
		return *x.ToHex
	}
	return ""
}

func (x *BfeduSignTx) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *BfeduSignTx) GetDataInitialChunk() []byte {
	if x != nil {
		return x.DataInitialChunk
	}
	return nil
}

func (x *BfeduSignTx) GetDataLength() uint32 {
	if x != nil && x.DataLength != nil {
		return *x.DataLength
	}
	return 0
}

func (x *BfeduSignTx) GetChainId() uint32 {
	if x != nil && x.ChainId != nil {
		return *x.ChainId
	}
	return 0
}

func (x *BfeduSignTx) GetTxType() uint32 {
	if x != nil && x.TxType != nil {
		return *x.TxType
	}
	return 0
}

//*
// Response: Device asks for more data from transaction payload, or returns the signature.
// If data_length is set, device awaits that many more bytes of payload.
// Otherwise, the signature_* fields contain the computed transaction signature. All three fields will be present.
// @end
// @next BfeduTxAck
type BfeduTxRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataLength *uint32 `protobuf:"varint,1,opt,name=data_length,json=dataLength" json:"data_length,omitempty"` // Number of bytes being requested (<= 1024)
	SignatureV *uint32 `protobuf:"varint,2,opt,name=signature_v,json=signatureV" json:"signature_v,omitempty"` // Computed signature (recovery parameter, limited to 27 or 28)
	SignatureR []byte  `protobuf:"bytes,3,opt,name=signature_r,json=signatureR" json:"signature_r,omitempty"`  // Computed signature R component (256 bit)
	SignatureS []byte  `protobuf:"bytes,4,opt,name=signature_s,json=signatureS" json:"signature_s,omitempty"`  // Computed signature S component (256 bit)
}

func (x *BfeduTxRequest) Reset() {
	*x = BfeduTxRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduTxRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduTxRequest) ProtoMessage() {}

func (x *BfeduTxRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduTxRequest.ProtoReflect.Descriptor instead.
func (*BfeduTxRequest) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{5}
}

func (x *BfeduTxRequest) GetDataLength() uint32 {
	if x != nil && x.DataLength != nil {
		return *x.DataLength
	}
	return 0
}

func (x *BfeduTxRequest) GetSignatureV() uint32 {
	if x != nil && x.SignatureV != nil {
		return *x.SignatureV
	}
	return 0
}

func (x *BfeduTxRequest) GetSignatureR() []byte {
	if x != nil {
		return x.SignatureR
	}
	return nil
}

func (x *BfeduTxRequest) GetSignatureS() []byte {
	if x != nil {
		return x.SignatureS
	}
	return nil
}

//*
// Request: Transaction payload data.
// @next BfeduTxRequest
type BfeduTxAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataChunk []byte `protobuf:"bytes,1,opt,name=data_chunk,json=dataChunk" json:"data_chunk,omitempty"` // Bytes from transaction payload (<= 1024 bytes)
}

func (x *BfeduTxAck) Reset() {
	*x = BfeduTxAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduTxAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduTxAck) ProtoMessage() {}

func (x *BfeduTxAck) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduTxAck.ProtoReflect.Descriptor instead.
func (*BfeduTxAck) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{6}
}

func (x *BfeduTxAck) GetDataChunk() []byte {
	if x != nil {
		return x.DataChunk
	}
	return nil
}

//*
// Request: Ask device to sign message
// @start
// @next BfeduMessageSignature
// @next Failure
type BfeduSignMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressN []uint32 `protobuf:"varint,1,rep,name=address_n,json=addressN" json:"address_n,omitempty"` // BIP-32 path to derive the key from master node
	Message  []byte   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`                    // message to be signed
}

func (x *BfeduSignMessage) Reset() {
	*x = BfeduSignMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduSignMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduSignMessage) ProtoMessage() {}

func (x *BfeduSignMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduSignMessage.ProtoReflect.Descriptor instead.
func (*BfeduSignMessage) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{7}
}

func (x *BfeduSignMessage) GetAddressN() []uint32 {
	if x != nil {
		return x.AddressN
	}
	return nil
}

func (x *BfeduSignMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

//*
// Response: Signed message
// @end
type BfeduMessageSignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressBin []byte  `protobuf:"bytes,1,opt,name=addressBin" json:"addressBin,omitempty"` // address used to sign the message (20 bytes, legacy firmware)
	Signature  []byte  `protobuf:"bytes,2,opt,name=signature" json:"signature,omitempty"`   // signature of the message
	AddressHex *string `protobuf:"bytes,3,opt,name=addressHex" json:"addressHex,omitempty"` // address used to sign the message (hex string, newer firmware)
}

func (x *BfeduMessageSignature) Reset() {
	*x = BfeduMessageSignature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduMessageSignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduMessageSignature) ProtoMessage() {}

func (x *BfeduMessageSignature) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduMessageSignature.ProtoReflect.Descriptor instead.
func (*BfeduMessageSignature) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{8}
}

func (x *BfeduMessageSignature) GetAddressBin() []byte {
	if x != nil {
		return x.AddressBin
	}
	return nil
}

func (x *BfeduMessageSignature) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *BfeduMessageSignature) GetAddressHex() string {
	if x != nil && x.AddressHex != nil {
		return *x.AddressHex
	}
	return ""
}

//*
// Request: Ask device to verify message
// @start
// @next Success
// @next Failure
type BfeduVerifyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressBin []byte  `protobuf:"bytes,1,opt,name=addressBin" json:"addressBin,omitempty"` // address to verify (20 bytes, legacy firmware)
	Signature  []byte  `protobuf:"bytes,2,opt,name=signature" json:"signature,omitempty"`   // signature to verify
	Message    []byte  `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`       // message to verify
	AddressHex *string `protobuf:"bytes,4,opt,name=addressHex" json:"addressHex,omitempty"` // address to verify (hex string, newer firmware)
}

func (x *BfeduVerifyMessage) Reset() {
	*x = BfeduVerifyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_bfedu_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BfeduVerifyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BfeduVerifyMessage) ProtoMessage() {}

func (x *BfeduVerifyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messages_bfedu_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BfeduVerifyMessage.ProtoReflect.Descriptor instead.
func (*BfeduVerifyMessage) Descriptor() ([]byte, []int) {
	return file_messages_bfedu_proto_rawDescGZIP(), []int{9}
}

func (x *BfeduVerifyMessage) GetAddressBin() []byte {
	if x != nil {
		return x.AddressBin
	}
	return nil
}

func (x *BfeduVerifyMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *BfeduVerifyMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *BfeduVerifyMessage) GetAddressHex() string {
	if x != nil && x.AddressHex != nil {
		return *x.AddressHex
	}
	return ""
}

var File_messages_bfedu_proto protoreflect.FileDescriptor

var file_messages_bfedu_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2d, 0x6f, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x68, 0x77, 0x2e, 0x74, 0x72, 0x65, 0x7a,
	0x6f, 0x72, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x6f, 0x72, 0x61, 0x6e,
	0x67, 0x65, 0x1a, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2d, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x12, 0x4f, 0x72, 0x61,
	0x6e, 0x67, 0x65, 0x47, 0x65, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x1b, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x6e, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0d, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x4e, 0x12, 0x21, 0x0a, 0x0c,
	0x73, 0x68, 0x6f, 0x77, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x77, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x22,
	0x60, 0x0a, 0x0f, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x12, 0x39, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x68, 0x77, 0x2e, 0x74, 0x72, 0x65, 0x7a, 0x6f, 0x72, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x48, 0x44, 0x4e,
	0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x78, 0x70, 0x75, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x78, 0x70, 0x75,
	0x62, 0x22, 0x52, 0x0a, 0x10, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x47, 0x65, 0x74, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x5f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x4e, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x77, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c,
	0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x77, 0x44, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x22, 0x4f, 0x0a, 0x0d, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x42, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x42, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x48, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x48, 0x65, 0x78, 0x22, 0xc0, 0x02, 0x0a, 0x0c, 0x4f, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x53, 0x69, 0x67, 0x6e, 0x54, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x5f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x4e, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61,
	0x73, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x67,
	0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x42, 0x69, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x74, 0x6f, 0x42, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x48, 0x65, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x48, 0x65, 0x78,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x10, 0x64, 0x61, 0x74, 0x61, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x4c,
	0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x74, 0x78, 0x54, 0x79, 0x70, 0x65, 0x22, 0x95, 0x01, 0x0a, 0x0f, 0x4f, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x54, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x1f,
	0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x76, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x56, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x53, 0x22, 0x2c, 0x0a, 0x0b, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x78, 0x41, 0x63, 0x6b,
	0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22,
	0x4a, 0x0a, 0x11, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f,
	0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x4e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x76, 0x0a, 0x16, 0x4f,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x42, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x42, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x48, 0x65,
	0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x48, 0x65, 0x78, 0x22, 0x8d, 0x01, 0x0a, 0x13, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x48, 0x65,
	0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x48, 0x65, 0x78, 0x42, 0x3e, 0x0a, 0x23, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x61, 0x74, 0x6f, 0x73,
	0x68, 0x69, 0x6c, 0x61, 0x62, 0x73, 0x2e, 0x74, 0x72, 0x65, 0x7a, 0x6f, 0x72, 0x2e, 0x6c, 0x69,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x42, 0x13, 0x54, 0x72, 0x65, 0x7a,
	0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5a,
	0x02, 0x2e, 0x2f,
}

var (
	file_messages_bfedu_proto_rawDescOnce sync.Once
	file_messages_bfedu_proto_rawDescData = file_messages_bfedu_proto_rawDesc
)

func file_messages_bfedu_proto_rawDescGZIP() []byte {
	file_messages_bfedu_proto_rawDescOnce.Do(func() {
		file_messages_bfedu_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_bfedu_proto_rawDescData)
	})
	return file_messages_bfedu_proto_rawDescData
}

var file_messages_bfedu_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_messages_bfedu_proto_goTypes = []interface{}{
	(*BfeduGetPublicKey)(nil),     // 0: hw.trezor.messages.bfedu.BfeduGetPublicKey
	(*BfeduPublicKey)(nil),        // 1: hw.trezor.messages.bfedu.BfeduPublicKey
	(*BfeduGetAddress)(nil),       // 2: hw.trezor.messages.bfedu.BfeduGetAddress
	(*BfeduAddress)(nil),          // 3: hw.trezor.messages.bfedu.BfeduAddress
	(*BfeduSignTx)(nil),           // 4: hw.trezor.messages.bfedu.BfeduSignTx
	(*BfeduTxRequest)(nil),        // 5: hw.trezor.messages.bfedu.BfeduTxRequest
	(*BfeduTxAck)(nil),            // 6: hw.trezor.messages.bfedu.BfeduTxAck
	(*BfeduSignMessage)(nil),      // 7: hw.trezor.messages.bfedu.BfeduSignMessage
	(*BfeduMessageSignature)(nil), // 8: hw.trezor.messages.bfedu.BfeduMessageSignature
	(*BfeduVerifyMessage)(nil),    // 9: hw.trezor.messages.bfedu.BfeduVerifyMessage
	(*HDNodeType)(nil),            // 10: hw.trezor.messages.common.HDNodeType
}
var file_messages_bfedu_proto_depIdxs = []int32{
	10, // 0: hw.trezor.messages.bfedu.BfeduPublicKey.node:type_name -> hw.trezor.messages.common.HDNodeType
	1,  // [1:1] is the sub-list for method output_type
	1,  // [1:1] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_messages_bfedu_proto_init() }
func file_messages_bfedu_proto_init() {
	if File_messages_bfedu_proto != nil {
		return
	}
	file_messages_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_messages_bfedu_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduGetPublicKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduPublicKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduGetAddress); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduAddress); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduSignTx); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduTxRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduTxAck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduSignMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduMessageSignature); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_bfedu_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BfeduVerifyMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messages_bfedu_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_bfedu_proto_goTypes,
		DependencyIndexes: file_messages_bfedu_proto_depIdxs,
		MessageInfos:      file_messages_bfedu_proto_msgTypes,
	}.Build()
	File_messages_bfedu_proto = out.File
	file_messages_bfedu_proto_rawDesc = nil
	file_messages_bfedu_proto_goTypes = nil
	file_messages_bfedu_proto_depIdxs = nil
}
