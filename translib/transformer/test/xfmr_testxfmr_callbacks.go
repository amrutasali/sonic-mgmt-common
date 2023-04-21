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
	"errors"
	"reflect"
	"github.com/openconfig/ygot/ygot"
        "github.com/Azure/sonic-mgmt-common/translib/db"
        "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
        log "github.com/golang/glog"
)

func init() {
	// Pre and Post transformer functions
        XlateFuncBind("test_pre_xfmr", test_pre_xfmr)
        XlateFuncBind("test_post_xfmr", test_post_xfmr)

	// Table transformer functions
	XlateFuncBind("testsensor_type_tbl_xfmr", testsensor_type_tbl_xfmr)

	// Key transformer functions
	XlateFuncBind("YangToDb_testsensor_type_key_xfmr", YangToDb_testsensor_type_key_xfmr)
	XlateFuncBind("DbToYang_testsensor_type_key_xfmr", DbToYang_testsensor_type_key_xfmr)
	XlateFuncBind("YangToDb_test_set_key_xfmr", YangToDb_test_set_key_xfmr)
	XlateFuncBind("DbToYang_test_set_key_xfmr", DbToYang_test_set_key_xfmr)

	// Key leafrefed Field transformer functions
	XlateFuncBind("DbToYang_test_sensor_group_id_field_xfmr", DbToYang_test_sensor_group_id_field_xfmr)
	XlateFuncBind("DbToYang_test_sensor_type_field_xfmr", DbToYang_test_sensor_type_field_xfmr)
	XlateFuncBind("DbToYang_test_set_name_field_xfmr", DbToYang_test_set_name_field_xfmr)

	// Field transformer functions
	XlateFuncBind("YangToDb_exclude_filter_field_xfmr", YangToDb_exclude_filter_field_xfmr)
	XlateFuncBind("DbToYang_exclude_filter_field_xfmr", DbToYang_exclude_filter_field_xfmr)
	XlateFuncBind("YangToDb_test_set_type_field_xfmr", YangToDb_test_set_type_field_xfmr)
	XlateFuncBind("DbToYang_test_set_type_field_xfmr", DbToYang_test_set_type_field_xfmr)

	//Subtree transformer function
	XlateFuncBind("YangToDb_test_port_bindings_xfmr", YangToDb_test_port_bindings_xfmr)
	XlateFuncBind("DbToYang_test_port_bindings_xfmr", DbToYang_test_port_bindings_xfmr)
	XlateFuncBind("Subscribe_test_port_bindings_xfmr", Subscribe_test_port_bindings_xfmr)
}

const (
	TEST_SET_TABLE = "TEST_SET_TABLE"
	TEST_SET_TYPE = "type"
	TEST_SET_PORTS = "ports"
)

/* E_OpenconfigTestXfmr_TEST_SET_TYPE */
var TEST_SET_TYPE_MAP = map[string]string{
	strconv.FormatInt(int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4), 10): "L3",
	strconv.FormatInt(int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV6), 10): "L3V6",
}

var test_pre_xfmr PreXfmrFunc = func(inParams XfmrParams) error {
        var err error
	log.Info("test_pre_xfmr:- Request URI path = ", inParams.requestUri)
	return err
}

var test_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {

        retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
        log.Info("Entering intf_post_xfmr requestUriPath = ", inParams.requestUri)
        return retDbDataMap, nil
}

var testsensor_type_tbl_xfmr TableXfmrFunc = func(inParams XfmrParams) ([]string, error) {
	var tblList []string
        pathInfo := NewPathInfo(inParams.uri)
        groupId := pathInfo.Var("id")
        sensorType := pathInfo.Var("type")

	log.Info("testsensor_type_tbl_xfmr inParams.uri ", inParams.uri)

	if len(groupId) == 0 {
		return tblList, nil
	}
	if len(sensorType) == 0 {
		if inParams.oper == GET || inParams.oper == DELETE {
			tblList = append(tblList, "TEST_SENSOR_A_TABLE")
			tblList = append(tblList, "TEST_SENSOR_B_TABLE")
		}
	} else {
		if strings.HasPrefix(sensorType, "sensora_") {
			tblList = append(tblList, "TEST_SENSOR_A_TABLE")
		} else if strings.HasPrefix(sensorType, "sensorb_") {
			tblList = append(tblList, "TEST_SENSOR_B_TABLE")
		}
	}
        log.Info("testsensor_type_tbl_xfmr tblList= ", tblList)
	return tblList, nil
}

var YangToDb_testsensor_type_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
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
                return rmap, errors.New("Incorrect dbKey : " + inParams.key)
	}

	data := (*inParams.dbDataMap)[inParams.curDb]
	tblName := TEST_SET_TABLE 
        if _, ok := data[tblName]; !ok {
                log.Info("DbToYang_test_set_key_xfmr table not found : ", tblName)
                return rmap, errors.New("table not found : " + tblName)
        }

        tsTbl := data[tblName]
        if _, ok := tsTbl[inParams.key]; !ok {
                log.Info("DbToYang_test_set_key_xfmr instance not found : ", inParams.key)
                return rmap, errors.New("Table Instance not found : " + inParams.key)
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
		rmap["type"] = ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE.ΛMap(typeVal)["ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE"][int64(typeVal)].Name
        }	
        rmap["name"] = testSetName

        log.Info("DbToYang_testsensor_type_key_xfmr rmap ", rmap)
        return rmap, err

}

var DbToYang_test_sensor_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
        var err error
        result := make(map[string]interface{})
        log.Info("DbToYang_test_sensor_group_id_field_xfmr - inParams.uri ", inParams.uri)

        if len(inParams.key) > 0 {
                result["id"] = inParams.key 
       }
        log.Info("DbToYang_test_sensor_group_id_field_xfmr returns ", result)

        return result, err
}

var DbToYang_test_sensor_type_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
        var err error
        result := make(map[string]interface{})

        log.Info("DbToYang_test_sensor_type_field_xfmr - inParams.uri ", inParams.uri)

        pathInfo := NewPathInfo(inParams.uri)
        groupId := pathInfo.Var("id")
        sensorType := pathInfo.Var("type")
        if groupId == "" || sensorType == "" {
                return result, err
        }
       if strings.HasPrefix(sensorType, "sensor") {
               result["type"] = sensorType
       } else {
		errStr := "Invalid Key in uri."
		return result, tlerr.InvalidArgsError{Format: errStr}
       }

        log.Info("DbToYang_test_sensor_type_field_xfmr returns ", result)

        return result, err
}

var DbToYang_test_set_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
       var err error
        result := make(map[string]interface{})

        log.Info("DbToYang_test_set_name_field_xfmr - inParams.uri ", inParams.uri)

        pathInfo := NewPathInfo(inParams.uri)
        setName := pathInfo.Var("name")
        setType := pathInfo.Var("type")
        if setName == "" || setType == "" {
                return result, err
        }
        result["type"] = setName

        log.Info("DbToYang_test_set_name_field_xfmrreturns ", result)

        return result, err
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
                                        if inParams.oper == DELETE {
                                                testSetTableMapNew[testSetName] = db.Value{Field: make(map[string]string)}
                                        } else {
                                                testSetType := findInMap(TEST_SET_TYPE_MAP, strconv.FormatInt(int64(inTestSetKey.Type), 10))
                                                testSetTableMapNew[testSetName] = db.Value{Field: map[string]string{"type":testSetType}}
                                        }
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

var DbToYang_test_port_bindings_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
        var err error
        var testSetTs *db.TableSpec
        var testSetTblMap map[string]db.Value
        var trustIntf, trustTestSet bool

        log.Info("DbToYang_test_port_bindings_xfmr")

        pathInfo := NewPathInfo(inParams.uri)

        testXfmrObj := getTestSetRoot(inParams.ygRoot)
        cdb := inParams.dbs[db.ConfigDB]
        testSetTs = &db.TableSpec{Name: TEST_SET_TABLE}
        testSetKeys, err := cdb.GetKeys(testSetTs)
        if err != nil {
            return err
        }
        testSetTblMap = make(map[string]db.Value)

        for key := range testSetKeys {
                testSetEntry, err := cdb.GetEntry(testSetTs, testSetKeys[key])
                if err != nil {
                        return err
                }
                testSetTblMap[(testSetKeys[key]).Get(0)] = testSetEntry
        }
	targetUriPath, _ := getYangPathFromUri(pathInfo.Path)
        interfaces := make(map[string]bool)
        if isSubtreeRequest(targetUriPath, "/openconfig-test-xfmr:test-xfmr/interfaces") || isSubtreeRequest(targetUriPath, "/openconfig-test-xfmr:test-xfmr/interfaces/interface") {
                intfSbt := testXfmrObj.Interfaces
                if nil == intfSbt.Interface || (nil != intfSbt.Interface && len(intfSbt.Interface) == 0) {
                        log.Info("Get request for all interfaces")
                        for testSetKey := range testSetTblMap {
                                testSetData := testSetTblMap[testSetKey]
                                if len(testSetData.GetList(TEST_SET_PORTS)) > 0 {
                                        testSetIntfs := testSetData.GetList(TEST_SET_PORTS)
                                        for intf := range testSetIntfs {
                                                interfaces[testSetIntfs[intf]] = true
                                        }
                                }
                        }

                        // No interface bindings present. Return. This is general Query ie No interface specified
                        // by the user. We should return no error
                        if len(interfaces) == 0 {
                                return nil
                        }
                        trustIntf = true
                        // For each binding present, create Ygot tree to process next level.
                        ygot.BuildEmptyTree(intfSbt)
                        for intfId := range interfaces {
                                ptr, _ := intfSbt.NewInterface(intfId)
                                ygot.BuildEmptyTree(ptr)
                        }
                } else {
                        log.Info("Get request for specific interface")
                }

		// For each interface present, Process it. The interface present could be created as part of
                // of the URI or created above
                for ifName, ocIntfPtr := range intfSbt.Interface {
                        log.Infof("Processing get request for %s", *ocIntfPtr.Id)
                        if !trustIntf {
                                if targetUriPath == "/openconfig-test-xfmr:test-xfmr/interfaces/interface" && strings.HasSuffix(inParams.requestUri, "]") {
                                        ygot.BuildEmptyTree(ocIntfPtr)
                                }
                        }

                        if nil != ocIntfPtr.Config {
                                ocIntfPtr.Config.Id = ocIntfPtr.Id
                        }
                        if nil != ocIntfPtr.State {
                                ocIntfPtr.State.Id = ocIntfPtr.Id
                        }

                        intfValPtr := reflect.ValueOf(ocIntfPtr)
                        intfValElem := intfValPtr.Elem()

                        testSets := intfValElem.FieldByName("IngressTestSets")
                        if !testSets.IsNil() {
                                testSet := testSets.Elem().FieldByName("IngressTestSet")
                                if testSet.IsNil() || (!testSet.IsNil() && testSet.Len() == 0) {
                                        log.Infof("Get all Ingress Test Sets for %s", ifName)

                                        // Check if any Test Set is applied
                                        for testSet, testSetData := range testSetTblMap {
                                                trustTestSet = true
                                                ports := testSetData.GetList(TEST_SET_PORTS)
						if contains(ports, ifName) {
                                                        testSetType := findInMap(TEST_SET_TYPE_MAP, testSetData.Get(TEST_SET_TYPE))
                                                        n, _ := strconv.ParseInt(testSetType, 10, 64)
                                                        testSetTypeDbkeyComp := "TEST_SET_IPV4"
                                                        if n == int64(ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV6) {
                                                             testSetTypeDbkeyComp = "TEST_SET_IPV6"
                                                        }
                                                        testSetName := getTestSetNameCompFromDbKey(testSet, testSetTypeDbkeyComp)
                                                        log.Infof("Port:%v TestSetName:%v TestSetype:%v ", ifName, testSetName, testSetTypeDbkeyComp)
                                                        testSetOrigType := convertSonicTestSetTypeToOC(testSetData.Get(TEST_SET_TYPE))
                                                        testSet := testSets.MethodByName("NewIngressTestSet").Call([]reflect.Value{reflect.ValueOf(testSetName), reflect.ValueOf(testSetOrigType)})
                                                        ygot.BuildEmptyTree(testSet[0].Interface().(ygot.ValidatedGoStruct))
                                                }
                                        }

                                } else {
                                        log.Info("Get for specific Test Set")
                                }

                                testSetMap := testSets.Elem().FieldByName("IngressTestSet")
                                testSetMapIter := testSetMap.MapRange()
                                for testSetMapIter.Next() {
                                        testSetKey := testSetMapIter.Key()
                                        testSetPtr := testSetMapIter.Value()

                                        if !trustTestSet {
                                                if targetUriPath == "/openconfig-test-xfmr:test-xfmr/interfaces/interface/ingress-test-sets/ingress-test-set" && strings.HasSuffix(inParams.requestUri, "]") {
                                                        ygot.BuildEmptyTree(testSetPtr.Interface().(ygot.ValidatedGoStruct))
                                                }
                                        }

                                        testSetName := testSetKey.FieldByName("SetName")
                                        testSetType := testSetKey.FieldByName("Type")
                                        testSetKeyStr := getTestSetKeyStrFromOCKey(testSetName.String(), testSetType.Interface().(ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE))
                                        testSetData, found := testSetTblMap[testSetKeyStr]
					if found && contains(testSetData.GetList(TEST_SET_PORTS), ifName) {
                                                testSetCfg := testSetPtr.Elem().FieldByName("Config")
                                                if !testSetCfg.IsNil() {
                                                        testSetCfg.Elem().FieldByName("SetName").Set(testSetPtr.Elem().FieldByName("SetName"))
                                                        testSetCfg.Elem().FieldByName("Type").Set(testSetPtr.Elem().FieldByName("Type"))
                                                }
                                                testSetState := testSetPtr.Elem().FieldByName("State")
                                                if !testSetState.IsNil() {
                                                        testSetState.Elem().FieldByName("SetName").Set(testSetPtr.Elem().FieldByName("SetName"))
                                                        testSetState.Elem().FieldByName("Type").Set(testSetPtr.Elem().FieldByName("Type"))
                                                }

                                        }
                                 }
                        }
                }

        }
        return err
}

var Subscribe_test_port_bindings_xfmr SubTreeXfmrSubscribe = func(inParams XfmrSubscInParams) (XfmrSubscOutParams, error) {
        var err error
        var result XfmrSubscOutParams

        pathInfo := NewPathInfo(inParams.uri)
        targetUriPath, _ := getYangPathFromUri(pathInfo.Path)
        print("Subscribe_test_port_bindings_xfmr targetUriPath:", targetUriPath)
        result.isVirtualTbl = true
        log.Info("Returning Subscribe_test_port_bindings_xfmr")
        return result, err
}

func convertSonicTestSetTypeToOC(testSetType string) ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE {
        var testSetOrigType ocbinds.E_OpenconfigTestXfmr_TEST_SET_TYPE

        if "IPV4" == testSetType {
                testSetOrigType = ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV4
        } else if "IPV6" == testSetType {
                testSetOrigType = ocbinds.OpenconfigTestXfmr_TEST_SET_TYPE_TEST_SET_IPV6
        } else {
                log.Infof("Unknown type %v", testSetType)
        }

        return testSetOrigType
}
