package utils;

import java.io.IOException;
import java.io.InputStream;


public class StreamUtils {
    /**
     * 关闭输入流
     */
    public static void closeInputStream(InputStream is){
        if (is != null){
            try {
                is.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}
