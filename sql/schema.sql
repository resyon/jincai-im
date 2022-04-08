use jincai;
# type User struct {
# 	Id       int    `json:"id" gorm:"column:id"`
# 	Username string `json:"username" gorm:"column:username"`
# 	Password string `json:"password" gorm:"column:password"`
# }

create table if not exists `user` (
    id int primary key auto_increment,
    username varchar(255),
    password varchar(255)
);


-- type Room struct {
-- 	RoomId string `json:"room_id" gorm:"room_id"`
-- 	// id of user who owns the room
-- 	OwnerId  int    `json:"owner_id" gorm:"owner_id"`
-- 	RoomName string `json:"room_name" gorm:"room_name"`
-- }
create table if not exists room (
    room_id char(36) primary key,
    owner_id int,
    room_name varchar(255),
    unique(owner_id, room_name),
    constraint foreign key user_id (owner_id) REFERENCES user(id)
);

# type RoomMate struct {
# 	RoomId string `json:"room_id" gorm:"room_id"`
# 	UserId int    `json:"user_id" gorm:"user_id"`
# }

create table if not exists room_mate(
    room_id char(36),
    user_id int,
    constraint foreign key room_mate_room (room_id) REFERENCES room(room_id),
    constraint foreign key room_mate_user (user_id) REFERENCES user(id)
);

# type Message struct {
# 	Id          int64  `json:"id" gorm:"column:id"`
# 	Time        int64  `json:"time" gorm:"column:time"`
# 	UserId      int    `json:"user_id" gorm:"column:user_id"`
# 	RoomId      string `json:"room_id" gorm:"column:room_id"` // 64B
# 	MessageType uint8  `json:"message_type" gorm:"message_type"`
# 	HasRead     bool   `json:"had_read" gorm:"column:has_read"`
# 	HasSend     bool   `json:"has_send" gorm:"column:has_send"`
# 	Text        string `json:"text" gorm:"column:text"`
# }
create table if not exists message(
    id bigint primary key,
    time bigint,
    user_id int,
    room_id char(36),
    message_type tinyint,
    has_read bool,
    has_send bool,
    text text,

    constraint foreign key  msg_user_id (user_id) references `user`(`id`),
    constraint foreign key msg_room_id (room_id) references `room`(`room_id`)
);

insert into user(`id`,`username`
)values (
         0, 'SYSTEM'
        );

insert into room(`room_id`, `room_name`, `owner_id`)values ('__sys_channel', 'SYS_CHANNEL', 0);