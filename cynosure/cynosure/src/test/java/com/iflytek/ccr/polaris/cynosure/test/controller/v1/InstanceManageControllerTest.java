package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;
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
 * 推送实例控制器测试
 *
 * @author sctang2
 * @create 2018-01-29 10:23
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class InstanceManageControllerTest {
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
     * 查询推送实例明细
     *
     * @throws Exception
     */
    @Test
    public void test_instance_manage_1_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/instanceManage/detail")
                .param("id", "3740064022783852544")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询推送实例列表
     *
     * @throws Exception
     */
    @Test
    public void test_instance_manage_2_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/instanceManage/List")
                .param("versionId", "3776343244166660096")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑推送实例
     *
     * @throws Exception
     */
    @Test
    public void test_instance_manage_3_edit() throws Exception {
        EditInstanceRequestBody body = new EditInstanceRequestBody();
        body.setGrayId("3783559869148168192");
        body.setVersionId("3783219485607985152");
        body.setContent("2222");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/instanceManage/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 4
     * 新增灰度组时查询推送实例
     *
     * @throws Exception
     */
    @Test
    public void test_instance_manage_4_appointList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/instanceManage/appointList")
                .param("project", "qq")
                .param("cluster", "qq")
                .param("service", "qq")
                .param("version", "1.0")
                .param("versionId", "3783219485607985152")
                .param("grayId", "0")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
