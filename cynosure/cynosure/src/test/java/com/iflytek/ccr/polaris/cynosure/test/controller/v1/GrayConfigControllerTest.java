package com.iflytek.ccr.polaris.cynosure.test.controller.v1;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.BatchPushServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.EditServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.PushServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.ServiceConfigFeedBackRequestBody;
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

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * 服务灰度配置控制器测试
 *
 * @author sctang2
 * @create 2018-01-30 16:02
 **/
@RunWith(SpringRunner.class)
@SpringBootTest
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class GrayConfigControllerTest {
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
     * 查询最近的配置列表
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_1_lastestList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/grayConfig/lastestList")
                .param("project", "project3")
                .param("cluster", "cluster1")
                .param("service", "service1")
                .param("version", "1.0.0.1")
                .param("gray", "test1")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 编辑配置
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_2_edit() throws Exception {
        EditServiceConfigRequestBody body = new EditServiceConfigRequestBody();
        body.setId("3782760274025512960");
        body.setContent("testAgin");
        body.setDesc("修改灰度配置后的描述");
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/grayConfig/edit")
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
     * 查询配置明细
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_3_detail() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/grayConfig/detail")
                .param("id", "3716882253993738240")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 推送配置
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_4_push() throws Exception {
        PushServiceConfigRequestBody body = new PushServiceConfigRequestBody();
        body.setId("3780959622257442816");
        List<String> regionIds = new ArrayList<>();
        regionIds.add("3777763353183649792");
        body.setRegionIds(regionIds);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/grayConfig/push")
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
     * 批量推送配置
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_5_batchPush() throws Exception {
        BatchPushServiceConfigRequestBody body = new BatchPushServiceConfigRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3780718589976248320");
        ids.add("3780718590009802752");
        ids.add("3780718590009802753");
        body.setIds(ids);
        List<String> regionIds = new ArrayList<>();
        regionIds.add("3777763353183649792");
        body.setRegionIds(regionIds);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/grayConfig/batchPush")
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
     * 查询灰度历史列表
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_6_historyList() throws Exception {
        String result = this.mockMvc.perform(MockMvcRequestBuilders.get("/api/v1/grayConfig/historyList")
                .param("project", "polaris_test")
                .param("cluster", "polaris_test")
                .param("service", "polaris_test")
                .param("version", "1.0")
                .contentType(MediaType.APPLICATION_JSON)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 回滚配置
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_9_rollback() throws Exception {
        IdsRequestBody body = new IdsRequestBody();
        List<String> ids = new ArrayList<>();
        ids.add("3731673452474531840");
        ids.add("3731673452474531841");
        body.setIds(ids);
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/rollback")
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
     * 更新反馈
     *
     * @throws Exception
     */
    @Test
    public void test_gray_config_10_feedback() throws Exception {
        ServiceConfigFeedBackRequestBody body = new ServiceConfigFeedBackRequestBody();
        body.setPushId("3717229737668509696");
        body.setProject("project3");
        body.setGroup("cluster1");
        body.setService("service1");
        body.setVersion("1.0.0.0");
        body.setConfig("1.yml");
        body.setAddr("192.168.0.2");
        body.setUpdateStatus(1);
        body.setUpdateTime(new Date());
        body.setLoadStatus(1);
        body.setLoadTime(new Date());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/service/config/feedback")
                .content(JSON.toJSONString(body))
                .contentType(MediaType.APPLICATION_JSON)
                .accept(MediaType.APPLICATION_JSON))
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andDo(MockMvcResultHandlers.print())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }

    /**
     * 新增灰度配置
     *
     * @throws Exception
     */
    @Test
    public void test_quick_start_3_addGrayConfig() throws Exception {
        ClassPathResource classPathResource1 = new ClassPathResource("1.yml");
        MockMultipartFile multipartFile1 = new MockMultipartFile("file", "1.yml", "application/octet-stream", classPathResource1.getInputStream());

        ClassPathResource classPathResource2 = new ClassPathResource("2.yml");
        MockMultipartFile multipartFile2 = new MockMultipartFile("file", "2.yml", "application/octet-stream", classPathResource2.getInputStream());

        ClassPathResource classPathResource3 = new ClassPathResource("3.yml");
        MockMultipartFile multipartFile3 = new MockMultipartFile("file", "3.yml", "application/octet-stream", classPathResource3.getInputStream());
        String result = this.mockMvc.perform(MockMvcRequestBuilders.fileUpload("/api/v1/grayConfig/addGrayConfig")
                .file(multipartFile1)
                .file(multipartFile2)
                .file(multipartFile3)
                .param("project", "测试项目")
                .param("cluster", "测试集群")
                .param("service", "测试服务")
                .param("version", "测试版本")
                .param("versionId", "3776343244166660096")
                .param("desc", "测试sdk")
                .param("grayId", "3780699522712207360")
                .param("name", "test")
                .contentType(MediaType.MULTIPART_FORM_DATA)
                .session((MockHttpSession) Utils.getLoginSession(this.mockMvc))
                .accept(MediaType.APPLICATION_JSON)).andExpect(MockMvcResultMatchers.status().isOk())
                .andReturn().getResponse().getContentAsString();
        Utils.assertResult(result);
    }
}
