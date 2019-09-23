package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IClusterCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceCondition;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackDetailRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.service.ITrackService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

/**
 * 轨迹控制器
 *
 * @author sctang2
 * @create 2017-11-24 11:56
 **/
@RestController
@RequestMapping(Constant.API + "/{v}/track")
public class TrackController {
	private final EasyLogger logger = EasyLoggerFactory.getInstance(TrackController.class);
	@Autowired
	private ITrackService trackServiceImpl;

	@Autowired
	private IProjectCondition projectConditionImpl;

	@Autowired
	private IClusterCondition clusterConditionImpl;

	@Autowired
	private IServiceCondition serviceConditionImpl;

	/**
	 * 查询最近的配置推送轨迹列表
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/config/lastestList", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> findConfigList(QueryTrackRequestBody body) {
		return this.trackServiceImpl.findLastestConfigList(body);
	}

	/**
	 * 查询最近的服务发现推送轨迹列表
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/discovery/lastestList", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> findDiscoveryList(QueryTrackRequestBody body) {
		return this.trackServiceImpl.findLastestDiscoveryList(body);
	}

	/**
	 * 查询配置推送轨迹明细
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/config/detail", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> findConfig(@Validated QueryTrackDetailRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.trackServiceImpl.findConfig(body);
	}

	/**
	 * 查询服务发现推送轨迹明细
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/discovery/detail", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> findDiscovery(@Validated QueryTrackDetailRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.trackServiceImpl.findDiscovery(body);
	}

	/**
	 * 删除配置推送轨迹
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/config/delete", method = RequestMethod.POST)
	public Response<String> deleteConfig(@Validated @RequestBody IdRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.trackServiceImpl.deleteConfig(body);
	}

	/**
	 * 批量删除配置推送轨迹
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/config/batchDelete", method = RequestMethod.POST)
	public Response<String> batchDeleteConfig(@Validated @RequestBody IdsRequestBody body) {
		List<String> ids = body.getIds();
		if (null == ids || ids.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_TRACK_IDS_NOT_NULL);
		}
		return this.trackServiceImpl.batchDeleteConfig(body);
	}

	/**
	 * 删除服务发现推送轨迹
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/discovery/delete", method = RequestMethod.POST)
	public Response<String> deleteDiscovery(@Validated @RequestBody IdRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.trackServiceImpl.deleteDiscovery(body);
	}

	/**
	 * 批量删除服务发现推送轨迹
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/discovery/batchDelete", method = RequestMethod.POST)
	public Response<String> batchDeleteDiscovery(@Validated @RequestBody IdsRequestBody body) {
		List<String> ids = body.getIds();
		if (null == ids || ids.isEmpty()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_TRACK_IDS_NOT_NULL);
		}
		return this.trackServiceImpl.batchDeleteDiscovery(body);
	}

	/**
	 * 快速查询服务推送轨迹下拉框
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/discovery/list", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> list(BaseRequestBody body) {
		return this.trackServiceImpl.findList(body);
	}
}
