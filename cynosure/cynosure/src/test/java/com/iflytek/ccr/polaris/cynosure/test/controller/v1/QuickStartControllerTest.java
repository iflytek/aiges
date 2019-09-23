package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddVersionRequestBodyByQuickStart;
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
 * 快速创建控制器测试
 *
 * @author sctang2
 * @create 2018-01-29 15:34
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class QuickStartControllerTest {
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
     * 快速新增版本
     *
     * @throws Exception
     */
    @Test
    public void test_quick_start_1_addservice() throws Exception {
        AddVersionRequestBodyByQuickStart body = new AddVersionRequestBodyByQuickStart();
        body.setProject("测试项目");
        body.setCluster("测试集群");
        body.setService("测试服务");
        body.setVersion("测试快速新增版本");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/quickStart/addVersion")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 快速新增服务版本
     *
     * @throws Exception
     */
    @Test
    public void test_quick_start_2_addserviceVersion() throws Exception {
//        ClassPathResource classPathResource1 = new ClassPathResource("1.yml");
//        MockMultipartFile multipartFile1 = new MockMultipartFile("file", "1.yml", "application/octet-stream", classPathResource1.getInputStream());
//
//        ClassPathResource classPathResource2 = new ClassPathResource("2.yml");
//        MockMultipartFile multipartFile2 = new MockMultipartFile("file", "2.yml", "application/octet-stream", classPathResource2.getInputStream());
//
//        ClassPathResource classPathResource3 = new ClassPathResource("3.yml");
//        MockMultipartFile multipartFile3 = new MockMultipartFile("file", "3.yml", "application/octet-stream", classPathResource3.getInputStream());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.fileUpload("/api/v1/quickStart/addServiceVersion")
//                .file(multipartFile1)
//                .file(multipartFile2)
//                .file(multipartFile3)
                .param("project", "project1")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("version", "1.0.0.2")
                .param("desc", "版本描述")
                .contentType(MediaType.MULTIPART_FORM_DATA)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON)).andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 新增服务配置
     *
     * @throws Exception
     */
    @Test
    public void test_quick_start_3_addserviceConfig() throws Exception {
        ClassPathResource classPathResource1 = new ClassPathResource("config/application.yml");
        MockMultipartFile multipartFile1 = new MockMultipartFile("file", "application.yml", "application/octet-stream", classPathResource1.getInputStream());

//        ClassPathResource classPathResource2 = new ClassPathResource("2.yml");
//        MockMultipartFile multipartFile2 = new MockMultipartFile("file", "2.yml", "application/octet-stream", classPathResource2.getInputStream());
//
//        ClassPathResource classPathResource3 = new ClassPathResource("3.yml");
//        MockMultipartFile multipartFile3 = new MockMultipartFile("file", "3.yml", "application/octet-stream", classPathResource3.getInputStream());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.fileUpload("/api/v1/quickStart/addServiceConfig")
                .file(multipartFile1)
//                .file(multipartFile2)
//                .file(multipartFile3)
                .param("project", "测试项目")
                .param("cluster", "测试集群")
                .param("service", "测试服务")
                .param("version", "测试版本")
                .param("desc", "测试服务配置")
                .contentType(MediaType.MULTIPART_FORM_DATA)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON)).andExpect(MockMvcResultMatchers.status().isOk())
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
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/quickStart/list")
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
    public void test_quick_start_5_list1() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/quickStart/list1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

}
