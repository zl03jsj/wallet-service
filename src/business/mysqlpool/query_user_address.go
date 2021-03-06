package mysqlpool

import (
	. "business/def"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"time"
)

func QueryUserAddress(queryMap map[string]interface{}) ([]UserAddress, bool) {
	sqls := "select user_key,user_class,asset_name,address,private_key,available_amount,frozen_amount,enabled," +
		" unix_timestamp(create_time), unix_timestamp(allocation_time), unix_timestamp(update_time) from user_address" +
		" where true"

	userAddress := make([]UserAddress, 0)
	params := make([]interface{}, 0)

	if len(queryMap) > 0 {
		sqls += andConditions(queryMap, &params)
		sqls += andPagination(queryMap, &params)
	}

	db := Get()
	rows, err := db.Query(sqls, params...)
	if err != nil {
		fmt.Println(err.Error())
		return userAddress, len(userAddress) > 0
	}

	var data UserAddress
	for rows.Next() {
		err := rows.Scan(&data.UserKey, &data.UserClass, &data.AssetName, &data.Address, &data.PrivateKey,
			&data.AvailableAmount, &data.FrozenAmount, &data.Enabled, &data.CreateTime, &data.AllocationTime,
			&data.UpdateTime)
		if err == nil {
			userAddress = append(userAddress, data)
		}
	}
	return userAddress, len(userAddress) > 0
}

func QueryUserAddressCount(queryMap map[string]interface{}) int {
	sqls := "select count(*) from user_address" +
		" where true"

	count := 0
	params := make([]interface{}, 0)

	if len(queryMap) > 0 {
		sqls += andConditions(queryMap, &params)
	}

	db := Get()
	db.QueryRow(sqls, params...).Scan(&count)
	return count
}

func QueryUserAddressByNameAddress(assetName string, address string) (UserAddress, bool) {
	queryMap := make(map[string]interface{})
	queryMap["asset_name"] = assetName
	queryMap["address"] = address
	if userAddress, ok := QueryUserAddress(queryMap); ok {
		return userAddress[0], true
	}
	return UserAddress{}, false
}

func QueryPayAddress(assetName string) (UserAddress, bool) {
	db := Get()
	row := db.QueryRow("select user_key,user_class,asset_name,address,private_key,"+
		" available_amount,frozen_amount,enabled, unix_timestamp(create_time), unix_timestamp(update_time)"+
		" from pay_address_view where asset_name = ?", assetName)

	data := UserAddress{}
	err := row.Scan(&data.UserKey, &data.UserClass, &data.AssetName, &data.Address, &data.PrivateKey, &data.AvailableAmount,
		&data.FrozenAmount, &data.Enabled,
		&data.CreateTime, &data.UpdateTime)
	if err != nil {
		return UserAddress{}, false
	}
	return data, true
}

func QueryPayAddressList(assetName []string) ([]UserAddress, bool) {
	sqls := "select user_key,user_class,asset_name,address,private_key,available_amount,frozen_amount,enabled, " +
		" unix_timestamp(create_time), unix_timestamp(update_time) from pay_address_view where true "

	params := make([]interface{}, 0)
	for _, v := range assetName {
		sqls += " and asset_name = ?"
		params = append(params, v)
	}

	arr := make([]UserAddress, 0)
	db := Get()
	rows, err := db.Query(sqls, params...)
	if err != nil {
		return arr, false
	}

	for rows.Next() {
		data := UserAddress{}
		err := rows.Scan(&data.UserKey, &data.UserClass, &data.AssetName, &data.Address, &data.PrivateKey, &data.AvailableAmount,
			&data.FrozenAmount, &data.Enabled,
			&data.CreateTime, &data.UpdateTime)
		if err == nil {
			arr = append(arr, data)
		}
	}
	return arr, len(arr) > 0
}

func SetPayAddress(assetName string, address string) error {
	db := Get()
	row := db.QueryRow("select a.private_key from user_address a"+
		" left join asset_property b on a.asset_name = b.asset_name"+
		" where a.user_class = 1 and b.is_token = 0 and a.asset_name = ? and a.address = ?",
		assetName, address)

	var privateKey string
	err := row.Scan(&privateKey)
	if err != nil {
		return err
	}

	_, err = db.Exec("insert pay_address (asset_name, address, private_key) values (?, ?, ?);",
		assetName, address, privateKey)
	if err != nil {
		_, err := db.Exec("update pay_address set address = ?, private_key = ? where asset_name = ?;",
			address, privateKey, assetName)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddUserAddress(userAddress []UserAddress) error {
	tx, err := Get().Begin()
	if err != nil {
		_, errMsg := CheckError(ErrorFailed, err.Error())
		l4g.Error(errMsg)
		return err
	}

	for _, v := range userAddress {
		_, err := tx.Exec("insert user_address (user_key, user_class, asset_name, address, private_key,"+
			" available_amount, frozen_amount, enabled, create_time, allocation_time, update_time)"+
			" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			v.UserKey, v.UserClass, v.AssetName, v.Address, v.PrivateKey, v.AvailableAmount, v.FrozenAmount, v.Enabled,
			time.Unix(v.CreateTime, 0).UTC().Format(TimeFormat),
			time.Unix(v.UpdateTime, 0).UTC().Format(TimeFormat),
			time.Unix(v.UpdateTime, 0).UTC().Format(TimeFormat))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_, errMsg := CheckError(ErrorFailed, err.Error())
		l4g.Error(errMsg)
		return err
	}

	return nil
}

func UpdatePayTokenAddress() error {
	db := Get()
	rows, err := db.Query("select a.user_key, a.user_class, a.available_amount, a.frozen_amount," +
		" a.asset_name, address, a.private_key, a.enabled, unix_timestamp(a.create_time)," +
		" unix_timestamp(a.allocation_time), unix_timestamp(a.update_time)" +
		" from user_address a left join asset_property b on a.asset_name = b.asset_name" +
		" where a.user_class = 1 and b.is_token = 0")
	defer rows.Close()
	if err != nil {
		return err
	}

	userAddress := make([]UserAddress, 0)
	var data UserAddress
	for rows.Next() {
		err := rows.Scan(&data.UserKey, &data.UserClass, &data.AvailableAmount, &data.FrozenAmount,
			&data.AssetName, &data.Address, &data.PrivateKey, &data.Enabled,
			&data.CreateTime, &data.AllocationTime, &data.UpdateTime)
		if err != nil {
			continue
		}
		userAddress = append(userAddress, data)
	}

	if len(userAddress) > 0 {
		CreateTokenAddress(userAddress)
		return nil
	}
	return nil
}

func CreateTokenAddress(userAddress []UserAddress) {
	assetPropertyMap := make(map[string]AssetProperty)
	if assetProperty, ok := QueryAssetProperty(nil); ok {
		for _, value := range assetProperty {
			assetPropertyMap[value.AssetName] = value
		}
	}

	var assetName string
	data := make([]UserAddress, 0)

	db := Get()
	for _, value := range userAddress {
		if value.UserClass == 1 {
			if v, ok := assetPropertyMap[value.AssetName]; ok {
				if v.IsToken == 0 {
					//设置默认支付地址
					db.Exec("insert pay_address (asset_name, address, private_key) values (?, ?, ?)",
						value.AssetName, value.Address, value.PrivateKey)

					rows, err := db.Query("select asset_name"+
						" from asset_property where is_token = 1 and parent_name = ?", v.AssetName)
					if err != nil {
						continue
					} else {
						for rows.Next() {
							err := rows.Scan(&assetName)
							if err != nil {
								continue
							}
							addr := value
							addr.AssetName = assetName
							addr.AvailableAmount = 0
							addr.FrozenAmount = 0
							data = append(data, addr)
						}
					}
					rows.Close()
				}
			}
		}
	}
	if len(data) > 0 {
		db := Get()
		for _, v := range data {
			db.Exec("insert user_address (user_key, user_class, asset_name, address, private_key,"+
				" available_amount, frozen_amount, enabled, create_time, allocation_time, update_time)"+
				" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
				v.UserKey, v.UserClass, v.AssetName, v.Address, v.PrivateKey, v.AvailableAmount, v.FrozenAmount, v.Enabled,
				time.Unix(v.CreateTime, 0).UTC().Format(TimeFormat),
				time.Unix(v.UpdateTime, 0).UTC().Format(TimeFormat),
				time.Unix(v.UpdateTime, 0).UTC().Format(TimeFormat))
		}
	}
}
