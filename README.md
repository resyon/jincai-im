# JinCai-IM

IM 项目的实践.

## Feature

### bug 

-[ ] 同名房间创建冲突

### feature

- [x] 登录授权
- [x] 创建房间
- [x] 收发消息
- [x] 房间上线提醒
- [ ] 用户在线状态
- [ ] 加入消息序列号, 确保消息有序, 可靠
- [ ] 漫游消息 (持久化)
- [ ] HTTP fallback 
- [ ] 使用消息队列替代 Redis Pub-Sub 
- [ ] 横向扩容可行


## Doc

```
// 登录, 不存在则创建
POST /login --data {"username": "", "password":""}

// 以下使用 JWT 授权, 
header:
Authorization: Bearer ${token}

query:
token=${token}

// 创建房间
POST /auth/room?room_name=""

// 加入房间
PATCH /auth/room?room_id=""

// 连接聊天室
GET ws://hostname:port/auth/ws?token

```
