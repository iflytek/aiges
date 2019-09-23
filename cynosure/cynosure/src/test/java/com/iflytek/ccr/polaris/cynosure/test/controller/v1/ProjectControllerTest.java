package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectMemberRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.DeleteProjectMemberRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.EditProjectRequestBody;
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
 * 项目控制器测试
 *
 * @author sctang2
 * @create 2017-11-28 12:06
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class ProjectControllerTest {
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
     * 新增项目
     *
     * @throws Exception
     */
    @Test
    public void test_project_1_add() throws Exception {
        AddProjectRequestBody body = new AddProjectRequestBody();
        body.setName("project2");
        body.setDesc("project2描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/project/add")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑项目
     *
     * @throws Exception
     */
    @Test
    public void test_project_2_edit() throws Exception {
        EditProjectRequestBody body = new EditProjectRequestBody();
        body.setId("3714708188751200256");
        body.setDesc("修改项目后的描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/project/edit")
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
     * 删除项目
     *
     * @throws Exception
     */
    @Test
    public void test_project_3_delete() throws Exception {
        IdRequestBody body = new IdRequestBody();
        body.setId("3714708188751200256");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/project/delete")
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
     * 查询项目明细
     *
     * @throws Exception
     */
    @Test
    public void test_project_4_find() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/project/detail")
                .param("id", "3714708188751200256")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询项目列表
     *
     * @throws Exception
     */
    @Test
    public void test_project_5_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/project/list")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 新增项目成员
     *
     * @throws Exception
     */
    @Test
    public void test_project_6_member_add() throws Exception {
        AddProjectMemberRequestBody body = new AddProjectMemberRequestBody();
        body.setId("3714708188751200256");
        body.setAccount("sctang");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/project/member/add")
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
     * 删除项目成员
     *
     * @throws Exception
     */
    @Test
    public void test_project_7_member_delete() throws Exception {
        DeleteProjectMemberRequestBody body = new DeleteProjectMemberRequestBody();
        body.setId("3714708188751200256");
        body.setUserId("3688919024172793856");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/project/member/delete")
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
     * 查询项目成员列表
     *
     * @throws Exception
     */
    @Test
    public void test_project_8_member_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/project/member/list")
                .param("id", "3688919024172793856")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
