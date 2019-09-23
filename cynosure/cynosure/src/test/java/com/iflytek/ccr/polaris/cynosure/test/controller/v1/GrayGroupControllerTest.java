package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.DeleteGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.test.controller.Utils;
import org.junit.Before;
import org.junit.FixMethodOrder;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.MethodSorters;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.MediaType;
import org.springframework.mock.web.MockHttpSession;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.test.context.junit4.SpringRunner;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultHandlers;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.context.WebApplicationContext;

/**
 * 灰度组控制器测试
 *
 * @author sctang2
 * @create 2018-01-29 10:23
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class GrayGroupControllerTest {
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
     * 新增灰度组
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_1_add() throws Exception {
        ClassPathResource classPathResource1 = new ClassPathResource("1.yml");
        MockMultipartFile multipartFile1 = new MockMultipartFile("file", "1.yml", "application/octet-stream", classPathResource1.getInputStream());

        ClassPathResource classPathResource2 = new ClassPathResource("2.yml");
        MockMultipartFile multipartFile2 = new MockMultipartFile("file", "2.yml", "application/octet-stream", classPathResource2.getInputStream());

        ClassPathResource classPathResource3 = new ClassPathResource("3.yml");
        MockMultipartFile multipartFile3 = new MockMultipartFile("file", "3.yml", "application/octet-stream", classPathResource3.getInputStream());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.fileUpload("/api/v1/gray/add")
                .file(multipartFile1)
                .file(multipartFile2)
                .file(multipartFile3)
                .param("project", "测试项目")
                .param("cluster", "测试集群")
                .param("service", "测试服务")
                .param("version", "测试版本")
                .param("versionId", "3776343244166660096")
                .param("desc", "测试灰度组")
                .param("name", "测试灰度组2")
                .contentType(MediaType.MULTIPART_FORM_DATA)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON)).andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询灰度组列表
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_2_list() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/gray/list")
                .param("versionId", "3763627581165797376")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询灰度组明细
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_3_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/gray/detail")
                .param("id", "3781019715757932544")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询最近的灰度组列表(取消)
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_4_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/gray/lastestList")
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
     * 删除灰度组
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_5_delete() throws Exception {
        DeleteGrayGroupRequestBody body = new DeleteGrayGroupRequestBody();
        body.setProject("project1");
        body.setCluster("cluster1");
        body.setService("service1");
        body.setVersion("1.0.0.0");
        body.setGrayId("3782772222335123456");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/gray/delete")
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
     * 编辑版本
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_6_edit() throws Exception {
        GrayGroupDetailResponseBody body = new GrayGroupDetailResponseBody();
        body.setId("3740064022783852544");
        body.setDesc("23");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/gray/edit")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 查询灰度配置订阅者列表
     *
     * @throws Exception
     */
    @Test
    public void test_gray_group_7_consumer() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/gray/consumer")
                .param("project", "project1")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("version", "version")
                .param("region", "R1")
                .param("grayId", "3740064022783852544")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
