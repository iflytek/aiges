package com.iflytek.ccr.polaris.cynosure.test.controller;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.user.LoginRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import org.junit.Assert;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

import javax.servlet.http.HttpSession;

/**
 * 登录
 *
 * @author sctang2
 * @create 2018-01-23 16:44
 **/
public class Utils {
    /**
     * 获取登录会话
     *
     * @param mockMvc
     * @return
     * @throws Exception
     */
    public static HttpSession getLoginSession(MockMvc mockMvc) throws Exception {
        LoginRequestBody body = new LoginRequestBody();
        body.setAccount("admin");
        body.setPassword("123456");
        return mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/user/login")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getRequest().getSession();
    }

    /**
     * 断言结果
     *
     * @param result
     */
    public static void assertResult(String result) {
        Response<Object> response = JSON.parseObject(result, Response.class);
        Assert.assertNotEquals(response.getCode().longValue(), 100001);
    }
}
