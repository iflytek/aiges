package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;
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

/**
 * 服务控制器测试
 *
 * @author sctang2
 * @create 2018-01-28 17:57
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class ServiceControllerTest {
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
    public void test_service_1_add() throws Exception {
        AddServiceRequestBody body = new AddServiceRequestBody();
        body.setName("service4");
        body.setDesc("服务4描述");
        body.setClusterId("3717225558505947136");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/add")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑服务
     *
     * @throws Exception
     */
    @Test
    public void test_service_2_edit() throws Exception {
        EditServiceRequestBody body = new EditServiceRequestBody();
        body.setId("3716423263660802048");
        body.setDesc("修改后的服务描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除服务
     *
     * @throws Exception
     */
    @Test
    public void test_service_3_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("3714708188751200256");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/delete")
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
     * 查询最近的服务列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_4_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/lastestList")
//                .param("project", null)
//                .param("cluster", null)
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务明细
     *
     * @throws Exception
     */
    @Test
    public void test_service_5_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/detail")
                .param("id", "3716423263660802048")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_6_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/list")
                .param("project", "project1")
                .param("cluster", "cluster1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
