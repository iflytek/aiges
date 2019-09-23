package com.iflytek.ccr.polaris.companion.utils;

import java.util.ArrayList;
import java.util.List;

public class ArrayUtilTest {
    public static void main(String[] args) {
        List<String> a=new ArrayList<>();
        List<String> b=new ArrayList<>();
        a.add("1");
        a.add("2");
        a.add("3");

        b.add("1");
        b.add("2");
        b.add("2");

        System.out.println(ArrayUtils.equals(b,a));
    }
}
