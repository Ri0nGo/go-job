local key = KEYS[1]
local cntKey = key..":cnt"

-- 用户输入的验证码
local inputCode = ARGV[1]
-- 实际redis中存储的验证码
local actualCode = redis.call("get", key)
-- 验证次数
local cnt = tonumber(redis.call("get", cntKey))

if cnt == nil or cnt <= 0 then
    -- 验证码过期或验证码验证次数已耗尽
    return -1
end

if inputCode == actualCode then
    redis.call("del", cntKey)
    redis.call("del", key)
    return
else
    -- 验证码输入错误，验证次数减1
    redis.call("decr", cntKey)
    return -2
end