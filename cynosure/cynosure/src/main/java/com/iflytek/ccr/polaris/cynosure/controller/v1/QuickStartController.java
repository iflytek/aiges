package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceConfigResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddVersionResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.service.IQuickStartService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceConfig;

import io.swagger.annotations.ApiOperation;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * 快速开始控制器
 *
 * @author sctang2
 * @create 2018-01-29 12:00
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/quickStart")
public class QuickStartController {
	@Autowired
	private IQuickStartService quickStartServiceImpl;

	@Autowired
	private IServiceConfig serviceConfigImpl;

	/**
	 * 新增服务
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/addService", method = RequestMethod.POST)
	public Response<AddServiceResponseBodyByQuickStart> addService(@Validated @RequestBody AddServiceRequestBodyByQuickStart body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.quickStartServiceImpl.addService(body);
	}

	/**
	 * 新增版本
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/addVersion", method = RequestMethod.POST)
	public Response<AddVersionResponseBodyByQuickStart> addVersion(@Validated @RequestBody AddVersionRequestBodyByQuickStart body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.quickStartServiceImpl.addVersion(body);
	}

	/**
	 * 新增服务版本
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/addServiceVersion", method = RequestMethod.POST)
	public Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersion(MultipartFile[] file, @Validated AddServiceVersionRequestBodyByQuickStart body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}

		// 取出配置Name的列表
		if (null != body.getIds() && !body.getIds().isEmpty()) {
			List<ServiceConfig> configList = this.serviceConfigImpl.findListByIds(body.getIds());
			List<String> serviceConfigNames = new ArrayList<>();
			for (ServiceConfig serviceConfig : configList) {
				serviceConfigNames.add(serviceConfig.getName());
			}

			Set<String> setServiceConfigName = new HashSet<>();
			if (serviceConfigNames.size() > 0) {
				// 检查配置文件名是否重复
				for (String name : serviceConfigNames) {
					if (!setServiceConfigName.add(name)) {
						return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
					}
				}
			}

		}

		// 检查是否存在拖拽文件
		if (null != file && file.length > 0) {
			List<FileContent> fileContentList = new ArrayList<>();
			List<String> fileNames = new ArrayList<>();
			FileContent fileContent;

			// 检查上传的文件中是否有重名的，同时还要排除是不是有和最近配置一样的文件名
			for (MultipartFile item : file) {
				// 检查文件名
				String fileName = item.getOriginalFilename();
				if (fileNames.contains(fileName)) {
					return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
				}

				// 检查上传文件名和最近的配置文件是否有重复
				if (null != body.getIds() && !body.getIds().isEmpty()) {
					List<ServiceConfig> configList = this.serviceConfigImpl.findListByIds(body.getIds());
					List<String> serviceConfigNames = new ArrayList<>();
					for (ServiceConfig serviceConfig : configList) {
						serviceConfigNames.add(serviceConfig.getName());
					}
					if (serviceConfigNames.contains(fileName)) {
						return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
					}
				}
				// 检查文件内容
				if (item.isEmpty()) {
					return new Response<>(SystemErrCode.ERRCODE_FILE_CONTENT_NOT_NULL, SystemErrCode.ERRMSG_FILE_CONTENT_NOT_NULL);
				}

				// 创建文件数据
				fileContent = new FileContent();
				fileContent.setFileName(fileName);
				try {
					byte[] content = item.getBytes();
					fileContent.setContent(content);
				} catch (IOException e) {
					GlobalExceptionUtil.log(e);
				}
				fileNames.add(fileName);
				fileContentList.add(fileContent);
			}

			// 调用含文件上传服务接口
			return this.quickStartServiceImpl.addServiceVersionAndFile(body, fileContentList);
		} else {
			// 调用不含文件上传新增服务接口
			return this.quickStartServiceImpl.addServiceVersion(body);
		}
	}

	/**
	 * 新增服务配置
	 *
	 * @param file
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/addServiceConfig", method = RequestMethod.POST)
	@ApiOperation("新增服务配置")
	public Response<AddServiceConfigResponseBodyByQuickStart> addServiceConfig(MultipartFile[] file, @Validated AddServiceConfigRequestBodyByQuickStart body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}

		// 检查是否存在文件
		int size = file.length;
		if (size <= 0) {
			return new Response<>(SystemErrCode.ERRCODE_NOT_FILE, SystemErrCode.ERRMSG_NOT_FILE);
		}

		List<FileContent> fileContentList = new ArrayList<>();
		List<String> fileNames = new ArrayList<>();
		FileContent fileContent;
		for (MultipartFile item : file) {
			// 检查文件名
			String fileName = item.getOriginalFilename();
			if (fileNames.contains(fileName)) {
				return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
			}

			// 检查文件内容
			if (item.isEmpty()) {
				return new Response<>(SystemErrCode.ERRCODE_FILE_CONTENT_NOT_NULL, SystemErrCode.ERRMSG_FILE_CONTENT_NOT_NULL);
			}

			// 创建文件数据
			fileContent = new FileContent();
			fileContent.setFileName(fileName);
			try {
				byte[] content = item.getBytes();
				fileContent.setContent(content);
			} catch (IOException e) {
				GlobalExceptionUtil.log(e);
			}
			fileNames.add(fileName);
			fileContentList.add(fileContent);
		}
		Response<AddServiceConfigResponseBodyByQuickStart> result1 = this.quickStartServiceImpl.addServiceConfig(body, fileContentList);
		return result1;
	}

	/**
	 * 快速查询
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/list", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> list(BaseRequestBody body) {
		return this.quickStartServiceImpl.findList(body);
	}

	/**
	 * 快速查询api版本
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/list1", method = RequestMethod.GET)
	public Response<QueryPagingListResponseBody> list1(BaseRequestBody body) {
		return this.quickStartServiceImpl.findList1(body);
	}

}
