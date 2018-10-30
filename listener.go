package main

type iListener interface {
	onProcess(msg *Msg)
}

func addListener(newConn *myConn) {
	newConn.add(&chat{conn: newConn})
	newConn.add(&readed{conn: newConn})
	newConn.add(&offLineMsg{conn: newConn})
	newConn.add(&addFriendRequest{conn: newConn})
	newConn.add(&addFriendResponse{conn: newConn})
	newConn.add(&deltFriend{conn: newConn})
	newConn.add(&getFriends{conn: newConn})
	newConn.add(&remarkFrindListener{conn: newConn})
	newConn.add(&addFriendGroupListener{conn: newConn})
	newConn.add(&remarkFrindGroupListener{conn: newConn})
	newConn.add(&deltFrindGroupListener{conn: newConn}) //登录完才能这些操作
}

func basicListener(newConn *myConn) { //新用户只能登录注册
	newConn.add(&login{conn: newConn})
	newConn.add(&regis{conn: newConn})
}

func (this *myConn) add(l iListener) {
	this.listenList = append(this.listenList, l)
}
