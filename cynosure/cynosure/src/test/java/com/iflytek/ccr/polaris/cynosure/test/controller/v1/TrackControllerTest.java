package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
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
import java.util.List;

/**
 * 轨迹控制器测试
 *
 * @author sctang2
 * @create 2018-02-02 17:03
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class TrackControllerTest {
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
     * 查询最近的配置推送轨迹列表
     *
     * @throws Exception
     */
    @Test
    public void test_track_1_config_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/track/config/lastestList")
                .param("project", "测试项目")
                .param("cluster", "测试集群")
                .param("service", "测试服务")
                .param("version", "测试版本")
                .param("filterGray", "-1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询配置推送轨迹明细
     *
     * @throws Exception
     */
    @Test
    public void test_track_2_config_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/track/config/detail")
                .param("pushId", "3783232632330911744")
                .param("isGray", "No")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除配置推送轨迹
     *
     * @throws Exception
     */
    @Test
    public void test_track_3_config_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("3717229737668509696");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/track/config/delete")
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
     * 批量删除配置推送轨迹
     *
     * @throws Exception
     */
    @Test
    public void test_track_4_config_batch_delete() throws Exception {
        IdsRequestBody body = new IdsRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3717229385577660416");
        ids.add("3717548557121617920");
        body.setIds(ids);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/track/config/batchDelete")
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
     * 查询最近的服务发现推送轨迹列表
     *
     * @throws Exception
     */
    @Test
    public void test_track_5_discovery_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/track/discovery/lastestList")
                .param("project", "project3")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("version", "1.0.0.0")
                .param("currentPage", "1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务 发现推送轨迹明细
     *
     * @throws Exception
     */
    @Test
    public void test_track_6_discovery_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/track/discovery/detail")
                .param("pushId", "3785389906340085760")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除服务发现推送轨迹
     *
     * @throws Exception
     */
    @Test
    public void test_track_7_config_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/track/discovery/delete")
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
     * 批量删除服务发现推送轨迹
     *
     * @throws Exception
     */
    @Test
    public void test_track_8_config_batch_delete() throws Exception {
        IdsRequestBody body = new IdsRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("");
        body.setIds(ids);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/track/discovery/batchDelete")
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
     * 快速查询
     *
     * @throws Exception
     */
    @Test
    public void test_quick_start_4_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/track/discovery/list")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
