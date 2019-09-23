package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.domain.DownloadFile;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.QueryCustomDetailRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.PushServiceConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IGrayService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceConfig;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import javax.servlet.http.HttpServletResponse;
import java.util.List;

/**
 * 服务配置控制器
 *
 * @author sctang2
 * @create 2017-11-21 11:45
 **/
@Api(value = "配置中心模块", description = "配置中心模块")
@RestController
@RequestMapping(Constant.API + "/{v}/service/config")
public class ServiceConfigController {
	@Autowired
	private IServiceConfig serviceConfigImpl;

	@Autowired
	private IGrayService grayServiceImpl;

	/**
	 * 查询最近的配置列表
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/lastestList", method = RequestMethod.GET)
	@ApiOperation("查询配置列表")
	public Response<QueryPagingListResponseBody> lastestList(@Validated QueryServiceConfigRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.findLastestList(body);
	}

	/**
	 * 编辑配置
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/edit", method = RequestMethod.POST)
	@ApiOperation(value = "编辑服务配置", notes = "编辑服务配置")
	public Response<ServiceConfigDetailResponseBody> edit(@Validated @RequestBody EditServiceConfigRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.edit(body);
	}

	/**
	 * 查询配置详情
	 *
	 * @param id
	 * @return
	 */
	@RequestMapping(value = "/detail", method = RequestMethod.GET)
	@ApiOperation(value = "查询配置详情", notes = "查询配置详情")
	public Response<ServiceConfigDetailResponseBody> find(@ApiParam("配置id") @RequestParam("id") String id) {
		if (StringUtils.isEmpty(id)) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
		}
		return this.serviceConfigImpl.find(id);
	}

	/**
	 * 推送配置
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/push", method = RequestMethod.POST)
	@ApiOperation("推送配置")
	public Response<PushServiceConfigResponseBody> push(@Validated @RequestBody PushServiceConfigRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.push(body);
	}

	/**
	 * 批量推送配置
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/batchPush", method = RequestMethod.POST)
	@ApiOperation("批量推送")
	public Response<PushServiceConfigResponseBody> batchPush(@RequestBody BatchPushServiceConfigRequestBody body) {
		List<String> configIds = body.getIds();
		if (null == configIds || configIds.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL);
		}
		List<String> ids = body.getRegionIds();
		if (null == ids || ids.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_REGION_IDS_NOT_NULL);
		}
		return this.serviceConfigImpl.batchPush(body);
	}

	/**
	 * 删除配置
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/delete", method = RequestMethod.POST)
	public Response<String> delete(@Validated @RequestBody IdRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.delete(body);
	}

	/**
	 * 批量删除配置
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/batchDelete", method = RequestMethod.POST)
	public Response<String> batchDelete(@RequestBody IdsRequestBody body) {
		List<String> ids = body.getIds();
		if (null == ids || ids.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL);
		}
		return this.serviceConfigImpl.batchDelete(body);
	}

	/**
	 * 查询配置历史列表
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/historyList", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> historyList(@Validated ServiceConfigHistoryListRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.findServiceConfigHistoryList(body);
	}

	/**
	 * 回滚配置
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/rollback", method = RequestMethod.POST)
	public Response<List<ServiceConfigDetailResponseBody>> rollback(@RequestBody IdsRequestBody body) {
		List<String> ids = body.getIds();
		if (null == ids || ids.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL);
		}
		return this.serviceConfigImpl.rollback(body);
	}

	/**
	 * 更新反馈
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/feedback", method = RequestMethod.POST)
	public Response<String> feedback(@Validated @RequestBody ServiceConfigFeedBackRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.serviceConfigImpl.feedback(body);
	}

	/**
	 * 服务订阅者查看
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/consumer", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> consumer(@Validated QueryCustomDetailRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.grayServiceImpl.consumer(body);
	}

	@RequestMapping(value = "/download", method = RequestMethod.GET)
	public Response<DownloadFile> download(@Validated DownloadServiceConfigRequestBody body, HttpServletResponse response, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_DOWNLOAD_FALSE, SystemErrCode.ERRCODE_SERVICE_CONFIG_DOWNLOAD_FALSE_MESSAGE);
		}
		return this.serviceConfigImpl.download(body);
	}

}
