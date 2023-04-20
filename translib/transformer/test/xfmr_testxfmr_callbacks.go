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
        XlateFuncBind("test_pre_xfmr", test_pre_xfmr)
        XlateFuncBind("test_post_xfmr", test_post_xfmr)
	XlateFuncBind("testsensor_type_tbl_xfmr", testsensor_type_tbl_xfmr)
	XlateFuncBind("YangToDb_testsensor_type_key_xfmr", YangToDb_testsensor_type_key_xfmr)
	XlateFuncBind("DbToYang_testsensor_type_key_xfmr", DbToYang_testsensor_type_key_xfmr)
	XlateFuncBind("YangToDb_test_set_key_xfmr", YangToDb_test_set_key_xfmr)
	XlateFuncBind("DbToYang_test_set_key_xfmr", DbToYang_test_set_key_xfmr)


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

var test_pre_xfmr PreXfmrFunc = func(inParams XfmrParams) error {
        var err error
        requestUriPath := (NewPathInfo(inParams.requestUri)).YangPath
       	if log.V(3) {
		log.Info("test_pre_xfmr:- Request URI path = ", requestUriPath)
	}	
	return err
}

var intf_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {

        requestUriPath := (NewPathInfo(inParams.requestUri)).YangPath
        retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
        log.Info("Entering intf_post_xfmr requestUriPath = ", requestUriPath)
        return retDbDataMap, nil
}

var testsensor_type_tbl_xfmr TableXfmrFunc = func(inParams XfmrParams) ([]string, error) {
	var tblList []string
        pathInfo := NewPathInfo(inParams.uri)
        groupId := pathInfo.Var("id")
        sensorType := pathInfo.Var("type")

	log.Info("testsensor_type_tbl_xfmr inParams.uri ", inParams.uri)

	if len(sensorType) == 0 {
		if inParams.oper == GET || inParams.oper == DELETE {
			tblList = append(tblList, "TEST_SENSOR_A_TABLE")
			tblList = append(tblList, "TEST_SENSOR_B_TABLE")
		}
	} else {
		if strings.HasPrefix(name, "sensora_") {
			tblList = append(tblList, "TEST_SENSOR_A_TABLE")
		} else if strings.HasPrefix(name, "sensorb_") {
			tblList = append(tblList, "TEST_SENSOR_B_TABLE")
		}
	}
        log.Info("testsensor_type_tbl_xfmr tblList= ", tblList)
	return tblList
}

var YangToDb_testsensor_type_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) 
        var sensor_type_key string
        var err error

        log.Info("YangToDb_testsensor_type_key_xfmr - inParams.uri ", inParams.uri)

        pathInfo := NewPathInfo(inParams.uri)
        groupId := pathInfo.Var("id")
        sensorType := pathInfo.Var("type")
        if groupId == "" || sensorType == "" {
                return sensor_type_key, err
        }
	if len(groupId) > 0 {
		sensor_type := ""
		if strings.HasPrefix(sensorType, "sensora_") {
			sensor_type = strings.Replace(sensorType, "sensora_", "sensor_type_a_", 1)
			sensor_type_key = groupId + "|" + sensor_type
		} else if strings.HasPrefix(sensorType, "sensorb_") {
			sensor_type = strings.Replace(sensorType, "sensorb_", "sensor_type_b_", 1)
			sensor_type_key = groupId + "|" + sensor_type
		} else {
			err_str := "Invalid key. Key not supported."
			err = tlerr.NotSupported(err_str)
		}
	}
        log.Info("YangToDb_testsensor_type_key_xfmr returns", sensor_type_key)
        return sensor_type_key, err
}

var DbToYang_testsensor_type_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
        var err error
        if log.V(3) {
                log.Info("Entering DbToYang_testsensor_type_key_xfmr inParams.uri ", inParams.uri)
        }
        var sensorType string

        if strings.Contains(inParams.key, "|") {
                key_split := strings.Split(inParams.key, "|")
                sensorType = key_split[1]
		if strings.HasPrefix(sensorType, "sensor_type_a_") {
			sensorType = strings.Replace(sensorType, "sensor_type_a_", "sensora_", 1)
		} else if strings.HasPrefix(sensorType, "sensor_type_b_") {
			sensorType = strings.Replace(sensorType, "sensor_type_b_", "sensorb_", 1)
		} else {
			sensorType = ""
			err_str := "Invalid key. Key not supported."
			err = tlerr.NotSupported(err_str)
		}
        }

        rmap["type"] = sensorType

        log.Info("DbToYang_testsensor_type_key_xfmr rmap ", rmap)
        return rmap, err
}

var YangToDb_test_set_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	testSetKey := ""
        if log.V(3) {
                log.Info("Entering DbToYang_testsensor_type_key_xfmr inParams.uri ", inParams.uri)
        }

        pathInfo := NewPathInfo(inParams.uri)
        testSetName := pathInfo.Var("name")
        testSetType := pathInfo.Var("type")

	if len(testSetName) > 0 && len(testSetType) > 0 {
		testSetKey = testSetName
	}
        log.Info(" YangToDb_test_set_key_xfmr returns ", testSetKey)
	return testSetKey, nil

}

var DbToYang_test_set_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
        rmap := make(map[string]interface{})
        var err error
        if log.V(3) {
		log.Info("DbToYang_test_set_key_xfmr invoked for uri: ", inParams.uri)
        }
        var testSetName string

        if len(inParams.key) > 0 && strings.Contains(inParams.key, "|") {
                key_split := strings.Split(inParams.key, "|")
                testSetName = key_split[0]
        } else {
                return result, errors.New("Incorrect dbKey : " + inParams.key)
	}

	data := (*inParams.dbDataMap)[inParams.curDb]
	tblName := TEST_SET_TABLE 
        if _, ok := data[tblName]; !ok {
                log.Info("DbToYang_test_set_key_xfmr table not found : ", tblName)
                return result, errors.New("table not found : " + tblName)
        }

        tsTbl := data[tblName]
        if _, ok := tsTbl[inParams.key]; !ok {
                log.Info("DbToYang_test_set_key_xfmr instance not found : ", inParams.key)
                return result, errors.New("Table Instance not found : " + inParams.key)
        }
        tsInst := tsTbl[inParams.key]
        typeStr, ok := tsInst.Field["type"]
        if ok {
        	var typeVal ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE
		if typeStr == "IPV4" {
			typeVal = ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4
		} else if typeStr == "IPV6" {
			typeVal = ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV6
		}
		result["type"] = ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE.ΛMap(typeVal)["ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE"][int64(typeVal)].Name
        }	
        rmap["name"] = testSetName

        log.Info("DbToYang_testsensor_type_key_xfmr rmap ", rmap)
        return rmap, err

}

func getTestSetRoot(s *ygot.GoStruct) *ocbinds.OpenconfigTestXfmr_TestXfmr {
        deviceObj := (*s).(*ocbinds.Device)
        return deviceObj.TestXfmr
}

func getTestSetKeyStrFromOCKey(setname string, settype ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE) string {
        //setT := settype.Map()["E_OpenconfigTestXfmr_TEST_SET_TYPE"][int64(settype)].Name
        setT := ""
        if settype == ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4 {
                setT = "TEST_SET_IPV4"
        } else {
                setT = "TEST_SET_IPV6"
        }

        return setname + "_" + setT
}

func getTestSetNameCompFromDbKey(testSetDbKey string, testSetType string) string {
        return testSetDbKey[:strings.LastIndex(testSetDbKey, "_" + testSetType)]
}

var YangToDb_exclude_filter_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
        res_map := make(map[string]string)
        var err error
        if inParams.param == nil {
            res_map["exclude-filter"] = ""
            return res_map, err
        }
        exflt, _ := inParams.param.(*string)
        if exflt != nil {
                res_map["exclude-filter"] = "filter_" + *exflt
                log.Info("YangToDb_exclude_filter_field_xfmr ", res_map["exclude-filter"])
        }
        return res_map, err

}

var DbToYang_exclude_filter_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
        var err error
        result := make(map[string]interface{})
        data := (*inParams.dbDataMap)[inParams.curDb]
        log.Info("DbToYang_exclude_filter_field_xfmr", data, inParams.ygRoot)
        pathInfo := NewPathInfo(inParams.uri)
        sensor_type := pathInfo.Var("type")
        tblNm := ""
        if strings.HasPrefix(sensor_type, "sensora") {
                tblNm = "TEST_SENSOR_A_TABLE"
        } else if strings.HasPrefix(sensor_type, "sensorb") {
                tblNm = "TEST_SENSOR_B_TABLE"
        }

        sensorData, ok := data[tblNm]
        if ok {
                sensorInst, instOk := sensorData[inParams.key]
                if instOk {
                        exFlt, fldOk := sensorInst.Field["exclude-filter"]
                        if fldOk {
                                result["exclude-filter"] = strings.Split(exFlt, "filter_")
                                log.Info("DbToYang_exclude_filter_field_xfmr - returning %v", result["exclude-filter"])
                        } else {
                                return nil, tlerr.NotFound("Resource Not Found")
                        }
                } else {
                        log.Info("DbToYang_exclude_filter_field_xfmr - sensor instance %v doesn't exist", inParams.key)
                }
        } else {
                log.Info("DbToYang_exclude_filter_field_xfmr - Table %v does exist in db Data", tblNm)
        }

        return result, err
}

var YangToDb_test_set_type_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
        res_map := make(map[string]string)
        var err error
        if inParams.param == nil {
            res_map[TEST_SET_TYPE] = ""
            return res_map, err
        }

        testSetType, _ := inParams.param.(ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE)
        log.Info("YangToDb_test_set_type_field_xfmr: ", inParams.ygRoot, " Xpath: ", inParams.uri, " Type: ", testSetType)
        res_map[TEST_SET_TYPE] = findInMap(TEST_SET_TYPE_MAP, strconv.FormatInt(int64(testSetType), 10))
        return res_map, err
}

var DbToYang_test_set_type_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
        var err error
        result := make(map[string]interface{})
        data := (*inParams.dbDataMap)[inParams.curDb]
        log.Info("DbToYang_test_set_type_field_xfmr", data, inParams.ygRoot)
        oc_testSetType := findInMap(TEST_SET_TYPE_MAP, data[TEST_SET_TABLE][inParams.key].Field[TEST_SET_TYPE])
        n, err := strconv.ParseInt(oc_testSetType, 10, 64)
        if n == int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4) {
             result[TEST_SET_TYPE] = "TEST_SET_IPV4"
        } else {
             result[TEST_SET_TYPE] = "TEST_SET_IPV6"
        }
        return result, err
}

var YangToDb_test_port_bindings_xfmr SubTreeXfmrYangToDb = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {
        var err error
        res_map := make(map[string]map[string]db.Value)
        testSetTableMap := make(map[string]db.Value)
        testSetTableMapNew := make(map[string]db.Value)
        log.Info("YangToDb_test_port_bindings_xfmr: ", inParams.ygRoot, inParams.uri)

        testXfmrObj := getTestSetRoot(inParams.ygRoot)
        if testXfmrObj.Interfaces == nil {
                return res_map, err
        }

        testSetTs := &db.TableSpec{Name: TEST_SET_TABLE}
        testSetKeys, err := inParams.d.GetKeys(testSetTs)
        if err != nil {
            return  res_map, err
        }

        for key := range testSetKeys {
                testSetEntry, err := inParams.d.GetEntry(testSetTs, testSetKeys[key])
                if err != nil {
                        return res_map, err
                }
                testSetTableMap[(testSetKeys[key].Get(0))] = testSetEntry
        }


        testSetInterfacesMap := make(map[string][]string)
        for intfId, _ := range testXfmrObj.Interfaces.Interface {
                intf := testXfmrObj.Interfaces.Interface[intfId]
                if intf != nil {
                        if intf.IngressTestSets != nil && len(intf.IngressTestSets.IngressTestSet) > 0 {
                                for inTestSetKey, _ := range intf.IngressTestSets.IngressTestSet {
                                        testSetName := getTestSetKeyStrFromOCKey(inTestSetKey.SetName, inTestSetKey.Type)
                                        testSetInterfacesMap[testSetName] = append(testSetInterfacesMap[testSetName], *intf.Id)
                                        _, ok := testSetTableMap[testSetName]
                                        if !ok {
                                                if inParams.oper == DELETE {
                                                        return res_map, tlerr.NotFound("Binding not found for test set  %v on %v", inTestSetKey.SetName, *intf.Id)
                                                }
                                        }
                                        testSetTableMapNew[testSetName] = db.Value{Field: make(map[string]string)}
                                }
                        } else {
                                for testSetKey, testSetData := range testSetTableMap {
                                        ports := testSetData.GetList(TEST_SET_PORTS)
                                        if contains(ports, *intf.Id) {
                                                testSetInterfacesMap[testSetKey] = append(testSetInterfacesMap[testSetKey], *intf.Id)
                                                testSetTableMapNew[testSetKey] = db.Value{Field: make(map[string]string)}
                                        }

                                }

                        }
                }
        }
        for k, _ := range testSetInterfacesMap {
                val := testSetTableMapNew[k]
                (&val).SetList(TEST_SET_PORTS + "@", testSetInterfacesMap[k])
        }
        res_map[TEST_SET_TABLE] = testSetTableMapNew
        if inParams.invokeCRUSubtreeOnce != nil {
                *inParams.invokeCRUSubtreeOnce = true
        }
        return res_map, err
}


