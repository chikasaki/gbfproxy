# 生成aes密钥
openssl rand -hex 32 | xxd -r -p | base64 | head -c 32 > aes.key

echo -n '输入用户名:'
read username

echo -n '输入密码:'
read password

echo -n '输入端口号(默认:53540):'
read port

if [ -z "${port}" ]; then
  port=53540
fi

# 复制生成新的配置文件，使用前面填写的用户名、密码、端口号进行替换
cp server_config_template.toml server_config.toml
sed -i "s/\[username\]/${username}/g" server_config.toml
sed -i "s/\[password\]/${password}/g" server_config.toml
sed -i "s/\[port\]/${port}/g" server_config.toml

# 如果有老进程存在的话，杀掉
pid=$(pgrep -f server_proxy)
if [ -n "$pid" ]; then
  kill -9 $pid
  sleep 1
fi

nohup ./server_proxy > nohup.out 2>&1 &
