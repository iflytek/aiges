if [ "$1" -eq "back" ]
then
    sed -i "s/github.com\/xfyun/git.xfyun.cn\/AIaaS/g" *
    sed -i "s/github.com\/xfyun/git.xfyun.cn\/AIaaS/g" server/*
    sed -i "s/github.com\/xfyun/git.xfyun.cn\/AIaaS/g" schemas/*
    sed -i "s/github.com\/xfyun/git.xfyun.cn\/AIaaS/g" common/*
    sed -i "s/github.com\/xfyun/git.xfyun.cn\/AIaaS/g" conf/*
    mv vendor/github.com/xfyun  vendor/git.xfyun.cn/AIaaS
else
    sed -i "s/git.xfyun.cn\/AIaaS/github.com\/xfyun/g" *
    sed -i "s/git.xfyun.cn\/AIaaS/github.com\/xfyun/g" server/*
    sed -i "s/git.xfyun.cn\/AIaaS/github.com\/xfyun/g" schemas/*
    sed -i "s/git.xfyun.cn\/AIaaS/github.com\/xfyun/g" common/*
    sed -i "s/git.xfyun.cn\/AIaaS/github.com\/xfyun/g" conf/*
    mv  vendor/git.xfyun.cn/AIaaS  vendor/github.com/xfyun
fi