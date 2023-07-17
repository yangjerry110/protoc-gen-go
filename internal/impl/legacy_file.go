// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"

	"github.com/yangjerry110/protoc-gen-go/internal/filedesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Every enum and message type generated by protoc-gen-go since commit 2fc053c5
// on February 25th, 2016 has had a method to get the raw descriptor.
// Types that were not generated by protoc-gen-go or were generated prior
// to that version are not supported.
//
// The []byte returned is the encoded form of a FileDescriptorProto message
// compressed using GZIP. The []int is the path from the top-level file
// to the specific message or enum declaration.
type (
	enumV1 interface {
		EnumDescriptor() ([]byte, []int)
	}
	messageV1 interface {
		Descriptor() ([]byte, []int)
	}
)

var legacyFileDescCache sync.Map // map[*byte]protoreflect.FileDescriptor

// legacyLoadFileDesc unmarshals b as a compressed FileDescriptorProto message.
//
// This assumes that b is immutable and that b does not refer to part of a
// concatenated series of GZIP files (which would require shenanigans that
// rely on the concatenation properties of both protobufs and GZIP).
// File descriptors generated by protoc-gen-go do not rely on that property.
func legacyLoadFileDesc(b []byte) protoreflect.FileDescriptor {
	// Fast-path: check whether we already have a cached file descriptor.
	if fd, ok := legacyFileDescCache.Load(&b[0]); ok {
		return fd.(protoreflect.FileDescriptor)
	}

	// Slow-path: decompress and unmarshal the file descriptor proto.
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	b2, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	fd := filedesc.Builder{
		RawDescriptor: b2,
		FileRegistry:  resolverOnly{protoregistry.GlobalFiles}, // do not register back to global registry
	}.Build().File
	if fd, ok := legacyFileDescCache.LoadOrStore(&b[0], fd); ok {
		return fd.(protoreflect.FileDescriptor)
	}
	return fd
}

type resolverOnly struct {
	reg *protoregistry.Files
}

func (r resolverOnly) FindFileByPath(path string) (protoreflect.FileDescriptor, error) {
	return r.reg.FindFileByPath(path)
}
func (r resolverOnly) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	return r.reg.FindDescriptorByName(name)
}
func (resolverOnly) RegisterFile(protoreflect.FileDescriptor) error {
	return nil
}