package utils;

import com.iflytek.ccr.finder.value.SubscribeRequestValue;
import org.junit.Assert;
import org.junit.Test;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ListUtil {

    public static  ArrayList<SubscribeRequestValue> collectAsList(SubscribeRequestValue ...obj){
        ArrayList<SubscribeRequestValue> list = new ArrayList<SubscribeRequestValue>();
        for (SubscribeRequestValue o:obj){
            list.add(o);
        }
        return list;
    }
    public static  ArrayList collectAsArrayList(Object ...obj){
        ArrayList list = new ArrayList();
        for (Object o:obj){
            list.add(o);
        }
        return list;
    }

    public static boolean equals(List<String > a,List<String > b){
        if (a == null || b == null)return false;
        if (a.size() != b.size())return false;

        Map<String ,String > ma = new HashMap<String ,String >();
        Map<String ,String > mb = new HashMap<String ,String >();
        b.forEach(o->ma.put((String) o,"1"));
        a.forEach(o->mb.put((String) o,"1"));

        if (ma.size() != mb.size())return false;
        for(Map.Entry<String ,String > e:ma.entrySet()){
            String s = mb.get(e.getKey());
            if (s == null || !s.equals("1")){
                return false;
            }
        }
            return true;
    }

    @Test
    public void test(){
       List a = ListUtil.collectAsArrayList("1","2","3");
       List b = ListUtil.collectAsArrayList("1","2","3");
       List c = ListUtil.collectAsArrayList("1","2","3","4");
       List d = ListUtil.collectAsArrayList("1","2","4","3");
       List e = ListUtil.collectAsArrayList("1","2","3","4");
       List f = ListUtil.collectAsArrayList("1","2","2","4");
        Assert.assertEquals(true,equals(a,b));
        Assert.assertEquals(true,equals(c,d));
        Assert.assertEquals(true,equals(e,d));
        Assert.assertEquals(false,equals(e,f));
    }
}
