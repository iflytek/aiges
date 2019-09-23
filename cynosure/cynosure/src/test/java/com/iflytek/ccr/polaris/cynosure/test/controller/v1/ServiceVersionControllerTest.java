package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
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
 * 服务版本控制器测试
 *
 * @author sctang2
 * @create 2018-01-29 10:23
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class ServiceVersionControllerTest {
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
     * 新增版本
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_1_add() throws Exception {
        AddServiceVersionRequestBody body = new AddServiceVersionRequestBody();
        body.setServiceId("3785028614114770944");
        body.setVersion("1.0.0.1.fsghdjftj");
        body.setDesc("1.0.0.1描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/version/add")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑版本
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_2_edit() throws Exception {
        EditServiceVersionRequestBody body = new EditServiceVersionRequestBody();
        body.setId("3716458136186388480");
        body.setDesc("修改后的版本描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/version/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 删除版本
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_3_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("3714708188751200256");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/version/delete")
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
     * 查询最近的服务版本列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_4_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/version/lastestList")
                .param("project", "project1")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务版本明细
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_5_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/version/detail")
                .param("id", "3716458136186388480")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询服务版本列表
     *
     * @throws Exception
     */
    @Test
    public void test_service_version_6_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/service/version/list")
                .param("project", "project1")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
