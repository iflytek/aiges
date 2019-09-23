#### a json script for go

#### inline functions

funs|use
---|---
append(str...)string|拼接字符串
split(str,sep,n)[]string|按照给定的切割符号 sep切分字符串 
len(str)int|求字符串的长度
sprintf(str,obj...)string|obj...|格式化字符串
add(nubmer ...)number|number array|数字求和
isnil(object)bool|判断对象是否为空
and(bool...)bool|对一系列的bool变量做 and操作 相当于 && ，任何变量都可以视为bool值<br>
or(bool ...)bool|对一系列的bool变量做 or 操作 相当于 '\\' ，任何变量都可以视为bool值<br>
eq(objectA,objectB)bool|判断两者转换为字符串后是否相等
gt(numberA,numberB)bool| 大于 > ,numberA>numberB 返回true，否则返回false
ge(numberA,numberB)bool| 大于等于 >= ,numberA>=numberB 返回true，否则返回false
lt(numberA,numberB)bool| 小于 < ,numberA<numberB 返回true，否则返回false
le(numberA,numberB)bool| 小于等于 <= ,numberA<=numberB 返回true，否则返回false
not(bool)bool|非 ！  如果bool=true 返回false 否则返回true
in(str,strs...)bool |如果str 在 strs中一个，返回true ，否则返回false  
contains(str,sub)| 如果str包含sub片段 ，返回true，否则返回false
join(str ... ,sep)string| 拼接字符串，用sep分割
index(array,idx)obj| 返回数组的第i个元素，越界会panic
return(code,msg)|终止script的执行并返回一个error，包含code，和message，类型为*ErrorReturn


