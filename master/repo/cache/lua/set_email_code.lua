local key = KEYS[1]
-- 验证码使用次数
local cntKey = key..":cnt"

-- 验证码
local code = ARGV[1]
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key存在，但是没有设置过期时间
    return -2
elseif ttl == -2 or ttl < 1740 then
    -- key不存在，获取过期时间使用超过60秒
    redis.call("set", key, code)
    -- 设置过期时间30分钟
    redis.call("expire", key, 1800)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 1800)
    return 0
else
    -- 发送太频繁了
    return -1
end