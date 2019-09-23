package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.CopyClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.EditClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.QueryClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IClusterService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 集群控制器
 *
 * @author sctang2
 * @create 2017-11-15 17:38
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/cluster")
public class ClusterController {
	@Autowired
	private IClusterService clusterServiceImpl;

	/**
	 * 查询最近的集群列表
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/lastestList", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> lastestList(QueryClusterRequestBody body) {
		return this.clusterServiceImpl.findLastestList(body);
	}

	/**
	 * 查询集群 列表
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/list", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> list(QueryClusterRequestBody body) {
		return this.clusterServiceImpl.findList(body);
	}

	/**
	 * 查询集群详情
	 *
	 * @param id
	 * @return
	 */
	@RequestMapping(value = "/detail", method = RequestMethod.GET)
	public Response<ClusterDetailResponseBody> find(@RequestParam("id") String id) {
		if (StringUtils.isEmpty(id)) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
		}
		return this.clusterServiceImpl.find(id);
	}

	/**
	 * 新增集群
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/add", method = RequestMethod.POST)
	public Response<ClusterDetailResponseBody> add(@Validated @RequestBody AddClusterRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.clusterServiceImpl.add(body);
	}

	/**
	 * 编辑 集群
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/edit", method = RequestMethod.POST)
	public Response<ClusterDetailResponseBody> edit(@Validated @RequestBody EditClusterRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.clusterServiceImpl.edit(body);
	}

	/**
	 * 删除集群
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
		return this.clusterServiceImpl.delete(body);
	}

	/**
	 * 复制集群
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/copy", method = RequestMethod.POST)
	public Response<ClusterDetailResponseBody> copy(@Validated @RequestBody CopyClusterRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.clusterServiceImpl.copy(body);
	}
}
