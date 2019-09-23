package com.iflytek.ccr.polaris.companion.common;

import org.apache.curator.framework.recipes.queue.QueueSerializer;

public class QueueItemSerializer implements QueueSerializer<String> {
    @Override
    public byte[] serialize(String item) {
        return item.getBytes();
    }

    @Override
    public String deserialize(byte[] bytes) {
        return new String(bytes);
    }
}
