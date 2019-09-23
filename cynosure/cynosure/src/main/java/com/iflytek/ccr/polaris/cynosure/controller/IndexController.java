package com.iflytek.ccr.polaris.cynosure.controller;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * 首页
 *
 * @author sctang2
 * @create 2017-12-01 17:21
 **/
@RestController
public class IndexController {
    @RequestMapping(value = "/", method = {RequestMethod.POST, RequestMethod.GET})
    public void index(HttpServletResponse response) throws IOException {
        response.sendRedirect("index.html");
    }
}
