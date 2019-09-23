package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.polaris.companion.utils.ArrayUtils;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

public class ArrayUtilsTest {
    static List<String> a=new ArrayList<>();
    static List<String> b=new ArrayList<>();
    @BeforeClass
    public static  void before(){
        a.add("1");
        a.add("3");
        a.add("2");
        b.add("1");
        b.add("1");
        b.add("3");



    }
    @Test
    public void go(){
        Assert.assertEquals(ArrayUtils.equals(b,a),true);
        Assert.assertEquals(ArrayUtils.equals(a,b),false);
    }
}
