package com.iflytek.ccr.finder;

import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

public class TTTest {
    List<String > a = new ArrayList<>();
    List<String > b = new ArrayList<>();

    @BeforeClass
    public void before(){
        a.add("1");
        a.add("2");
        a.add("3");
        b.add("1");
        b.add("2");
        b.add("2");
    }
    @Test
    public void ca(){
        Assert.assertEquals(a.size()==b.size(),true);
    }

}
