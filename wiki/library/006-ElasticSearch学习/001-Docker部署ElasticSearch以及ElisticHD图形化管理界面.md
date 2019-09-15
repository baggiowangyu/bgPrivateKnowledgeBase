# Docker部署ElasticSearch以及ElisticHD图形化管理界面

## Docker命令

```
docker run -p 9200:9200 -d --name elasticsearch elasticsearch
docker run -p 9800:9800 -d --link elasticsearch:demo containerize/elastichd
```

在浏览器中访问 http://[docker-ip]:9800
连接到 http://demo:9200
