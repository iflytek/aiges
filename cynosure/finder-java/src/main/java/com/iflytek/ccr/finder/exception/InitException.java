package com.iflytek.ccr.finder.exception;

/**
 * 初始化异常
 */
public class InitException extends Exception {

    public InitException(String msg) {
        super(msg);
    }

    public InitException(int code, String msg) {
        super(msg);
    }

    public InitException(String msg, Throwable cause) {
        super(msg, cause);
    }
}
