/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED BY codecgen.
// ************************************************************

package v1

import (
	"errors"
	"fmt"
	codec1978 "github.com/ugorji/go/codec"
	pkg3_v1 "k8s.io/api/core/v1"
	pkg1_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkg2_types "k8s.io/apimachinery/pkg/types"
	"reflect"
	"runtime"
	time "time"
)

const (
	// ----- content types ----
	codecSelferC_UTF81234 = 1
	codecSelferC_RAW1234  = 0
	// ----- value types used ----
	codecSelferValueTypeArray1234 = 10
	codecSelferValueTypeMap1234   = 9
	// ----- containerStateValues ----
	codecSelfer_containerMapKey1234    = 2
	codecSelfer_containerMapValue1234  = 3
	codecSelfer_containerMapEnd1234    = 4
	codecSelfer_containerArrayElem1234 = 6
	codecSelfer_containerArrayEnd1234  = 7
)

var (
	codecSelferBitsize1234                         = uint8(reflect.TypeOf(uint(0)).Bits())
	codecSelferOnlyMapOrArrayEncodeToStructErr1234 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer1234 struct{}

func init() {
	if codec1978.GenVersion != 5 {
		_, file, _, _ := runtime.Caller(0)
		err := fmt.Errorf("codecgen version mismatch: current: %v, need %v. Re-generate file: %v",
			5, codec1978.GenVersion, file)
		panic(err)
	}
	if false { // reference the types, but skip this branch at build/run time
		var v0 pkg3_v1.PersistentVolumeReclaimPolicy
		var v1 pkg1_v1.TypeMeta
		var v2 pkg2_types.UID
		var v3 time.Time
		_, _, _, _ = v0, v1, v2, v3
	}
}

func (x *StorageClass) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym1 := z.EncBinary()
		_ = yym1
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			var yyq2 [8]bool
			_, _, _ = yysep2, yyq2, yy2arr2
			const yyr2 bool = false
			yyq2[0] = x.Kind != ""
			yyq2[1] = x.APIVersion != ""
			yyq2[2] = true
			yyq2[4] = len(x.Parameters) != 0
			yyq2[5] = x.ReclaimPolicy != nil
			yyq2[6] = len(x.MountOptions) != 0
			yyq2[7] = x.AllowVolumeExpand != nil
			var yynn2 int
			if yyr2 || yy2arr2 {
				r.EncodeArrayStart(8)
			} else {
				yynn2 = 1
				for _, b := range yyq2 {
					if b {
						yynn2++
					}
				}
				r.EncodeMapStart(yynn2)
				yynn2 = 0
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[0] {
					yym4 := z.EncBinary()
					_ = yym4
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.Kind))
					}
				} else {
					r.EncodeString(codecSelferC_UTF81234, "")
				}
			} else {
				if yyq2[0] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("kind"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yym5 := z.EncBinary()
					_ = yym5
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.Kind))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[1] {
					yym7 := z.EncBinary()
					_ = yym7
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.APIVersion))
					}
				} else {
					r.EncodeString(codecSelferC_UTF81234, "")
				}
			} else {
				if yyq2[1] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("apiVersion"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yym8 := z.EncBinary()
					_ = yym8
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.APIVersion))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[2] {
					yy10 := &x.ObjectMeta
					yym11 := z.EncBinary()
					_ = yym11
					if false {
					} else if z.HasExtensions() && z.EncExt(yy10) {
					} else {
						z.EncFallback(yy10)
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[2] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("metadata"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yy12 := &x.ObjectMeta
					yym13 := z.EncBinary()
					_ = yym13
					if false {
					} else if z.HasExtensions() && z.EncExt(yy12) {
					} else {
						z.EncFallback(yy12)
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				yym15 := z.EncBinary()
				_ = yym15
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81234, string(x.Provisioner))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1234)
				r.EncodeString(codecSelferC_UTF81234, string("provisioner"))
				z.EncSendContainerState(codecSelfer_containerMapValue1234)
				yym16 := z.EncBinary()
				_ = yym16
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81234, string(x.Provisioner))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[4] {
					if x.Parameters == nil {
						r.EncodeNil()
					} else {
						yym18 := z.EncBinary()
						_ = yym18
						if false {
						} else {
							z.F.EncMapStringStringV(x.Parameters, false, e)
						}
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[4] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("parameters"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					if x.Parameters == nil {
						r.EncodeNil()
					} else {
						yym19 := z.EncBinary()
						_ = yym19
						if false {
						} else {
							z.F.EncMapStringStringV(x.Parameters, false, e)
						}
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[5] {
					if x.ReclaimPolicy == nil {
						r.EncodeNil()
					} else {
						yy21 := *x.ReclaimPolicy
						yysf22 := &yy21
						yysf22.CodecEncodeSelf(e)
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[5] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("reclaimPolicy"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					if x.ReclaimPolicy == nil {
						r.EncodeNil()
					} else {
						yy23 := *x.ReclaimPolicy
						yysf24 := &yy23
						yysf24.CodecEncodeSelf(e)
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[6] {
					if x.MountOptions == nil {
						r.EncodeNil()
					} else {
						yym26 := z.EncBinary()
						_ = yym26
						if false {
						} else {
							z.F.EncSliceStringV(x.MountOptions, false, e)
						}
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[6] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("mountOptions"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					if x.MountOptions == nil {
						r.EncodeNil()
					} else {
						yym27 := z.EncBinary()
						_ = yym27
						if false {
						} else {
							z.F.EncSliceStringV(x.MountOptions, false, e)
						}
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[7] {
					if x.AllowVolumeExpand == nil {
						r.EncodeNil()
					} else {
						yy29 := *x.AllowVolumeExpand
						yym30 := z.EncBinary()
						_ = yym30
						if false {
						} else {
							r.EncodeBool(bool(yy29))
						}
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[7] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("allowVolumeExpand"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					if x.AllowVolumeExpand == nil {
						r.EncodeNil()
					} else {
						yy31 := *x.AllowVolumeExpand
						yym32 := z.EncBinary()
						_ = yym32
						if false {
						} else {
							r.EncodeBool(bool(yy31))
						}
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1234)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1234)
			}
		}
	}
}

func (x *StorageClass) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap1234 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1234)
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray1234 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1234)
		}
	}
}

func (x *StorageClass) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys3Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys3Slc
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1234)
		yys3Slc = r.DecodeBytes(yys3Slc, true, true)
		yys3 := string(yys3Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1234)
		switch yys3 {
		case "kind":
			if r.TryDecodeAsNil() {
				x.Kind = ""
			} else {
				yyv4 := &x.Kind
				yym5 := z.DecBinary()
				_ = yym5
				if false {
				} else {
					*((*string)(yyv4)) = r.DecodeString()
				}
			}
		case "apiVersion":
			if r.TryDecodeAsNil() {
				x.APIVersion = ""
			} else {
				yyv6 := &x.APIVersion
				yym7 := z.DecBinary()
				_ = yym7
				if false {
				} else {
					*((*string)(yyv6)) = r.DecodeString()
				}
			}
		case "metadata":
			if r.TryDecodeAsNil() {
				x.ObjectMeta = pkg1_v1.ObjectMeta{}
			} else {
				yyv8 := &x.ObjectMeta
				yym9 := z.DecBinary()
				_ = yym9
				if false {
				} else if z.HasExtensions() && z.DecExt(yyv8) {
				} else {
					z.DecFallback(yyv8, false)
				}
			}
		case "provisioner":
			if r.TryDecodeAsNil() {
				x.Provisioner = ""
			} else {
				yyv10 := &x.Provisioner
				yym11 := z.DecBinary()
				_ = yym11
				if false {
				} else {
					*((*string)(yyv10)) = r.DecodeString()
				}
			}
		case "parameters":
			if r.TryDecodeAsNil() {
				x.Parameters = nil
			} else {
				yyv12 := &x.Parameters
				yym13 := z.DecBinary()
				_ = yym13
				if false {
				} else {
					z.F.DecMapStringStringX(yyv12, false, d)
				}
			}
		case "reclaimPolicy":
			if r.TryDecodeAsNil() {
				if x.ReclaimPolicy != nil {
					x.ReclaimPolicy = nil
				}
			} else {
				if x.ReclaimPolicy == nil {
					x.ReclaimPolicy = new(pkg3_v1.PersistentVolumeReclaimPolicy)
				}
				x.ReclaimPolicy.CodecDecodeSelf(d)
			}
		case "mountOptions":
			if r.TryDecodeAsNil() {
				x.MountOptions = nil
			} else {
				yyv15 := &x.MountOptions
				yym16 := z.DecBinary()
				_ = yym16
				if false {
				} else {
					z.F.DecSliceStringX(yyv15, false, d)
				}
			}
		case "allowVolumeExpand":
			if r.TryDecodeAsNil() {
				if x.AllowVolumeExpand != nil {
					x.AllowVolumeExpand = nil
				}
			} else {
				if x.AllowVolumeExpand == nil {
					x.AllowVolumeExpand = new(bool)
				}
				yym18 := z.DecBinary()
				_ = yym18
				if false {
				} else {
					*((*bool)(x.AllowVolumeExpand)) = r.DecodeBool()
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	z.DecSendContainerState(codecSelfer_containerMapEnd1234)
}

func (x *StorageClass) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj19 int
	var yyb19 bool
	var yyhl19 bool = l >= 0
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.Kind = ""
	} else {
		yyv20 := &x.Kind
		yym21 := z.DecBinary()
		_ = yym21
		if false {
		} else {
			*((*string)(yyv20)) = r.DecodeString()
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.APIVersion = ""
	} else {
		yyv22 := &x.APIVersion
		yym23 := z.DecBinary()
		_ = yym23
		if false {
		} else {
			*((*string)(yyv22)) = r.DecodeString()
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.ObjectMeta = pkg1_v1.ObjectMeta{}
	} else {
		yyv24 := &x.ObjectMeta
		yym25 := z.DecBinary()
		_ = yym25
		if false {
		} else if z.HasExtensions() && z.DecExt(yyv24) {
		} else {
			z.DecFallback(yyv24, false)
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.Provisioner = ""
	} else {
		yyv26 := &x.Provisioner
		yym27 := z.DecBinary()
		_ = yym27
		if false {
		} else {
			*((*string)(yyv26)) = r.DecodeString()
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.Parameters = nil
	} else {
		yyv28 := &x.Parameters
		yym29 := z.DecBinary()
		_ = yym29
		if false {
		} else {
			z.F.DecMapStringStringX(yyv28, false, d)
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		if x.ReclaimPolicy != nil {
			x.ReclaimPolicy = nil
		}
	} else {
		if x.ReclaimPolicy == nil {
			x.ReclaimPolicy = new(pkg3_v1.PersistentVolumeReclaimPolicy)
		}
		x.ReclaimPolicy.CodecDecodeSelf(d)
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.MountOptions = nil
	} else {
		yyv31 := &x.MountOptions
		yym32 := z.DecBinary()
		_ = yym32
		if false {
		} else {
			z.F.DecSliceStringX(yyv31, false, d)
		}
	}
	yyj19++
	if yyhl19 {
		yyb19 = yyj19 > l
	} else {
		yyb19 = r.CheckBreak()
	}
	if yyb19 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		if x.AllowVolumeExpand != nil {
			x.AllowVolumeExpand = nil
		}
	} else {
		if x.AllowVolumeExpand == nil {
			x.AllowVolumeExpand = new(bool)
		}
		yym34 := z.DecBinary()
		_ = yym34
		if false {
		} else {
			*((*bool)(x.AllowVolumeExpand)) = r.DecodeBool()
		}
	}
	for {
		yyj19++
		if yyhl19 {
			yyb19 = yyj19 > l
		} else {
			yyb19 = r.CheckBreak()
		}
		if yyb19 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1234)
		z.DecStructFieldNotFound(yyj19-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
}

func (x *StorageClassList) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym1 := z.EncBinary()
		_ = yym1
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			var yyq2 [4]bool
			_, _, _ = yysep2, yyq2, yy2arr2
			const yyr2 bool = false
			yyq2[0] = x.Kind != ""
			yyq2[1] = x.APIVersion != ""
			yyq2[2] = true
			var yynn2 int
			if yyr2 || yy2arr2 {
				r.EncodeArrayStart(4)
			} else {
				yynn2 = 1
				for _, b := range yyq2 {
					if b {
						yynn2++
					}
				}
				r.EncodeMapStart(yynn2)
				yynn2 = 0
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[0] {
					yym4 := z.EncBinary()
					_ = yym4
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.Kind))
					}
				} else {
					r.EncodeString(codecSelferC_UTF81234, "")
				}
			} else {
				if yyq2[0] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("kind"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yym5 := z.EncBinary()
					_ = yym5
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.Kind))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[1] {
					yym7 := z.EncBinary()
					_ = yym7
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.APIVersion))
					}
				} else {
					r.EncodeString(codecSelferC_UTF81234, "")
				}
			} else {
				if yyq2[1] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("apiVersion"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yym8 := z.EncBinary()
					_ = yym8
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81234, string(x.APIVersion))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if yyq2[2] {
					yy10 := &x.ListMeta
					yym11 := z.EncBinary()
					_ = yym11
					if false {
					} else if z.HasExtensions() && z.EncExt(yy10) {
					} else {
						z.EncFallback(yy10)
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[2] {
					z.EncSendContainerState(codecSelfer_containerMapKey1234)
					r.EncodeString(codecSelferC_UTF81234, string("metadata"))
					z.EncSendContainerState(codecSelfer_containerMapValue1234)
					yy12 := &x.ListMeta
					yym13 := z.EncBinary()
					_ = yym13
					if false {
					} else if z.HasExtensions() && z.EncExt(yy12) {
					} else {
						z.EncFallback(yy12)
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1234)
				if x.Items == nil {
					r.EncodeNil()
				} else {
					yym15 := z.EncBinary()
					_ = yym15
					if false {
					} else {
						h.encSliceStorageClass(([]StorageClass)(x.Items), e)
					}
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1234)
				r.EncodeString(codecSelferC_UTF81234, string("items"))
				z.EncSendContainerState(codecSelfer_containerMapValue1234)
				if x.Items == nil {
					r.EncodeNil()
				} else {
					yym16 := z.EncBinary()
					_ = yym16
					if false {
					} else {
						h.encSliceStorageClass(([]StorageClass)(x.Items), e)
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1234)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1234)
			}
		}
	}
}

func (x *StorageClassList) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap1234 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1234)
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray1234 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1234)
		}
	}
}

func (x *StorageClassList) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys3Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys3Slc
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1234)
		yys3Slc = r.DecodeBytes(yys3Slc, true, true)
		yys3 := string(yys3Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1234)
		switch yys3 {
		case "kind":
			if r.TryDecodeAsNil() {
				x.Kind = ""
			} else {
				yyv4 := &x.Kind
				yym5 := z.DecBinary()
				_ = yym5
				if false {
				} else {
					*((*string)(yyv4)) = r.DecodeString()
				}
			}
		case "apiVersion":
			if r.TryDecodeAsNil() {
				x.APIVersion = ""
			} else {
				yyv6 := &x.APIVersion
				yym7 := z.DecBinary()
				_ = yym7
				if false {
				} else {
					*((*string)(yyv6)) = r.DecodeString()
				}
			}
		case "metadata":
			if r.TryDecodeAsNil() {
				x.ListMeta = pkg1_v1.ListMeta{}
			} else {
				yyv8 := &x.ListMeta
				yym9 := z.DecBinary()
				_ = yym9
				if false {
				} else if z.HasExtensions() && z.DecExt(yyv8) {
				} else {
					z.DecFallback(yyv8, false)
				}
			}
		case "items":
			if r.TryDecodeAsNil() {
				x.Items = nil
			} else {
				yyv10 := &x.Items
				yym11 := z.DecBinary()
				_ = yym11
				if false {
				} else {
					h.decSliceStorageClass((*[]StorageClass)(yyv10), d)
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	z.DecSendContainerState(codecSelfer_containerMapEnd1234)
}

func (x *StorageClassList) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj12 int
	var yyb12 bool
	var yyhl12 bool = l >= 0
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.Kind = ""
	} else {
		yyv13 := &x.Kind
		yym14 := z.DecBinary()
		_ = yym14
		if false {
		} else {
			*((*string)(yyv13)) = r.DecodeString()
		}
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.APIVersion = ""
	} else {
		yyv15 := &x.APIVersion
		yym16 := z.DecBinary()
		_ = yym16
		if false {
		} else {
			*((*string)(yyv15)) = r.DecodeString()
		}
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.ListMeta = pkg1_v1.ListMeta{}
	} else {
		yyv17 := &x.ListMeta
		yym18 := z.DecBinary()
		_ = yym18
		if false {
		} else if z.HasExtensions() && z.DecExt(yyv17) {
		} else {
			z.DecFallback(yyv17, false)
		}
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1234)
	if r.TryDecodeAsNil() {
		x.Items = nil
	} else {
		yyv19 := &x.Items
		yym20 := z.DecBinary()
		_ = yym20
		if false {
		} else {
			h.decSliceStorageClass((*[]StorageClass)(yyv19), d)
		}
	}
	for {
		yyj12++
		if yyhl12 {
			yyb12 = yyj12 > l
		} else {
			yyb12 = r.CheckBreak()
		}
		if yyb12 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1234)
		z.DecStructFieldNotFound(yyj12-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1234)
}

func (x codecSelfer1234) encSliceStorageClass(v []StorageClass, e *codec1978.Encoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem1234)
		yy2 := &yyv1
		yy2.CodecEncodeSelf(e)
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd1234)
}

func (x codecSelfer1234) decSliceStorageClass(v *[]StorageClass, d *codec1978.Decoder) {
	var h codecSelfer1234
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []StorageClass{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else if yyl1 > 0 {
		var yyrr1, yyrl1 int
		var yyrt1 bool
		_, _ = yyrl1, yyrt1
		yyrr1 = yyl1 // len(yyv1)
		if yyl1 > cap(yyv1) {

			yyrg1 := len(yyv1) > 0
			yyv21 := yyv1
			yyrl1, yyrt1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 328)
			if yyrt1 {
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]StorageClass, yyrl1)
				}
			} else {
				yyv1 = make([]StorageClass, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			if yyrg1 {
				copy(yyv1, yyv21)
			}
		} else if yyl1 != len(yyv1) {
			yyv1 = yyv1[:yyl1]
			yyc1 = true
		}
		yyj1 := 0
		for ; yyj1 < yyrr1; yyj1++ {
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				yyv1[yyj1] = StorageClass{}
			} else {
				yyv2 := &yyv1[yyj1]
				yyv2.CodecDecodeSelf(d)
			}

		}
		if yyrt1 {
			for ; yyj1 < yyl1; yyj1++ {
				yyv1 = append(yyv1, StorageClass{})
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					yyv1[yyj1] = StorageClass{}
				} else {
					yyv3 := &yyv1[yyj1]
					yyv3.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj1 := 0
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, StorageClass{}) // var yyz1 StorageClass
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					yyv1[yyj1] = StorageClass{}
				} else {
					yyv4 := &yyv1[yyj1]
					yyv4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []StorageClass{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}
