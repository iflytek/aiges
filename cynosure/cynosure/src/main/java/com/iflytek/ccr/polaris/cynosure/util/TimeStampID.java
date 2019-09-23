package com.iflytek.ccr.polaris.cynosure.util;

import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.text.*;
import java.util.Calendar;

/**
 * Created by admin
 * on 2018/10/24 10:11
 */
public class TimeStampID {
    /** The FieldPosition. */
    private static final FieldPosition HELPER_POSITION = new FieldPosition(0);
    /**This Format for format the data to special format. */
    private final static Format dateFormat = new SimpleDateFormat("YYYYMMddHHmmssS");
    /** This Format for format the number tospecial format. */
    private final static NumberFormat numberFormat = new DecimalFormat("0000");
    /**This int is the sequence number ,the default value is 0. */
    private static int seq = 0;
    private static final int MAX = 9999;

    /**
     * 生成唯一的时间戳ID
     * @return String
     */
    public static synchronized String getStampID() {
        Calendar rightNow = Calendar.getInstance();
        StringBuffer sb = new StringBuffer();
        dateFormat.format(rightNow.getTime(), sb,HELPER_POSITION);
        numberFormat.format(seq, sb, HELPER_POSITION);
        if (seq == MAX) {
            seq = 0;
        }else {
            seq++;
        }
        return sb.toString() + getServiceIP();
    }

    //获取本机ip的哈希值作为机器标识
    private static String getServiceIP(){
        InetAddress addr = null;
        try{
            addr = InetAddress.getLocalHost();
        }catch (UnknownHostException e){
            GlobalExceptionUtil.log(e);
        }
        String ip = addr.getHostAddress().toString();
        //地址取哈希
        return String.valueOf(ip.hashCode());
    }
}
