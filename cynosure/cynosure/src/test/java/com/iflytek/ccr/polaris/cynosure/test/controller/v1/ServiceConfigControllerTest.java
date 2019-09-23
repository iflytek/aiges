package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.test.controller.Utils;
import org.junit.Before;
import org.junit.FixMethodOrder;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.MethodSorters;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.mock.web.MockHttpSession;
import org.springframework.test.context.junit4.SpringRunner;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultHandlers;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.context.WebApplicationContext;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * 服务配置控制器测试
 *
 * @author sctang2
 * @create 2018-01-30 16:02
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class ServiceConfigControllerTest {
    @Autowired
    private WebApplicationContext context;

    private MockMvc mockMvc;
    private MockHttpSession session;

    @Before
    public void setup() {
        //初始化
        this.mockMvc = MockMvcBuilders.webAppContextSetup(this.context).build();
        this.session = new MockHttpSession();
    }

    /**
     * 查询最近的配置列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_1_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/config/lastestList")
                .param("project", "project3")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("version", "1.0.0.1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_2_edit() throws Exception {
        EditServiceConfigRequestBody body = new EditServiceConfigRequestBody();
        body.setId("3753486267568881664");
        body.setContent("my friend xiBigBig");
        body.setDesc("修改配置后的描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询配置明细
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_3_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/config/detail")
                .param("id", "3716882253993738240")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 推送配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_4_push() throws Exception {
        PushServiceConfigRequestBody body = new PushServiceConfigRequestBody();
        body.setId("3753486267610824705");
        List<String> regionIds = new ArrayList<>();
        regionIds.add("3777763353183649792");
        body.setRegionIds(regionIds);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/push")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 批量推送配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_5_batchPush() throws Exception {
        BatchPushServiceConfigRequestBody body = new BatchPushServiceConfigRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3717226752750125056");
        ids.add("3717226752750125057");
        ids.add("3717226752750125058");
        body.setIds(ids);
        List<String> regionIds = new ArrayList<>();
        regionIds.add("3717224882795184128");
        regionIds.add("3717225000453799936");
        body.setRegionIds(regionIds);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/batchPush")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_6_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("3780732398677786624");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/delete")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 批量删除配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_7_batchdelete() throws Exception {
        IdsRequestBody body = new IdsRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3697214230441754624");
        ids.add("3708849106630737920");
        body.setIds(ids);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/batchDelete")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询配置历史列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_8_historyList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/config/historyList")
                .param("project", "GG")
                .param("cluster", "GG")
                .param("service", "GG")
                .param("version", "1.0")
                .param("isPage", "1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 回滚配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_9_rollback() throws Exception {
        IdsRequestBody body = new IdsRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3731673452474531840");
        ids.add("3731673452474531841");
        body.setIds(ids);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/rollback")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 更新反馈
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_10_feedback() throws Exception {
        ServiceConfigFeedBackRequestBody body = new ServiceConfigFeedBackRequestBody();
        body.setPushId("3783233233995431953");
        body.setProject("qq");
        body.setGroup("qq");
        body.setService("qq");
        body.setVersion("2.0");
        body.setConfig("test2.yml");
        body.setAddr("55.11.2.2:7010");
        body.setUpdateStatus(1);
        body.setUpdateTime(new Date());
        body.setLoadStatus(1);
        body.setLoadTime(new Date());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/feedback")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询订阅者列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_11_consumer() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/config/consumer")
                .param("project", "测试项目")
                .param("cluster", "测试集群")
                .param("service", "测试服务")
                .param("version", "测试版本")
                .param("region", "测试区域")
                .param("grayId", "0")

                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除配置
     *
     * @throws Exception
     */
    @Test
    public void test_service_config_12_download() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/config/download")
                .param("id", "3783228579949576192")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        System.out.println(result);
        Utils.assertResult(result);
    }

}
