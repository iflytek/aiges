package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.AddApiVersionResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceDiscoveryResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IMyServiceDiscovery;
import com.iflytek.ccr.polaris.cynosure.service.IServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.service.IServiceDiscovery;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import java.util.ArrayList;
import java.util.List;

/**
 * 服务发现控制器
 *
 * @author sctang2
 * @create 2017-12-05 16:38
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/service/discovery")
public class ServiceDiscoveryController {
    @Autowired
    private IMyServiceDiscovery myServiceDiscovery;
    @Autowired
    private IServiceDiscovery  serviceDiscoveryImpl;
    @Autowired
    private IServiceApiVersion serviceApiVersionImpl;

    /**
     * 新增版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public Response<AddApiVersionResponseBody> add(@Validated @RequestBody AddServiceDiscoveryRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.add(body);
    }

    /**
     * 查询最近的服务发现列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/lastestList", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> lastestList(@Validated QueryServiceDiscoveryListRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.findLastestList1(body);
    }

    /**
     * 查询服务发现明细
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ServiceDiscoveryResponseBody> detail(@Validated QueryServiceDiscoveryDetailRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.find(body);
    }

//    /**
//     * 编辑服务发现
//     *
//     * @param body
//     * @param result
//     * @return
//     */
//    @RequestMapping(value = "/edit", method = RequestMethod.POST)
//    public Response<String> edit(@Validated @RequestBody EditServiceDiscoveryRequestBody body, BindingResult result) {
//        if (result.hasErrors()) {
//            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
//        }
//
//        if (null != body.getParams() && !body.getParams().isEmpty()) {
//            //校验自定义规则的key值重复
//            List<String> keyList = new ArrayList<>();
//            List<ServiceParam> params = body.getParams();
//            for (ServiceParam serviceParam : params) {
//                String key = serviceParam.getKey();
//                keyList.add(key);
//            }
//            long count = keyList.stream().distinct().count();
//            boolean isRepeat = count < keyList.size();
//            if (isRepeat) {
//                return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_PARAMS_REPEAT);
//            }
//        }
//        return this.serviceDiscoveryImpl.edit(body);
//    }

    /**
     * 编辑服务发现
     *
     * @param body
     * @param  result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<String> edit1(@Validated @RequestBody EditServiceDiscoveryRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }

        if (null != body.getParams() && !body.getParams().isEmpty()) {
            //校验自定义规则的key值重复
            List<String> keyList = new ArrayList<>();
            List<ServiceParam> params = body.getParams();
            for (ServiceParam serviceParam : params) {
                String key = serviceParam.getKey();
                keyList.add(key);
            }
            long count = keyList.stream().distinct().count();
            boolean isRepeat = count < keyList.size();
            if (isRepeat) {
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_PARAMS_REPEAT);
            }
        }
        return this.myServiceDiscovery.edit1(body);
    }

    /**
     * 查询服务提供者列表
     *
     * @param  body
     * @param result
     * @return
     */
    @RequestMapping(value = "/provider", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> provider(@Validated QueryServiceDiscoveryDetailRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.provider(body);
    }

    /**
     * 编辑服务提供者
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/provider/edit", method = RequestMethod.POST)
    public Response<String> editProvider(@Validated @RequestBody EditServiceDiscoveryProviderRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.editProvider(body);
    }


    /**
     * 查询服务消费方列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/consumer", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> consumer(@Validated QueryServiceDiscoveryDetailRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.consumer(body);
    }

    /**
     * 更新反馈
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/feedback", method = RequestMethod.POST)
    public Response<String> feedback(@Validated @RequestBody ServiceDiscoveryFeedBackRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceDiscoveryImpl.feedback(body);
    }

    /**
     * 查询服务发现负载均衡规则
     *
     * @return
     */
    @RequestMapping(value = "/loadBalanceList", method = RequestMethod.GET)
    public Response<List<LoadBalance>> loadBalanceList() {
        return this.serviceDiscoveryImpl.findBalanceList();
    }

    /**
     * 新增版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/addApiVersion", method = RequestMethod.POST)
    public Response<ServiceApiVersionDetailResponseBody> add(@Validated @RequestBody AddServiceApiVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceApiVersionImpl.add(body);
    }
}
