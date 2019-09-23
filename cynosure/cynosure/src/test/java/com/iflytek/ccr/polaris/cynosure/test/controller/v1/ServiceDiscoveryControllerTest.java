package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.*;
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
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.context.WebApplicationContext;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * 服务发现控制器测试
 *
 * @author sctang2
 * @create 2018-02-01 15:48
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class ServiceDiscoveryControllerTest {
    @Autowired
    private WebApplicationContext context;

    private MockMvc         mockMvc;
    private MockHttpSession session;

    @Before
    public void setup() {
        //初始化
        this.mockMvc = MockMvcBuilders.webAppContextSetup(this.context).build();
        this.session = new MockHttpSession();
    }

    /**
     * 新增服务
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_1_add() throws Exception {
        AddServiceDiscoveryRequestBody body = new AddServiceDiscoveryRequestBody();
        body.setProject("测试项目");
        body.setGroup("测试集群");
        body.setService("测试服务");
        body.setApiVersion("测试api版本");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/discovery/add")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询最近的服务发现列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_2_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/discovery/lastestList")
                .param("project", "test0803")
                .param("cluster", "test0803")
                .param("service","test0803")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务发现明细
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_3_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/discovery/detail")
                .param("project", "project3")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("apiVersion","测试版本")
                .param("region", "test")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑服务发现
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_4_edit() throws Exception {
        EditServiceDiscoveryRequestBody body = new EditServiceDiscoveryRequestBody();
        body.setProject("project1");
        body.setCluster("cluster1");
        body.setService("service1");
        body.setApiVersion("1.0.0.1");
        body.setRegion("test");
        body.setLoadbalance("roundrobin");

        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/discovery/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务提供者列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_5_provider() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/discovery/provider")
                .param("project", "test0803")
                .param("cluster", "test0803")
                .param("service", "test0803")
                .param("apiVersion","1.0")
                .param("region", "test")
                .param("isPage", "0")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑服务提供者
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_6_edit_provider() throws Exception {
        EditServiceDiscoveryProviderRequestBody body = new EditServiceDiscoveryProviderRequestBody();
        body.setProject("test0803");
        body.setCluster("test0803");
        body.setService("test0803");
        body.setApiVersion("1.0");
        body.setRegion("test");
        body.setAddr("80.1.86.21:311");
        body.setValid(true);
        body.setWeight(10);
        List<ServiceParam> params = new ArrayList<>();
        ServiceParam param1 = new ServiceParam();
        param1.setKey("key1");
        param1.setVal("value1");
        ServiceParam param2 = new ServiceParam();
        param2.setKey("key2");
        param2.setVal("value2");
        params.add(param1);
        params.add(param2);
        body.setParams(params);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/discovery/provider/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务订阅者列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_7_consumer() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/discovery/consumer")
                .param("project", "test0803")
                .param("cluster", "test0803")
                .param("service", "test0803")
                .param("apiVersion","1.0")
                .param("region", "test")
                .param("isPage", "0")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 更新反馈
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_8_feedback() throws Exception {
        ServiceDiscoveryFeedBackRequestBody body = new ServiceDiscoveryFeedBackRequestBody();
        body.setPushId("3717961262139703296");
        body.setProject("project3");
        body.setGroup("cluster1");
        body.setProvider("service1");
        body.setProviderVersion("1.0.0.0");
        body.setConsumer("service1");
        body.setConsumerVersion("1.0.0.0");
        body.setAddr("192.168.0.2");
        body.setUpdateStatus(1);
        body.setUpdateTime(new Date());
        body.setLoadStatus(1);
        body.setLoadTime(new Date());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/discovery/feedback")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务发现负载均衡列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_discovery_9_loadBalanceList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/discovery/loadBalanceList")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
