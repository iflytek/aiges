package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IRegionCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.AddRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.EditRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.region.RegionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IRegionService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.*;

/**
 * 区域业务接口实现
 *
 * @author sctang2
 * @create 2017-11-14 20:53
 **/
@Service
public class RegionServiceImpl extends BaseService implements IRegionService {
    @Autowired
    private IRegionCondition regionConditionImpl;

    @Override
    public Response<QueryPagingListResponseBody> findList(BaseRequestBody body) {
        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);

        //查询总数
        int totalCount = this.regionConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<RegionDetailResponseBody> list = new ArrayList<>();
        Optional<List<Region>> clusterList = Optional.ofNullable(this.regionConditionImpl.findList(map));
        clusterList.ifPresent(x -> {
            x.forEach(y -> {
                //创建区域结果
                RegionDetailResponseBody clusterDetail = this.createRegionResult(y);
                list.add(clusterDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<RegionDetailResponseBody> find(String id) {
        //根据id查询区域信息
        Region region = this.regionConditionImpl.findById(id);
        if (null == region) {
            //不存在该区域
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //创建区域结果
        RegionDetailResponseBody result = this.createRegionResult(region);
        return new Response<>(result);
    }

    @Override
    public Response<RegionDetailResponseBody> add(AddRegionRequestBody body) {
        //通过区域名称查询区域信息
        String name = body.getName();
        Region region = this.regionConditionImpl.findByName(name);
        if (null != region) {
            //已存在该区域
            return new Response<>(SystemErrCode.ERRCODE_REGION_EXISTS, SystemErrCode.ERRMSG_REGION_EXISTS);
        }

        //通过companion查询区域信息
        String pushUrl = body.getPushUrl();
        Region regionCompare = this.regionConditionImpl.findByPushUrl(pushUrl);

        if (null != regionCompare) {
            //已存在该区域
            return new Response<>(SystemErrCode.ERRCODE_COMPANION_EXISTS, SystemErrCode.ERRMSG_COMPANION_EXISTS);
        }

        //创建区域
        Region newRegion = this.regionConditionImpl.add(body);

        //创建区域结果
        RegionDetailResponseBody result = this.createRegionResult(newRegion);
        return new Response<>(result);
    }

    @Override
    public Response<RegionDetailResponseBody> edit(EditRegionRequestBody body) {
        //根据id查询区域信息
        String id = body.getId();
        Region region = this.regionConditionImpl.findById(id);
        if (null == region) {
            //不存在该区域
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //根据id更新集群
        Region updateRegion = this.regionConditionImpl.updateById(id, body);

        //创建区域结果
        updateRegion.setName(region.getName());
        updateRegion.setCreateTime(region.getCreateTime());
        RegionDetailResponseBody result = this.createRegionResult(updateRegion);
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //根据id删除区域
        String id = body.getId();
        int success = this.regionConditionImpl.deleteById(id);
        if (success <= 0) {
            //不存在该区域
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }
        return new Response<>(null);
    }

    /**
     * 创建区域结果
     *
     * @param region
     * @return
     */
    private RegionDetailResponseBody createRegionResult(Region region) {
        RegionDetailResponseBody result = new RegionDetailResponseBody();
        result.setCreateTime(region.getCreateTime());
        result.setId(region.getId());
        result.setName(region.getName());
        result.setPushUrl(region.getPushUrl());
        result.setUpdateTime(region.getUpdateTime());
        return result;
    }
}
