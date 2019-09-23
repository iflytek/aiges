package com.iflytek.ccr.polaris.cynosure.controller.v2;

import com.iflytek.ccr.polaris.cynosure.apicompatible.ApiVersion;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * 用户控制器V2版本
 *
 * @author sctang2
 * @create 2017-11-09 15:51
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/user")
@ApiVersion(2)
public class UserControllerV2 {
}
