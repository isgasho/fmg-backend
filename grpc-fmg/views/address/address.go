package address

import (
	_ "context"
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	accountException "grpc-demo/exceptions/account"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"

	//"encoding/json"
	//"errors"
	//"fmt"
	_ "grpc-demo/exceptions/system"
	_ "grpc-demo/models/db"

	_ "net/http"
	//"strings"
)

type Result struct {
	CityName     string `json:"city_name"`
	ProvinceName string `json:"province_name"`
	DistrictName  string `json:"district_name"`
}

func MGetAddress(ctx iris.Context, auth authbase.AuthAuthorization, uid int) {

	var addresses []db.Address
	//var results []Result
	db.Driver.Where("account_id = ?", uid).Find(&addresses)


	data := make([]interface{}, 0, len(addresses))
	for _, address := range addresses {
		func(data *[]interface{}) {
			aid := address.ID
			if address.IsDelete == false {
				names := GetAddressStr(aid)
				info := paramsUtils.ModelToDict(address, []string{"ID", "ProvinceID", "CountryId", "CityID",
					"DistrictID", "Detail", "Name", "Phone"})
				info["city_name"] = names.CityName
				info["province_name"] = names.ProvinceName
				info["district_name"] = names.DistrictName
				*data = append(*data, info)
			}
			//db.Driver.Select("city","province","district").Where("id=?,id=?,id=?",city,province,district).Scan(&names)

			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)

}


func GetAddressStr( aid int) *Result {
	var address db.Address
	//db.Driver.Where("id = ?", aid).First(&address)

		err := db.Driver.GetOne("address", aid, &address)
		if err != nil {
			fmt.Print(err)
			panic(accountException.AddressNotFount())
		}
		results :=Result{}
	db.Driver.Table("city, province, district").Debug().Select("city.name as city_name, province.name as province_name ,district.name as district_name").Where("province.id = ? and city.id = ? and district.id = ?", address.ProvinceID, address.CityID, address.DistrictID).Find(&results)
	fmt.Println(results)

	//if len(results) < 1 {
	//	ctx.JSON(iris.Map{})
	//} else {
	//	ctx.JSON(results[0])
	//}
	return &results

}
func GetAddress(ctx iris.Context, auth authbase.AuthAuthorization, aid int) {
	var address db.Address
	db.Driver.Where("id = ?", aid).First(&address)
	names := GetAddressStr(aid)
	info := paramsUtils.ModelToDict(address, []string{"ID", "ProvinceID", "CountryId", "CityID",
		"DistrictID", "Detail", "Name", "Phone"})

	info["city_name"] = names.CityName
	info["province_name"] = names.ProvinceName
	info["district_name"] = names.DistrictName

	ctx.JSON(info)

}
func ListAddress(ctx iris.Context, auth authbase.AuthAuthorization) {

	if country_id := ctx.URLParamIntDefault("country_id", 0); country_id != 0 {
		var provinces []db.Province
		db.Driver.Where("country_id = ?", 1).Find(&provinces)
		names := make([]interface{}, 0, len(provinces))
		for _, v := range provinces {
			func(names *[]interface{}) {
				*names = append(*names, paramsUtils.ModelToDict(v, []string{"ID", "Name"}))
				defer func() {
					recover()
				}()
			}(&names)
		}
		ctx.JSON(names)
	}

	if province_id := ctx.URLParamIntDefault("province_id", 0); province_id != 0 {
		var citys []db.City

		db.Driver.Where("province_id = ?", province_id).Find(&citys)

		names := make([]interface{}, 0, len(citys))
		for _, v := range citys {
			names = append(names, paramsUtils.ModelToDict(v, []string{"ID", "Name"}))
		}
		ctx.JSON(names)
	}

	if city_id := ctx.URLParamIntDefault("city_id", 0); city_id != 0 {
		var districts []db.District
		db.Driver.Where("city_id = ?", city_id).Find(&districts)
		names := make([]interface{}, 0, len(districts))
		for _, v := range districts {
			names = append(names, paramsUtils.ModelToDict(v, []string{"ID", "Name"}))
		}
		ctx.JSON(names)
	}

}

func CreatAddress(ctx iris.Context, auth authbase.AuthAuthorization, uid int) {

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	province := params.Int("province_id", "province_id")
	country := params.Int("country_id", "country_id")
	city := params.Int("city_id", "city_id")
	district := params.Int("district_id", "district_id")
	detail := params.Str("detail", "detail")
	name := params.Str("name", "地址使用用户名")
	phone := params.Str("phone", "电话")

	var address db.Address
	address = db.Address{
		CountryId:  country,
		ProvinceID: province,
		CityID:     city,
		DistrictID: district,
		Detail:     detail,
		AccountId:  uid,
		Name:       name,
		Phone:      phone,
	}

	db.Driver.Create(&address)
	ctx.JSON(iris.Map{
		"id": address.ID,
	})

}

func PutAddress(ctx iris.Context, auth authbase.AuthAuthorization, aid int) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	var address db.Address
	params.Diff(address)
	if err := db.Driver.GetOne("address", aid, &address); err != nil {
		panic(accountException.AddressNotFount())
	}
	if params.Has("city_id") {
		newCity := params.Int("city_id", "新城市")
		address.CityID = newCity
	}
	if params.Has("name") {
		name := params.Str("name", "新名字")
		address.Name = name
	}
	if params.Has("district_id") {
		district_id := params.Int("district_id", "新地区")
		address.DistrictID = district_id
	}
	if params.Has("phone") {
		province := params.Int("province_id", "province_id")
		address.ProvinceID = province
	}
	if params.Has("detail") {
		detail := params.Str("detail", "新地址")
		address.Detail = detail
	}
	if params.Has("phone") {
		phone := params.Str("phone", "新手机")
		address.Phone = phone
	}
	if params.Has("is_deleted") {
		isDelete := params.Bool("is_deleted", "是否删除")
		address.IsDelete = isDelete
	}
	db.Driver.Save(&address)

	ctx.JSON(iris.Map{
		"id": address.ID,
	})

}

func DeleteAddress(ctx iris.Context, auth authbase.AuthAuthorization, aid int) {
	var address db.Address
	if err := db.Driver.GetOne("address", aid, &address); err == nil {

		db.Driver.Delete(address)
	}else{
		panic(accountException.AddressNotFount())
	}

	ctx.JSON(iris.Map{
		"id": aid,
	})
}
