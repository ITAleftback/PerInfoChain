/**
 * @Author: Anpw
 * @Description:
 * @File:  module_code
 * @Version: 1.0.0
 * @Date: 2021/5/26 22:56
 */

package errcode

var (
	ErrorCreateUserFail = NewError(20010001, "创建用户失败")
	ErrorLoginUserFail  = NewError(20010002, "用户登录失败")
	ErrorUpdateCodeFail = NewError(20010003, "修改密码失败")

	ErrorAddUserInfoFail    = NewError(20020001, "添加用户基本信息失败")
	ErrorDeleteUserInfoFail = NewError(20020002, "删除用户基本信息失败")
	ErrorUpdateUserInfoFail = NewError(20020003, "更改用户基本信息失败")
	ErrorSelectUserInfoFail = NewError(20020004, "查询用户基本信息失败")

	ErrorAddEducationUserInfoFail    = NewError(20030001, "添加教育基本信息失败")
	ErrorDeleteEducationUserInfoFail = NewError(20030002, "删除教育基本信息失败")
	ErrorUpdateEducationUserInfoFail = NewError(20030003, "更改教育基本信息失败")
	ErrorSelectEducationUserInfoFail = NewError(20030004, "查询教育基本信息失败")

	ErrorAddEducationScoreInfoFail    = NewError(20050001, "添加成绩基本信息失败")
	ErrorDeleteEducationScoreInfoFail = NewError(20050002, "删除成绩基本信息失败")
	ErrorUpdateEducationScoreInfoFail = NewError(20050003, "更改成绩基本信息失败")
	ErrorSelectEducationScoreInfoFail = NewError(20050004, "查询成绩基本信息失败")

	ErrorAddCriminalCaseFail    = NewError(20060001, "添加公安基本信息失败")
	ErrorDeleteCriminalCaseFail = NewError(20060002, "删除公安基本信息失败")
	ErrorUpdateCriminalCaseFail = NewError(20060003, "更改公安基本信息失败")
	ErrorSelectCriminalCaseFail = NewError(20060004, "查询公安基本信息失败")

	ErrorAddHouseInfoFail    = NewError(20070001, "添加住房基本信息失败")
	ErrorDeleteHouseInfoFail = NewError(20070002, "删除住房基本信息失败")
	ErrorUpdateHouseInfoFail = NewError(20070003, "更改住房基本信息失败")
	ErrorSelectHouseInfoFail = NewError(20070004, "查询住房基本信息失败")

	ErrorAddBankInfoFail    = NewError(20080001, "添加银行基本信息失败")
	ErrorDeleteBankInfoFail = NewError(20080002, "删除银行基本信息失败")
	ErrorUpdateBankInfoFail = NewError(20080003, "更改银行基本信息失败")
	ErrorSelectBankInfoFail = NewError(20080004, "查询银行基本信息失败")

	ErrorAddBankDetailsFail    = NewError(20090001, "添加银行明细基本信息失败")
	ErrorDeleteBankDetailsFail = NewError(20090002, "删除银行明细基本信息失败")
	ErrorUpdateBankDetailsFail = NewError(20090003, "更改银行明细基本信息失败")
	ErrorSelectBankDetailsFail = NewError(20090004, "查询银行明细基本信息失败")

	ERROR_UPLOAD_FILE_FAIL = NewError(20040001, "上传文件失败")
	ERROR_UPDATE_FILE_FAIL = NewError(20040002, "更新文件失败")
	ERROR_DELETE_FILE_FAIL = NewError(20040003, "删除文件失败")

)
