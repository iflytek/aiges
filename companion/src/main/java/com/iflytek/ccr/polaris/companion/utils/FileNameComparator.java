package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.polaris.companion.common.Constants;
import org.apache.http.client.utils.DateUtils;

import java.util.Comparator;
import java.util.Date;

public class FileNameComparator implements Comparator<String> {
    String[] dateFormats = new String[]{Constants.DATE_PATTERN};

    public static void main(String[] args) {
        System.out.println(DateUtils.parseDate("20180101", new String[]{Constants.DATE_PATTERN}));
    }

    @Override
    public int compare(String o1, String o2) {
        Date date1 = DateUtils.parseDate(o1, dateFormats);
        Date date2 = DateUtils.parseDate(o2, dateFormats);
        return date2.compareTo(date1);
    }
}
