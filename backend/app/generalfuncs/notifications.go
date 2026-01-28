package generalfuncs

import (
	"log"
	"strconv"
	"time"

	"social-network/app/handlers/websocket"
	"social-network/app/models"
	"social-network/db"
)

// // CREATE TABLE notifications (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     user_id INTEGER NOT NULL,
//     follow_request_id INTEGER,
//     type TEXT NOT NULL CHECK(type IN ('follow_request','new_follower','user_accepted_follow','new_comment', 'group_invitation', 'group_join_request_approved', 'event_created','group_join_request','event_invitation')),
//     sender_id INTEGER,
//     sender_name TEXT,
//     sender_avatar TEXT,
//     post_id INTEGER,
//     group_id INTEGER,
//     group_name TEXT,
//     event_id INTEGER,
//     is_read BOOLEAN DEFAULT 0,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
//     FOREIGN KEY (sender_id) REFERENCES users(id),
//     FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
//     FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
//     FOREIGN KEY (follow_request_id) REFERENCES follow_user_requests(id) ON DELETE CASCADE
// );

func CreateNotification(n models.Notification) error {

	//horizontal expansion WORK BETTER
	_, err := db.Database.Exec(`INSERT INTO notifications 
	(user_id, follow_request_id, type, sender_id, sender_name, sender_avatar, post_id, group_id, group_name, event_id, is_read, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID,
		n.FollowRequestID,
		n.Type,
		n.SenderID,
		n.SenderName,
		n.SenderAvatar,
		n.PostID,
		n.GroupID,
		n.GroupName,
		n.EventID,
		false,
		time.Now(),
	)

	//check if user is online and send websocket notification
	log.Printf("Creating notification for user %d of type %s", n.UserID, n.Type)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return err
	}
	if websocket.IsUserOnline(strconv.Itoa(n.UserID)) {
		// for notification just a signal is enough
		websocket.SendToUser(
			strconv.Itoa(n.UserID),
			websocket.WebSocketMessage{
				Type: "notification",
				Data: "new_notification",
			},
		)

	}

	return err
}

// keep it simple for now and query the members of the group each time
// exclude the user who triggered the notification
// to avoid sending notification to self
// but i have to add the ws message for user that sent the message
func GetGroupMemberIDs(groupID int, excludeUserID int) ([]int, error) {
	var memberIDs []int
	rows, err := db.Database.Query(`SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ?`, groupID, excludeUserID)
	if err != nil {
		log.Println("Error querying group members:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Println("Error scanning group member ID:", err)
			continue
		}
		memberIDs = append(memberIDs, userID)
	}
	return memberIDs, nil
}
