////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2023 Dell, Inc.                                                 //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//  http://www.apache.org/licenses/LICENSE-2.0                                //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

//go:build !campus_pkg

package transformer

import (
        "strconv"
        "strings"
	"github.com/openconfig/ygot/ygot"
        "github.com/Azure/sonic-mgmt-common/translib/db"
        "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
        log "github.com/golang/glog"
)

func init() {
	XlateFuncBind("YangToDb_exclude_filter_field_xfmr", YangToDb_exclude_filter_field_xfmr)
	XlateFuncBind("DbToYang_exclude_filter_field_xfmr", DbToYang_exclude_filter_field_xfmr)
	XlateFuncBind("YangToDb_type_field_xfmr", YangToDb_type_field_xfmr)
	XlateFuncBind("DbToYang_type_field_xfmr", DbToYang_type_field_xfmr)
	XlateFuncBind("YangToDb_test_port_bindings_xfmr", YangToDb_test_port_bindings_xfmr)
	XlateFuncBind("DbToYang_test_port_bindings_xfmr", DbToYang_test_port_bindings_xfmr)
}

const (
	TEST_SET_TABLE = "TEST_SET_TABLE"
	TEST_SET_TYPE = "type"
)

/* E_OpenconfigTestXfmr_TEST_SET_TYPE */
var TEST_SET_TYPE_MAP = map[string]string{
	strconv.FormatInt(int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4), 10): "L3",
	strconv.FormatInt(int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV6), 10): "L3V6",
}

