package com.iflytek.ccr.finder.utils;

import java.util.Collection;
import java.util.Iterator;

public class StringUtils {

    public static String join(Collection var0, String var1) {
        StringBuffer var2 = new StringBuffer();

        for (Iterator var3 = var0.iterator(); var3.hasNext(); var2.append((String) var3.next())) {
            if (var2.length() != 0) {
                var2.append(var1);
            }
        }

        return var2.toString();
    }

    public static boolean isNullOrEmpty(String str) {
        return null == str || str.isEmpty();
    }

    public static boolean isNOtNullOrEmpty(String str) {
        return !(null == str || str.isEmpty());
    }

    public static boolean isEmpty(Collection coll) {
        return coll == null || coll.isEmpty();
    }

    public static boolean isNotEmpty(Collection coll) {
        return !isEmpty(coll);
    }
}
