package com.iflytek.ccr.polaris.companion.utils;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class ArrayUtils {


    /**
     * 判断两个集合的内容是否相等
     *
     * @param a
     * @param b
     * @return
     */
    public static boolean equals(List<String> a, List<String> b) {
        boolean flag = true;
        if (null == a || null == b || a.size() != b.size()) {
            flag = false;
        } else {
            Set<String> a1 = new HashSet<>();
            a1.addAll(a);
            Set<String> a2 = new HashSet<>();
            a2.addAll(b);
            if (a1.size() != a2.size()) {
                return flag;
            }
            for (String temp : a1) {
                flag = false;
                for (String tempB : a2) {
                    if (temp.equals(tempB)) {
                        flag = true;
                        break;
                    }
                }
                if (!flag) {
                    break;
                }
            }
        }
        return flag;
    }
}
