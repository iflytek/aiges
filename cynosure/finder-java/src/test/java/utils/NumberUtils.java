package utils;

public class NumberUtils {
    /**
     * 判断字符串是否是数字串
     * @param s
     * @return
     */

    public static boolean isNumber(String s){

        char cs[] = s.toCharArray();
        for (char a:cs){
            if (a < '0' || a > '9'){
                return false;
            }
        }

        return true;
    }

    /**
     * 判断字符串是否是整数
     * @param s
     * @return
     */
    public static boolean isInteger(String s){
        try {
            Integer.parseInt(s);
        }catch (Exception e){
            return false;
        }
        return true;
    }

}
