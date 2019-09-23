package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.user.AddUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.EditUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.LoginRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.ModifyPasswordRequestBody;
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
 * 用户模块控制器测试
 *
 * @author sctang2
 * @create 2018-01-13 13:51
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class UserControllerTest {
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
     * 登录
     *
     * @throws Exception
     */
    @Test
    public void test_user_1_login() throws Exception {
        LoginRequestBody body = new LoginRequestBody();
        body.setAccount("admin");
        body.setPassword("123456");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/login")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询用户列表
     *
     * @throws Exception
     */
    @Test
    public void test_user_2_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/user/list")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 新增用户
     *
     * @throws Exception
     */
    @Test
    public void test_user_3_add() throws Exception {
        AddUserRequestBody body = new AddUserRequestBody();
        body.setAccount("sctang");
        body.setEmail("sctang2@iflytek.com");
        body.setPhone("13739263609");
        body.setUserName("sctang");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/add")
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
     * 编辑用户
     *
     * @throws Exception
     */
    @Test
    public void test_user_4_edit() throws Exception {
        EditUserRequestBody body = new EditUserRequestBody();
        body.setId("3714413091023224832");
        body.setEmail("sctang2@iflytek.com");
        body.setPhone("13739263609");
        body.setUserName("sctang");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/edit")
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
     * 查询用户详情
     *
     * @throws Exception
     */
    @Test
    public void test_user_5_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/user/detail")
                .param("id", "3714413091023224832")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 修改密码
     *
     * @throws Exception
     */
    @Test
    public void test_user_6_modifyPassword() throws Exception {
        ModifyPasswordRequestBody body = new ModifyPasswordRequestBody();
        body.setOldPassword("123456");
        body.setNewPassword("123456");
        body.setConfirmPassword("123456");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/modifyPassword")
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
     * 查询个人信息
     *
     * @throws Exception
     */
    @Test
    public void test_user_7_info() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/user/info")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑个人信息
     *
     * @throws Exception
     */
    @Test
    public void test_user_8_editInfo() throws Exception {
        EditUserRequestBody body = new EditUserRequestBody();
        body.setEmail("sctang2@iflytek.com");
        body.setPhone("13739263609");
        body.setUserName("sctang");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/editInfo")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
